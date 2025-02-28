package game

import (
	"bytes"
	_ "embed"
	"go-game-space-shooter/internal/assets"

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
		best:      0,
		current:   0,
		font:      font,
		fontBytes: fontTrainOneRegularTTF,
	}
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
