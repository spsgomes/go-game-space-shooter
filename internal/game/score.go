package game

import (
	"bytes"
	_ "embed"
	"go-game-space-shooter/internal/assets"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func NewScore() *Score {

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

	return &Score{
		best:    0,
		current: 0,
		font:    font,
	}
}

func (s *Score) Draw(screen *ebiten.Image) {

	wsX, _ := GetWindowSize()

	s.font.Size = 80
	op := &text.DrawOptions{}
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignCenter
	op.GeoM.Translate(float64(wsX)/2.0, 0.00)

	str := strconv.FormatInt(s.GetScore(), 10)
	text.Draw(screen, str, s.font, op)
	_, textH := text.Measure(str, s.font, op.LineSpacing)

	s.font.Size = 30
	op = &text.DrawOptions{}
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/2/255.0)
	op.PrimaryAlign = text.AlignCenter
	op.GeoM.Translate(float64(wsX)/2.0, textH)

	text.Draw(screen, "best: "+strconv.FormatInt(s.GetHighScore(), 10), s.font, op)
}

func (s *Score) GetScore() int64 {
	return s.current
}

func (s *Score) GetHighScore() int64 {
	return s.best
}

func (s *Score) AddScore(add int64) {
	s.current += add

	if s.IsHighScore() {
		s.best = s.current
	}
}

func (s *Score) IsHighScore() bool {
	return s.current > s.best
}
