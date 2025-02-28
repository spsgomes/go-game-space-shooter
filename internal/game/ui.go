package game

import (
	"bytes"
	_ "embed"
	"go-game-space-shooter/internal/assets"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func NewUi(game *Game) *Ui {

	fontTrainOneRegularTTF, err := assets.LoadFont("", "TrainOne-Regular.ttf")
	if err != nil {
		HandleError(err)
	}

	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontTrainOneRegularTTF))
	if err != nil {
		HandleError(err)
	}

	font := &text.GoTextFace{
		Source: fontSource,
		Size:   80,
	}

	return &Ui{
		game: game,
		background: &Background{
			filename: "background",
			ticker:   0.0,
			velocity: 0.5,
			oDx:      1.0,
			oDy:      1.0,
		},
		font:      font,
		fontBytes: fontTrainOneRegularTTF,
	}
}

func (u *Ui) Draw(screen *ebiten.Image) {

	// Player is dead
	if u.game.player.disabled {
		u.drawDeathScreen(screen)

	} else {
		u.drawHpBar(screen)
	}

	// Draw score
	u.drawScore(screen)
}

func (u *Ui) DrawBackground(screen *ebiten.Image) {
	sprite, err := assets.NewSprite("background")
	if err != nil {
		HandleError(err)
	}

	wsX, wsY := GetWindowSize()

	op := &ebiten.DrawImageOptions{}

	for x := -1; x < int(math.Ceil(wsX/background_size)); x++ {
		for y := -1; y < int(math.Ceil(wsY/background_size)); y++ {
			op.GeoM.Translate(float64(x*background_size)+u.background.ticker*u.background.oDx*u.background.velocity, float64(y*background_size)+u.background.ticker*u.background.oDy*u.background.velocity)
			screen.DrawImage(sprite.Image, op)
			op.GeoM.Reset()
		}
	}

	u.background.ticker++

	if u.background.ticker*u.background.oDx*u.background.velocity > background_size {
		u.background.ticker = 0
	}
}

func (u *Ui) drawDeathScreen(screen *ebiten.Image) {
	wsX, wsY := GetWindowSize()

	op := &text.DrawOptions{}
	op.LineSpacing = 20
	u.font.Size = 100
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 0/255.0, 0/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignCenter

	str := "YOU ARE DEAD"

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/2.0-textH/2.0)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 0/255.0, 0/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = "nice try though"

	_, textH = text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/2.0+textH+op.LineSpacing)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()
}

func (u *Ui) drawHpBar(screen *ebiten.Image) {
	_, wsY := GetWindowSize()

	op := &text.DrawOptions{}
	u.font.Size = 20
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignStart

	strs := []string{
		strconv.FormatFloat(u.game.player.character.hp.current, 'f', -1, 64),
		strconv.FormatFloat(u.game.player.character.hp.max, 'f', -1, 64),
	}
	str := "HP: " + strings.Join(strs, "/")

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(window_padding, wsY-window_padding-textH)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()
}

func (u *Ui) drawScore(screen *ebiten.Image) {

	wsX, _ := GetWindowSize()

	op := &text.DrawOptions{}
	u.font.Size = 100
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignCenter

	str := strconv.FormatInt(u.game.score.GetScore(), 10)

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(float64(wsX)/2.0, 0)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = strconv.FormatInt(u.game.score.GetHighScore(), 10)

	op.GeoM.Translate(float64(wsX)/2.0, textH-20)
	text.Draw(screen, "best: "+str, u.font, op)
	op.GeoM.Reset()
}
