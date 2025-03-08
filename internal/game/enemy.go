package game

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/audio"
	"image/color"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewEnemy(enemyType string, spriteName string, x float64, y float64, angle float64) *Enemy {
	sprite, err := assets.NewSprite(spriteName)
	if err != nil {
		HandleError(err)
	}

	hitAudio, err := audio.NewAudio("damage1.mp3", "mp3")
	if err != nil {
		HandleError(err)
	}

	enemy := Enemy{
		character: &Character{
			position: &CharacterVector{
				vector: &Vector{
					x: float64(x),
					y: float64(y),
				},
				angle: angle,
				scale: 0.6,
			},
			movement: &Movement{
				velocity: 1.0,
			},
			sprite: sprite,
			hp: &Health{
				max:     10.0,
				current: 10.0,
			},
		},
		attack: &Attack{
			spriteName:       "laser_red",
			fireRate:         0.5,
			velocity:         5.0,
			damage:           10.0,
			criticalChance:   0.0,
			criticalModifier: 0.0,
			hitAudio:         hitAudio,
		},
		enemyType:           enemyType,
		worthPoints:         10,
		minLengthFromPlayer: 200.0,
		isRunningAway:       false,
		isStopped:           false,
		disabled:            false,
	}

	// Apply configs
	enemy.applyConfigs()

	// Enemy type
	switch enemy.enemyType {
	case "tank":
		enemy.character.position.scale = 1.0
		enemy.character.hp.max *= 3.0
		enemy.character.hp.current = enemy.character.hp.max
	case "boss":
		enemy.character.position.scale = 1.0
		enemy.character.hp.max *= 20.0
		enemy.character.hp.current = enemy.character.hp.max
		enemy.attack.fireRate *= 6.0
		enemy.attack.damage *= 3.0
	default:
		// Pass
	}

	// Create attack timer
	enemy.attack.timer = NewTimer(time.Millisecond * time.Duration(1.0/enemy.attack.fireRate*1000))

	return &enemy
}

func (e *Enemy) Update(g *Game, p *Player) {

	if e.disabled {
		return
	}

	e.updateMovement(p)
	e.updateAttack(g)
}

func (e *Enemy) Draw(screen *ebiten.Image) {

	if e.disabled {
		return
	}

	op := &ebiten.DrawImageOptions{}

	e.character.sprite.Rotate(op, e.character.position.angle)
	e.character.sprite.Scale(op, e.character.position.scale)
	e.character.sprite.Translate(op, e.character.position.scale, e.character.position.vector.x, e.character.position.vector.y)

	screen.DrawImage(e.character.sprite.Image, op)

	op.GeoM.Reset()

	// Config: Draw Colission Rects
	if Configs["DRAW_COLLISION_RECTS"] == "1" && e.character.position.collision != nil {
		// Draw collision rectangle
		vector.StrokeRect(screen, float32(e.character.position.collision.x0), float32(e.character.position.collision.y0), float32(e.character.position.collision.x1-e.character.position.collision.x0), float32(e.character.position.collision.y1-e.character.position.collision.y0), 1.0, color.RGBA{255, 0, 0, 255}, true)
	}
}

func SpawnEnemies(random *rand.Rand, enemies []*Enemy, currentWave int, max int) []*Enemy {

	const OFFSET_Y int = 200

	// Enemy Spawner: Boss!
	if currentWave > 0 && currentWave%10 == 0 {
		eX, eY := GetRandomSpawnPosition(random, OFFSET_Y)
		enemies = append(enemies, NewEnemy("boss", "boss", eX, eY, 0))

		// Don't spawn other enemies in boss encounter
		return enemies
	}

	// Enemy Spawner: Basic
	for range random.Intn(max) + 1 {
		eX, eY := GetRandomSpawnPosition(random, OFFSET_Y)
		enemies = append(enemies, NewEnemy("basic", "enemy", eX, eY, 0))
	}

	// Enemy Spawner: Tank (50% chance after wave 5)
	if currentWave >= 5 && random.Float64()*100.0 <= 50 {
		eX, eY := GetRandomSpawnPosition(random, OFFSET_Y)
		enemies = append(enemies, NewEnemy("tank", "enemy", eX, eY, 0))
	}

	return enemies
}

func (e *Enemy) OffsetHp(offset float64) {

	tmp := e.character.hp.current

	e.character.hp.current += offset

	if e.character.hp.current <= 0 {
		e.character.hp.current = 0
		e.disabled = true

	} else if e.character.hp.current > e.character.hp.max {
		e.character.hp.current = e.character.hp.max
	}

	if tmp > e.character.hp.current {
		// ? DEBUG
		// fmt.Println("Enemy hit:", tmp, "->", e.character.hp.current)
	}
}

