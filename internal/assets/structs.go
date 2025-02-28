package assets

import "github.com/hajimehoshi/ebiten/v2"

type SpriteInfo struct {
	name     string // populated automatically
	path     string // defaults to "./"
	filename string
}

type Sprite struct {
	Image *ebiten.Image
	Info  *SpriteInfo
}
