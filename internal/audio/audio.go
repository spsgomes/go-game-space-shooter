package audio

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"path/filepath"

	ebitenAudio "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

//go:embed *.mp3 *.wav
var audioFiles embed.FS

const sampleRate = 48000

var context = ebitenAudio.NewContext(sampleRate)

func NewAudio(filename string, filetype string) (*Audio, error) {

	file, err := loadAudioFile("", filename)
	if err != nil {
		return nil, err
	}

	if filetype == "" {
		return nil, errors.New("filetype cannot be empty")
	}

	r := bytes.NewReader(file)

	var player *ebitenAudio.Player

	var d audioStream
	switch filetype {
	case "mp3":
		d, err = mp3.DecodeF32(r)
		if err != nil {
			return nil, err
		}

	case "wav":
		d, err = wav.DecodeF32(r)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("filetype " + filetype + " not supported")
	}

	player, err = context.NewPlayerF32(d)
	if err != nil {
		return nil, err
	}

	newAudio := &Audio{
		context: context,
		player:  player,
		volume:  1,
	}

	return newAudio, nil
}

func (a *Audio) GetVolume() float64 {
	return a.volume
}

func (a *Audio) SetVolume(volume float64) {
	a.player.SetVolume(volume)
	a.volume = volume
}

func (a *Audio) Play() {
	a.player.Rewind()
	a.player.Play()
}

func (a *Audio) Continue() {
	if !a.player.IsPlaying() {
		a.player.Play()
	}
}

func (a *Audio) Pause() {
	if a.player.IsPlaying() {
		a.player.Pause()
	}
}

func loadAudioFile(path string, filename string) ([]byte, error) {

	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	f, err := audioFiles.ReadFile(filepath.Join(path, filename))

	if err != nil {
		return nil, errors.New("cannot open audio with filename \"" + path + filename + "\"")
	}

	return f, nil
}
