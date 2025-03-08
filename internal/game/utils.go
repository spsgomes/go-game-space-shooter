package game

import (
	"go-game-space-shooter/internal/assets"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func HandleError(error any) {
	panic(error)
}

func GetWindowSize() (float64, float64) {

	var windowWidth, windowHeight int

	if ebiten.IsFullscreen() {
		windowWidth, windowHeight = ebiten.Monitor().Size()
	} else {
		windowWidth, windowHeight = ebiten.WindowSize()
	}
	return float64(windowWidth), float64(windowHeight)

}

func CheckWithinBounds(x float64, y float64, w float64, h float64, s float64) (float64, float64) {

	wsX, wsY := GetWindowSize()

	if x < w*s/2+WINDOW_PADDING {
		x = w*s/2 + WINDOW_PADDING

	} else if x > float64(wsX)-(w*s/2)-WINDOW_PADDING {
		x = float64(wsX) - (w * s / 2) - WINDOW_PADDING
	}

	if y < h*s/2+WINDOW_PADDING {
		y = h*s/2 + WINDOW_PADDING

	} else if y > float64(wsY)-(h*s/2)-WINDOW_PADDING {
		y = float64(wsY) - (h * s / 2) - WINDOW_PADDING
	}

	return x, y
}

func DistanceBetweenTwoPoints(p1 *Vector, p2 *Vector) (dx float64, dy float64, length float64) {

	// Difference vector
	dx = p2.x - p1.x
	dy = p2.y - p1.y

	// Normalize (direction vector; a direction vector has a length of 1)
	length = math.Sqrt((dx * dx) + (dy * dy))
	if length > 0 {
		dx /= length
		dy /= length
	}

	return dx, dy, length
}

func GetObjectRectCoords(vector *Vector, sprite *assets.Sprite, scale float64) (x0, y0, x1, y1 int) {

	if scale == 0.0 {
		scale = 1.0
	}

	// Get source rectangle to check collisions with
	x0 = int(vector.x - (float64(sprite.Image.Bounds().Dx())*scale)/2.0)
	y0 = int(vector.y - (float64(sprite.Image.Bounds().Dy())*scale)/2.0)

	x1 = int(float64(sprite.Image.Bounds().Dx()) * scale)
	y1 = int(float64(sprite.Image.Bounds().Dy()) * scale)

	return x0, y0, x0 + x1, y0 + y1
}

func GetCursorVector() *Vector {
	cursorPosX, cursorPosY := ebiten.CursorPosition()
	return &Vector{
		x: float64(cursorPosX),
		y: float64(cursorPosY),
	}
}

func TrimTrailingZeros(str string) string {
	return strings.TrimRight(strings.TrimRight(str, "0"), ".")
}
