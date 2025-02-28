package game

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/audio"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Score struct {
	best    int64
	current int64
}

type Background struct {
	filename string
	ticker   float64
	velocity float64
	oDx      float64
	oDy      float64
}

type Ui struct {
	game       *Game
	background *Background
	font       *text.GoTextFace
	fontBytes  []byte
}

type Game struct {
	// Utils
	random *rand.Rand
	music  *audio.Audio

	// Mechanics
	score           *Score
	ui              *Ui
	enemySpawnTimer *Timer

	// Entities
	player      *Player
	enemies     []*Enemy
	projectiles []*Projectile

	// Misc.
	oneSecondTimer *Timer
}

type Timer struct {
	currentTicks int
	targetTicks  int
}

type Vector struct {
	x     float64
	y     float64
	angle float64
	scale float64
}

type CollisionRect struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

type Movement struct {
	velocity        float64
	turningVelocity float64
}

type MovementDirection struct {
	oDx   float64
	oDy   float64
	angle float64
}

type Health struct {
	max     float64
	current float64
}

type CharacterVector struct {
	vector    *Vector
	collision *CollisionRect
	angle     float64
	scale     float64
}

type Character struct {
	position *CharacterVector
	movement *Movement
	sprite   *assets.Sprite
	hp       *Health
}

type Attack struct {
	spriteName string
	fireRate   float64
	velocity   float64
	damage     float64
	timer      *Timer
	audio      *audio.Audio
}

type Player struct {
	character *Character
	attack    *Attack
	disabled  bool
}

type Enemy struct {
	character           *Character
	attack              *Attack
	worthPoints         int64
	minLengthFromPlayer float64
	isRunningAway       bool
	isStopped           bool
	disabled            bool
}

type Projectile struct {
	position  *Vector
	collision *CollisionRect
	sprite    *assets.Sprite
	movement  *Movement
	direction *MovementDirection
	ownerTag  string
	owner     Character
	damage    float64
	disabled  bool
}
