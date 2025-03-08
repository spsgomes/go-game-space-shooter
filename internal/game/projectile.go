package game

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/audio"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewProjectile(ownerTag string, owner *Character, spriteName string, x float64, y float64, initialAngle float64, velocity float64, damage float64, critical bool, hitAudio *audio.Audio) *Projectile {
	sprite, err := assets.NewSprite(spriteName)
	if err != nil {
		HandleError(err)
	}

	projectile := Projectile{
		position: &Vector{
			x:     x,
			y:     y,
			angle: initialAngle,
			scale: 1,
		},
		sprite: sprite,
		movement: &Movement{
			velocity: velocity,
		},
		direction: &MovementDirection{
			oDx:   0,
			oDy:   0,
			angle: initialAngle,
		},
		hitAudio: hitAudio,
		ownerTag: ownerTag,
		owner:    *owner,
		damage:   damage,
		critical: critical,
		disabled: false,
	}

	return &projectile
}

func (p *Projectile) Update(g *Game) {

	if p.disabled {
		return
	}

	p.updateMovement()
	p.checkCollisions(g)
}

func (p *Projectile) Draw(screen *ebiten.Image) {

	if p.disabled {
		return
	}

	op := &ebiten.DrawImageOptions{}

	p.sprite.Rotate(op, p.position.angle)
	p.sprite.Scale(op, p.position.scale)
	p.sprite.Translate(op, p.position.scale, p.position.x, p.position.y)

	screen.DrawImage(p.sprite.Image, op)

	op.GeoM.Reset()

	// Config: Draw Colission Rects
	if Configs["DRAW_COLLISION_RECTS"] == "1" {
		// Draw collision rectangle
		vector.StrokeRect(screen, float32(p.collision.x0), float32(p.collision.y0), float32(p.collision.x1-p.collision.x0), float32(p.collision.y1-p.collision.y0), 1.0, color.RGBA{0, 255, 0, 255}, true)
	}
}

func (p *Projectile) SetProjectileDirection(target *Vector) {
	p.direction.oDx, p.direction.oDy, _ = DistanceBetweenTwoPoints(p.position, target)
	p.position.angle = ((math.Atan2(p.direction.oDy, p.direction.oDx) * 180) / math.Pi) - 90
}

func (p *Projectile) IsOutOfBounds() bool {

	wsX, wsY := GetWindowSize()

	Dx := p.sprite.Image.Bounds().Dx()
	Dy := p.sprite.Image.Bounds().Dy()

	return p.position.x < -float64(Dx) || p.position.x > wsX+float64(Dx) || p.position.y < -float64(Dy) || p.position.y > wsY+float64(Dy)
}

func (p *Projectile) updateMovement() {

	// Set position based on directional movement
	p.position.x += p.direction.oDx * p.movement.velocity
	p.position.y += p.direction.oDy * p.movement.velocity

	// Update collision rectangle
	x0, y0, x1, y1 := GetObjectRectCoords(p.position, p.sprite, p.position.scale)
	p.collision = &CollisionRect{x0: x0 - 20, y0: y0 + 20, x1: x1 + 20, y1: y1 + 20}

}

func (p *Projectile) checkCollisions(g *Game) {

	if p.disabled {
		return
	}

	srcRect := image.Rect(p.collision.x0, p.collision.y0, p.collision.x1, p.collision.y1)

	// Check collisions with Player
	if p.ownerTag == "enemy" {
		if !g.player.disabled && srcRect.Overlaps(image.Rect(g.player.character.position.collision.x0, g.player.character.position.collision.y0, g.player.character.position.collision.x1, g.player.character.position.collision.y1)) {

			// Play the hit audio
			p.hitAudio.Play()

			// Remove from player's HP
			g.player.OffsetHp(-p.damage)

			// Disable the projectile
			p.disabled = true
		}

	} else if p.ownerTag == "player" {

		for _, enemy := range g.enemies {

			if !enemy.disabled && srcRect.Overlaps(image.Rect(enemy.character.position.collision.x0, enemy.character.position.collision.y0, enemy.character.position.collision.x1, enemy.character.position.collision.y1)) {

				// Play the hit audio
				p.hitAudio.Play()

				// Remove from enemy's HP
				enemy.OffsetHp(-p.damage)

				// Disable the projectile
				p.disabled = true

				// If enemy was killed, add to score
				if enemy.disabled {
					// ? DEBUG
					// fmt.Println("enemy killed, awards points:", enemy.worthPoints)

					g.score.AddScore(enemy.worthPoints)
				}
			}
		}
	}
}
