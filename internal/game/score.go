package game

import (
	_ "embed"
)

func NewScore() *Score {
	return &Score{
		best:    0,
		current: 0,
	}
}

func (s *Score) GetScore() int64 {
	return s.current
}

func (s *Score) ResetScore() {
	s.current = 0
}

func (s *Score) GetHighScore() int64 {
	return s.best
}

func (s *Score) SetHighScore(highscore int64) {
	s.best = highscore
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
