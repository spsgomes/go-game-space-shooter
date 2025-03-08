package game

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/audio"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
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

type ButtonState int

const (
	ButtonStateDefault ButtonState = iota
	ButtonStateHover   ButtonState = iota
)

type Button struct {
	text      string
	tag       string
	position  *Vector
	collision *CollisionRect
	state     ButtonState
}

type Ui struct {
	game              *Game
	background        *Background
	mainMenuButtons   []Button
	pausedMenuButtons []Button
	deathMenuButtons  []Button
	forceCursorShape  ebiten.CursorShapeType
	font              *text.GoTextFace
	fontBytes         []byte
}

type Save struct {
	path     string
	filename string
	data     map[string]any
}

type GameState int

const (
	GameStateInitial GameState = iota
	GameStatePlaying GameState = iota
	GameStatePaused  GameState = iota
	GameStateDeath   GameState = iota
)

type DamageNumber struct {
	damage      float64
	x           float64
	y           float64
	effect      string
	ticksPassed int
}

type Game struct {
	// Utils
	random *rand.Rand
	music  *audio.Audio
	save   *Save
	state  GameState

	// Mechanics
	score            *Score
	ui               *Ui
	enemySpawnTimer  *Timer
	pickupSpawnTimer *Timer
	damageNumbers    []DamageNumber

	// Entities
	player      *Player
	enemies     []*Enemy
	projectiles []*Projectile
	pickups     []*Pickup

	// Flags
	hasSavedOnDeath bool

	// Counters
	currentWave int

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
	spriteName       string
	fireRate         float64
	velocity         float64
	damage           float64
	criticalChance   float64
	criticalModifier float64
	timer            *Timer
	audio            *audio.Audio
	hitAudio         *audio.Audio
}

type Player struct {
	character *Character
	attack    *Attack
	disabled  bool
}

type Enemy struct {
	character           *Character
	attack              *Attack
	enemyType           string
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
	hitAudio  *audio.Audio
	ownerTag  string
	owner     Character
	damage    float64
	critical  bool
	disabled  bool
}

type Pickup struct {
	position     *Vector
	collision    *CollisionRect
	sprite       *assets.Sprite
	audio        *audio.Audio
	effectName   string
	effectAmount float64
	disabled     bool
}
