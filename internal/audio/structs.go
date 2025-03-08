package audio

import "github.com/hajimehoshi/ebiten/v2/audio"

type Audio struct {
	context *audio.Context
	Player  *audio.Player
	volume  float64
}
