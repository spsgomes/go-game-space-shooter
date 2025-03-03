package game

import (
	"go-game-space-shooter/internal/audio"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	window_padding  = 40
	background_size = 256
)

var max_enemies_per_wave int

func (g *Game) Update() error {
	g.player.Update(g)

	// Enemy: Update
	if len(g.enemies) > 0 {
		for _, enemy := range g.enemies {
			enemy.Update(g, g.player)
		}
	}

	// Projectile: Update
	if len(g.projectiles) > 0 {
		for _, projectile := range g.projectiles {
			projectile.Update(g)
		}
	}

	// Enemy: Spawn timer
	g.enemySpawnTimer.Update()
	if g.enemySpawnTimer.IsReady() {
		g.enemySpawnTimer.Reset()

		// Only spawn enemies if the player isn't dead
		if !g.player.disabled {
			g.enemies = SpawnEnemies(g.random, g.enemies, max_enemies_per_wave)
		}
	}

	// A 1-second timer to check for bounds, and HP
	g.oneSecondTimer.Update()
	if g.oneSecondTimer.IsReady() {
		g.oneSecondTimer.Reset()

		// Check Projectiles
		if len(g.projectiles) > 0 {
			var tmp []*Projectile
			for _, projectile := range g.projectiles {
				if !projectile.IsOutOfBounds() && !projectile.disabled {
					tmp = append(tmp, projectile)
				}
			}
			g.projectiles = tmp
		}

		// Check Enemies enabled
		if len(g.enemies) > 0 {
			var tmp []*Enemy
			for _, enemy := range g.enemies {
				if !enemy.disabled {
					tmp = append(tmp, enemy)
				}
			}
			g.enemies = tmp
		}
	}

	// Save Game on Death
	if g.player.disabled && !g.hasSavedOnDeath {
		g.hasSavedOnDeath = true

		_, err := g.save.Save(g)
		if err != nil {
			HandleError(err)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Background Draw
	g.ui.DrawBackground(screen)

	// Projectile: Draw
	if len(g.projectiles) > 0 {
		for _, projectile := range g.projectiles {
			projectile.Draw(screen)
		}
	}

	// Player: Draw
	g.player.Draw(screen)

	// Enemy: Draw
	if len(g.enemies) > 0 {
		for _, enemy := range g.enemies {
			enemy.Draw(screen)
		}
	}

	// Ui: Draw
	g.ui.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

var Configs map[string]string

func NewGame(configs map[string]string) *Game {

	Configs = configs

	// Game music
	music, err := audio.NewAudio("music.mp3", "mp3")
	if err != nil {
		HandleError(err)
	}

	// Config: Game Seed
	game_seed, err := strconv.ParseInt(Configs["GAME_SEED"], 10, 64)
	if err != nil {
		HandleError(err)
	}

	// Config: Enemy Spawn Time
	enemy_spawn_time, err := strconv.ParseInt(Configs["ENEMY_SPAWN_TIME"], 10, 0)
	if err != nil {
		HandleError(err)
	}

	// Config: Maximum Enemies per Wave
	tmp, err := strconv.ParseInt(Configs["MAX_ENEMIES_PER_WAVE"], 10, 0)
	if err != nil {
		HandleError(err)
	}

	max_enemies_per_wave = int(tmp)

	g := &Game{
		// Utils
		random: rand.New(rand.NewSource(game_seed)),
		music:  music,
		save:   NewSave(),

		// Mechanics
		score:           NewScore(),
		enemySpawnTimer: NewTimer(time.Duration(enemy_spawn_time) * time.Second),

		// Entities
		player: NewPlayer(),

		// Flags
		hasSavedOnDeath: false,

		// Misc.
		oneSecondTimer: NewTimer(1000 * time.Millisecond),
	}

	// Attach the game UI
	g.ui = NewUi(g)

	// Trigger enemy spawner once on init
	g.enemySpawnTimer.TriggerNow()

	// Config: Music Volume
	music_volume, err := strconv.ParseFloat(Configs["MUSIC_VOLUME"], 64)
	if err != nil {
		HandleError(err)
	}

	// Play the music
	g.music.SetVolume(music_volume)
	g.music.Play()

	// Load the Save
	g.save.LoadSave(g, false)

	return g
}
