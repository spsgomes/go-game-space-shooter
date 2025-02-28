package main

import (
	"go-game-space-shooter/internal/assets"
	"go-game-space-shooter/internal/game"
	"path"
	"runtime"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joho/godotenv"
)

const window_title = "Space Shooter in Go! - by Simão Gomes at simaogomes.com"

func main() {
	var Configs map[string]string
	Configs, err := godotenv.Read(path.Join(getRuntimeDirectory(), "..", "..", "configs/configs.env"))
	if err != nil {
		panic(err)
	}

	// Config: Window Width
	window_width, err := strconv.Atoi(Configs["WINDOW_WIDTH"])
	if err != nil {
		panic(err)
	}

	// Config: Window Height
	window_height, err := strconv.Atoi(Configs["WINDOW_HEIGHT"])
	if err != nil {
		panic(err)
	}

	window_icon, _ := assets.GetWindowIconImages()

	ebiten.SetWindowSize(window_width, window_height)
	ebiten.SetWindowTitle(window_title)
	ebiten.SetWindowIcon(window_icon)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)

	// Config: Fullscreen Enabled
	if Configs["FULLSCREEN_ENABLED"] == "1" {
		ebiten.SetFullscreen(true)
	}

	// Initialize a new game
	g := game.NewGame(Configs)

	// Run the game
	err = ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}

func getRuntimeDirectory() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(file)
	}

	return ""
}
