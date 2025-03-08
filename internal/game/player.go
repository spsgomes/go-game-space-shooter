package game

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/audio"
	"image/color"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewPlayer() *Player {
	sprite, err := assets.NewSprite("player")
	if err != nil {
		HandleError(err)
	}

	attackAudio, err := audio.NewAudio("laser.wav", "wav")
	if err != nil {
		HandleError(err)
	}

	hitAudio, err := audio.NewAudio("damage2.wav", "wav")
	if err != nil {
		HandleError(err)
	}

	wsX, wsY := GetWindowSize()

	player := Player{
		character: &Character{
			position: &CharacterVector{
				vector: &Vector{
					x: float64(wsX) / 2.0,
					y: float64(wsY) / 2.0,
				},
				angle: 0.0,
				scale: 0.6,
			},
			movement: &Movement{
				velocity:        10.0,
				turningVelocity: 5.0,
			},
			sprite: sprite,
			hp: &Health{
				max:     100.0,
				current: 100.0,
			},
		},
		attack: &Attack{
			spriteName:       "laser_blue",
			fireRate:         6.0,
			velocity:         10.0,
			damage:           5.0,
			criticalChance:   5.0,
			criticalModifier: 2.0,
			audio:            attackAudio,
			hitAudio:         hitAudio,
		},
	}

	// Apply configs
	player.applyConfigs()

	// Create attack timer
	player.attack.timer = NewTimer(time.Millisecond * time.Duration(1.0/player.attack.fireRate*1000))

	return &player
}

