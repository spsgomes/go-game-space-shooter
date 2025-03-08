package game

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/audio"
	"image"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func NewPickup(effectName string, effectAmount float64, spriteName string, x float64, y float64, audio *audio.Audio) *Pickup {
	sprite, err := assets.NewSprite(spriteName)
	if err != nil {
		HandleError(err)
	}

	pickup := Pickup{
		position: &Vector{
			x:     x,
			y:     y,
			angle: 0,
			scale: 1,
		},
		sprite:       sprite,
		audio:        audio,
		effectName:   effectName,
		effectAmount: effectAmount,
	}

	// Update collision rectangle
	x0, y0, x1, y1 := GetSpriteRectCoords(pickup.position, pickup.sprite, pickup.position.scale)
	pickup.collision = &CollisionRect{x0: x0, y0: y0, x1: x1, y1: y1}

	return &pickup
}

func SpawnPickups(random *rand.Rand, pickups []*Pickup, max int) []*Pickup {

	var spawn_types = make([]map[string]any, 0)

	// Spawn type: health
	spawn_types = append(spawn_types, map[string]any{
		"typeName":    "health",
		"amount":      20.0,
		"sprite":      "pill_blue",
		"audio":       "pickup.wav",
		"audioType":   "wav",
		"audioVolume": 0.5,
	})

	// Spawn type: damage
	spawn_types = append(spawn_types, map[string]any{
		"typeName":    "damage",
		"amount":      1.1,
		"sprite":      "bolt_bronze",
		"audio":       "pickup.wav",
		"audioType":   "wav",
		"audioVolume": 0.5,
	})

	// Spawn type: critical_modifier
	spawn_types = append(spawn_types, map[string]any{
		"typeName":    "critical_modifier",
		"amount":      0.5,
		"sprite":      "blue_box_bolt",
		"audio":       "pickup.wav",
		"audioType":   "wav",
		"audioVolume": 0.5,
	})

	// Spawn type: critical_chance
	spawn_types = append(spawn_types, map[string]any{
		"typeName":    "critical_chance",
		"amount":      5.0,
		"sprite":      "blue_box_star",
		"audio":       "pickup.wav",
		"audioType":   "wav",
		"audioVolume": 0.5,
	})

	const OFFSET float64 = 200.0

	qty := random.Intn(max) + 1
	wsX, wsY := GetWindowSize()

	for range qty {

		// Get random type to spawn
		spawn_type := spawn_types[random.Intn(len(spawn_types))]

		// Get random position for spawn
		eX := random.Intn(int(wsX-OFFSET*2)) + int(OFFSET) // Generate an integer number between [OFFSET] and [wsX-OFFSET]
		eY := random.Intn(int(wsY-OFFSET*2)) + int(OFFSET) // Generate an integer number between [OFFSET] and [wsY-OFFSET]

		spawnAudio, err := audio.NewAudio(spawn_type["audio"].(string), spawn_type["audioType"].(string))
		if err != nil {
			HandleError(err)
		}
		spawnAudio.SetVolume(spawn_type["audioVolume"].(float64))

		pickups = append(pickups, NewPickup(
			spawn_type["typeName"].(string),
			spawn_type["amount"].(float64),
			spawn_type["sprite"].(string),
			float64(eX),
			float64(eY),
			spawnAudio,
		))
	}

	return pickups
}

func (p *Pickup) Update(g *Game) {

	if p.disabled {
		return
	}

	p.checkCollisions(g)
}

func (p *Pickup) Draw(screen *ebiten.Image) {

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
		vector.StrokeRect(screen, float32(p.collision.x0), float32(p.collision.y0), float32(p.collision.x1-p.collision.x0), float32(p.collision.y1-p.collision.y0), 1.0, color.RGBA{255, 255, 0, 255}, true)
	}
}

func (p *Pickup) checkCollisions(g *Game) {

	if p.disabled {
		return
	}

	srcRect := image.Rect(p.collision.x0, p.collision.y0, p.collision.x1, p.collision.y1)

	// Check collisions with Player
	if !g.player.disabled && srcRect.Overlaps(image.Rect(g.player.character.position.collision.x0, g.player.character.position.collision.y0, g.player.character.position.collision.x1, g.player.character.position.collision.y1)) {

		// Play the audio
		p.audio.Play()

		// Decide what to do based on effect type
		switch p.effectName {
		case "health":
			// Add to player's HP
			g.player.OffsetHp(p.effectAmount)

		case "damage":
			// Add to player's HP
			g.player.attack.damage *= p.effectAmount

		case "critical_chance":
			// Add to player's attack critical chance
			g.player.attack.criticalChance += p.effectAmount

		case "critical_modifier":
			// Add to player's attack critical modifier
			g.player.attack.criticalModifier += p.effectAmount
		}

		// Disable the pickup
		p.disabled = true
	}
}