func (e *Enemy) applyConfigs() {

	configs := e.getConfigs()

	// Apply configs: Enemy Scale
	if configs["enemy_scale"] > 0.00 {
		e.character.position.scale = configs["enemy_scale"]
	}

	// Apply configs: Enemy HP
	if configs["enemy_hp"] > 0.00 {
		e.character.hp.max = configs["enemy_hp"]
		e.character.hp.current = configs["enemy_hp"]
	}

	// Apply config: Enemy Fire Rate
	if configs["enemy_fire_rate"] > 0.00 {
		e.attack.fireRate = configs["enemy_fire_rate"]
	}

	// Apply config: Enemy Projectile Speed
	if configs["enemy_projectile_speed"] >= 0.00 {
		e.attack.velocity = configs["enemy_projectile_speed"]
	}

	// Apply config: Enemy Projectile Speed
	if configs["enemy_projectile_damage"] > 0.00 {
		e.attack.damage = configs["enemy_projectile_damage"]
	}

	// Apply config: Enemy Point Worth
	if configs["enemy_point_worth"] > 0.00 {
		e.worthPoints = int64(configs["enemy_point_worth"])
	}
}

func (e *Enemy) getConfigs() map[string]float64 {
	var val float64
	var err error

	configs := make(map[string]float64)

	// Config: Enemy Scale
	val, err = strconv.ParseFloat(Configs["ENEMY_SCALE"], 64)
	if err == nil {
		configs["enemy_scale"] = val
	}

	// Config: Enemy HP
	val, err = strconv.ParseFloat(Configs["ENEMY_HP"], 64)
	if err == nil {
		configs["enemy_hp"] = val
	}

	// Config: Enemy Fire Rate
	val, err = strconv.ParseFloat(Configs["ENEMY_FIRE_RATE"], 64)
	if err == nil {
		configs["enemy_fire_rate"] = val
	}

	// Config: Enemy Projectile Speed
	val, err = strconv.ParseFloat(Configs["ENEMY_PROJECTILE_SPEED"], 64)
	if err == nil {
		configs["enemy_projectile_speed"] = val
	}

	// Config: Enemy Projectile Speed
	val, err = strconv.ParseFloat(Configs["ENEMY_PROJECTILE_DAMAGE"], 64)
	if err == nil {
		configs["enemy_projectile_damage"] = val
	}

	// Config: Enemy Point Worth
	val, err = strconv.ParseFloat(Configs["ENEMY_POINT_WORTH"], 64)
	if err == nil {
		configs["enemy_point_worth"] = val
	}

	return configs
}

func (e *Enemy) updateMovement(p *Player) {

	dx, dy, length := DistanceBetweenTwoPoints(e.character.position.vector, p.character.position.vector)

	velocityModifier := 0.0

	// Player is dead, enemy will go back to home base
	if p.disabled {
		dx *= -1
		dy *= -1
		e.isRunningAway = true
		e.isStopped = false

	} else {

		// If the enemy is within the minLengthFromPlayer, move back by inverting it's trajectory
		if length < e.minLengthFromPlayer {
			dx *= -1
			dy *= -1
			e.isRunningAway = true
			e.isStopped = false

			// If it's between minLengthFromPlayer and minLengthFromPlayer+100, it stops and still shoots
		} else if (e.isRunningAway || e.isStopped) && length < e.minLengthFromPlayer+100 {
			e.isRunningAway = false
			e.isStopped = true

			// Gets closer to the player
		} else {
			e.isRunningAway = false
			e.isStopped = false
		}
	}

	// if the enemy is running away, increase it's velocity
	if e.isRunningAway {
		velocityModifier = 1.0
	}

	if !e.isStopped {
		e.character.position.vector.x += dx * (e.character.movement.velocity + velocityModifier)
		e.character.position.vector.y += dy * (e.character.movement.velocity + velocityModifier)
	}

	e.character.position.angle = ((math.Atan2(dy, dx) * 180) / math.Pi) - 90

	// Update collision rectangle
	x0, y0, x1, y1 := GetObjectRectCoords(e.character.position.vector, e.character.sprite, e.character.position.scale)
	e.character.position.collision = &CollisionRect{x0: x0 - 10, y0: y0 - 10, x1: x1 + 15, y1: y1 + 10}

	// ? DEBUG
	// fmt.Println(e.character.position.collision.x0, e.character.position.collision.y0, e.character.position.collision.x1, e.character.position.collision.y1, float64(float64(e.character.sprite.Image.Bounds().Dx())*e.character.position.scale))
}

func (e *Enemy) updateAttack(g *Game) {

	if g.player.disabled {
		return
	}

	e.attack.timer.Update()

	// If the enemy is running away, reset the timer so it doesn't shoot, and to wait for the next timer target
	if e.isRunningAway {
		e.attack.timer.Reset()
	}

	if e.attack.timer.IsReady() {
		e.attack.timer.Reset()

		// Find middle of character position vector
		projectileX := e.character.position.vector.x
		projectileY := e.character.position.vector.y - ((float64(e.character.sprite.Image.Bounds().Dy())) / 2.0)

		// Create a new projectile
		projectile := NewProjectile("enemy", e.character, e.attack.spriteName, projectileX, projectileY, e.character.position.angle, e.attack.velocity, e.attack.damage, false, e.attack.hitAudio)
		projectile.SetProjectileDirection(g.player.character.position.vector)

		// Add to projectile list
		g.projectiles = append(g.projectiles, projectile)
	}
}
