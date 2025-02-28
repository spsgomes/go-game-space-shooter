package assets

import (
	"embed"
	"errors"
	"image"
	_ "image/png"
	"math"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.png *.ttf
var assets embed.FS

var SpriteMap = map[string]SpriteInfo{
	"background": {
		filename: "background.png",
	},
	"player": {
		filename: "player.png",
	},
	"enemy": {
		filename: "enemy.png",
	},
	"laser_blue": {
		filename: "laser_blue.png",
	},
	"laser_red": {
		filename: "laser_red.png",
	},
}

var cache = map[string]*Sprite{}

func NewSprite(name string) (*Sprite, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	// Fetch sprite from cache if it exists
	spriteCacheValue, ok := cache[name]
	if ok {
		return spriteCacheValue, nil
	}

	spriteInfo, ok := SpriteMap[name]
	if !ok {
		return nil, errors.New(name + " was not found in the sprite map")
	}

	spriteInfo.name = name
	if spriteInfo.path == "" {
		spriteInfo.path = ""
	}

	loadedImage, err := loadImage(spriteInfo.path, spriteInfo.filename)
	if err != nil {
		return nil, err
	}

	sprite := &Sprite{
		Image: loadedImage,
		Info:  &spriteInfo,
	}

	cache[name] = sprite

	return sprite, nil
}

func (s *Sprite) Rotate(op *ebiten.DrawImageOptions, angle float64) {
	s.centerSpriteXY(op)
	op.GeoM.Rotate(angle * math.Pi / 180.0)
	s.resetSpriteXY(op, 1)
}

func (s *Sprite) Scale(op *ebiten.DrawImageOptions, scale float64) {
	s.centerSpriteXY(op)
	op.GeoM.Scale(scale, scale)
	s.resetSpriteXY(op, 1)
}

func (s *Sprite) Translate(op *ebiten.DrawImageOptions, scale float64, x float64, y float64) {
	s.centerSpriteXY(op)
	op.GeoM.Translate(x-float64(s.Image.Bounds().Dx())*scale/2, y-float64(s.Image.Bounds().Dx())*scale/2)
	s.resetSpriteXY(op, scale)
}

func LoadFont(path string, filename string) ([]byte, error) {
	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	f, err := assets.ReadFile(filepath.Join(path, filename))

	if err != nil {
		return nil, errors.New("cannot open font with filename \"" + path + filename + "\"")
	}

	return f, nil
}

func loadImage(path string, filename string) (*ebiten.Image, error) {

	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	f, err := assets.Open(filepath.Join(path, filename))

	if err != nil {
		return nil, errors.New("cannot open image with filename \"" + path + filename + "\"")
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.New("cannot decode image with filename \"" + path + filename + "\"")
	}

	return ebiten.NewImageFromImage(img), nil
}

func GetWindowIconImages() ([]image.Image, error) {

	var images []image.Image

	f, err := assets.Open("window_icon.png")

	if err != nil {
		return nil, errors.New("cannot open window icon image")
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, errors.New("cannot decode window icon image")
	}

	images = append(images, img)

	return images, nil
}

func (s *Sprite) centerSpriteXY(op *ebiten.DrawImageOptions) {
	// Center the sprite for GeoM operations
	spriteWidth := s.Image.Bounds().Dx()
	spriteWidthHalf := float64(spriteWidth / 2)
	spriteHeight := s.Image.Bounds().Dy()
	spriteHeightHalf := float64(spriteHeight / 2)

	op.GeoM.Translate(-spriteWidthHalf, -spriteHeightHalf)
}
func (s *Sprite) resetSpriteXY(op *ebiten.DrawImageOptions, scale float64) {
	// Center the sprite for GeoM operations
	spriteWidth := s.Image.Bounds().Dx()
	spriteWidthHalf := float64(spriteWidth) * scale / 2
	spriteHeight := s.Image.Bounds().Dy()
	spriteHeightHalf := float64(spriteHeight) * scale / 2

	op.GeoM.Translate(spriteWidthHalf, spriteHeightHalf)
}