func (p *Player) Update(g *Game) {

	if p.disabled {
		return
	}

	if g.state == GameStatePlaying {
		p.updateMovement()
		p.updateAttack(g)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {

	if p.disabled {
		return
	}

	op := &ebiten.DrawImageOptions{}

	p.character.sprite.Rotate(op, p.character.position.angle)
	p.character.sprite.Scale(op, p.character.position.scale)
	p.character.sprite.Translate(op, p.character.position.scale, p.character.position.vector.x, p.character.position.vector.y)

	screen.DrawImage(p.character.sprite.Image, op)

	op.GeoM.Reset()

	// Config: Draw Colission Rects
	if Configs["DRAW_COLLISION_RECTS"] == "1" && p.character.position.collision != nil {
		// Draw collision rectangle
		vector.StrokeRect(screen, float32(p.character.position.collision.x0), float32(p.character.position.collision.y0), float32(p.character.position.collision.x1-p.character.position.collision.x0), float32(p.character.position.collision.y1-p.character.position.collision.y0), 1.0, color.RGBA{0, 0, 255, 255}, true)
	}
}

func (p *Player) OffsetHp(offset float64) {

	tmp := p.character.hp.current

	p.character.hp.current += offset

	if p.character.hp.current <= 0 {
		p.character.hp.current = 0
		p.disabled = true

	} else if p.character.hp.current > p.character.hp.max {
		p.character.hp.current = p.character.hp.max
	}

	if tmp > p.character.hp.current {
		// ? DEBUG
		// fmt.Println("Player hit:", tmp, "->", p.character.hp.current)
	}
}

func (p *Player) applyConfigs() {

	configs := p.getConfigs()

	// Apply configs: Attack Volume
	if configs["attack_volume"] >= 0.00 {
		p.attack.audio.SetVolume(configs["attack_volume"])
	}

	// Apply configs: Player Scale
	if configs["player_scale"] > 0.00 {
		p.character.position.scale = configs["player_scale"]
	}

	// Apply configs: Player HP
	if configs["player_hp"] > 0.00 {
		p.character.hp.max = configs["player_hp"]
		p.character.hp.current = configs["player_hp"]
	}

	// Apply config: Player Fire Rate
	if configs["player_fire_rate"] > 0.00 {
		p.attack.fireRate = configs["player_fire_rate"]
	}

	// Apply config: Player Projectile Speed
	if configs["player_projectile_speed"] >= 0.00 {
		p.attack.velocity = configs["player_projectile_speed"]
	}

	// Apply config: Player Projectile Speed
	if configs["player_projectile_damage"] > 0.00 {
		p.attack.damage = configs["player_projectile_damage"]
	}
}

func (p *Player) getConfigs() map[string]float64 {
	var val float64
	var err error

	configs := make(map[string]float64)

	// Config: Attack Volume
	val, err = strconv.ParseFloat(Configs["ATTACK_VOLUME"], 64)
	if err == nil {
		configs["attack_volume"] = val
	}

	// Config: Player Scale
	val, err = strconv.ParseFloat(Configs["PLAYER_SCALE"], 64)
	if err == nil {
		configs["player_scale"] = val
	}

	// Config: Player HP
	val, err = strconv.ParseFloat(Configs["PLAYER_HP"], 64)
	if err == nil {
		configs["player_hp"] = val
	}

	// Config: Player Fire Rate
	val, err = strconv.ParseFloat(Configs["PLAYER_FIRE_RATE"], 64)
	if err == nil {
		configs["player_fire_rate"] = val
	}

	// Config: Player Projectile Speed
	val, err = strconv.ParseFloat(Configs["PLAYER_PROJECTILE_SPEED"], 64)
	if err == nil {
		configs["player_projectile_speed"] = val
	}

	// Config: Player Projectile Speed
	val, err = strconv.ParseFloat(Configs["PLAYER_PROJECTILE_DAMAGE"], 64)
	if err == nil {
		configs["player_projectile_damage"] = val
	}

	return configs
}

func (p *Player) updateMovement() {

	// Flag to check if the player is turning
	var turning int8 = 0

	// Player Controls: Up
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.character.position.vector.y -= p.character.movement.velocity
	}

	// Player Controls: Down
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.character.position.vector.y += p.character.movement.velocity
	}

	// Player Controls: Left
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.character.position.vector.x -= p.character.movement.velocity
		turning = -1
	}

	// Player Controls: Right
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.character.position.vector.x += p.character.movement.velocity
		turning = 1
	}

	// Rotate back to the original position
	if turning == 0 {
		if p.character.position.angle > 0 {
			turning = -1
		} else if p.character.position.angle < 0 {
			turning = 1
		}
	}

	p.character.position.angle = p.getPositionAngle()

	// Guarantee the player remains within bounds
	p.character.position.vector.x, p.character.position.vector.y = CheckWithinBounds(
		p.character.position.vector.x,
		p.character.position.vector.y,
		float64(p.character.sprite.Image.Bounds().Dx()),
		float64(p.character.sprite.Image.Bounds().Dy()),
		p.character.position.scale,
	)

	// Update collision rectangle
	x0, y0, x1, y1 := GetSpriteRectCoords(p.character.position.vector, p.character.sprite, p.character.position.scale)
	p.character.position.collision = &CollisionRect{x0: x0 - 5, y0: y0 - 25, x1: x1 + 2, y1: y1 + 5}

	// ? DEBUG
	// fmt.Println(p.character.position.collision.x0, p.character.position.collision.y0, p.character.position.collision.x1, p.character.position.collision.y1, float64(float64(p.character.sprite.Image.Bounds().Dx())*p.character.position.scale))
}

func (p *Player) updateAttack(g *Game) {

	p.attack.timer.Update()
	if p.attack.timer.IsReady() {

		// Player Controls: Shoot
		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			p.attack.timer.Reset()

			// Find middle of character  position vector
			projectileX := p.character.position.vector.x
			projectileY := p.character.position.vector.y - ((float64(p.character.sprite.Image.Bounds().Dy())) / 2.0)

			// Attack values
			attackDamage := p.attack.damage
			attackCritical := false

			// Calculate critical
			if g.random.Float64()*100.0 <= p.attack.criticalChance {
				attackCritical = true
				attackDamage *= p.attack.criticalModifier
			}

			// Create a new projectile
			projectile := NewProjectile("player", p.character, p.attack.spriteName, projectileX, projectileY, p.character.position.angle, p.attack.velocity, attackDamage, attackCritical, p.attack.hitAudio)
			projectile.SetProjectileDirection(GetCursorVector())

			// Add to projectile list
			g.projectiles = append(g.projectiles, projectile)

			// Play the attack audio
			p.attack.audio.Play()
		}
	}
}

func (p *Player) getPositionAngle() float64 {
	dx, dy, _ := DistanceBetweenTwoPoints(p.character.position.vector, GetCursorVector())
	return ((math.Atan2(dy, dx) * 180) / math.Pi) + 90
}
