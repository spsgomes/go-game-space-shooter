package game

import (
	"bytes"
	"go-game-space-shooter/internal/assets"
	"image/color"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	WINDOW_PADDING  = 40
	BACKGROUND_SIZE = 256
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

func (u *Ui) Update() error {
	// Pause/Unpause
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if u.game.state == GameStatePlaying {
			u.game.state = GameStatePaused
		} else if u.game.state == GameStatePaused {
			u.game.state = GameStatePlaying
		}
	}

	// Start game
	if u.game.state == GameStateInitial && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		u.game.state = GameStatePlaying
	}

	// Exit game
	if (u.game.state == GameStateInitial || u.game.state == GameStateDeath) && inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	// Restart game
	if u.game.state == GameStateDeath && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		u.game.Restart()
		u.game.state = GameStatePlaying
	}

	return nil
}

func (u *Ui) Draw(screen *ebiten.Image) {

	switch u.game.state {
	case GameStateInitial:
		u.drawMainMenu(screen)
	case GameStateDeath:
		u.drawDeathScreen(screen)
	case GameStatePaused:
		u.drawPauseScreen(screen)
		u.drawEnemiesHpBar(screen)
		u.drawPlayerHpBar(screen)
		u.drawDamageNumbers(screen)
		u.drawPlayerStats(screen)

	case GameStatePlaying:
		u.drawEnemiesHpBar(screen)
		u.drawPlayerHpBar(screen)
		u.drawDamageNumbers(screen)
		u.drawPlayerStats(screen)
	}

	if u.game.state != GameStateInitial {
		u.drawScore(screen)
	}
}

func (u *Ui) DrawBackground(screen *ebiten.Image) {
	sprite, err := assets.NewSprite("background")
	if err != nil {
		HandleError(err)
	}

	wsX, wsY := GetWindowSize()

	op := &ebiten.DrawImageOptions{}

	for x := -1; x < int(math.Ceil(wsX/BACKGROUND_SIZE)); x++ {
		for y := -1; y < int(math.Ceil(wsY/BACKGROUND_SIZE)); y++ {
			op.GeoM.Translate(float64(x*BACKGROUND_SIZE)+u.background.ticker*u.background.oDx*u.background.velocity, float64(y*BACKGROUND_SIZE)+u.background.ticker*u.background.oDy*u.background.velocity)
			screen.DrawImage(sprite.Image, op)
			op.GeoM.Reset()
		}
	}

	if slices.Contains([]GameState{GameStateInitial, GameStatePlaying, GameStateDeath}, u.game.state) {
		u.background.ticker++

		if u.background.ticker*u.background.oDx*u.background.velocity > BACKGROUND_SIZE {
			u.background.ticker = 0
		}
	}
}

func (u *Ui) drawMainMenu(screen *ebiten.Image) {
	wsX, wsY := GetWindowSize()

	op := &text.DrawOptions{}
	op.LineSpacing = 30
	u.font.Size = 100
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/255.0)
	op.PrimaryAlign = text.AlignCenter

	str := "Space Shooter!"

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/4.0-textH/2.0)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 255/255.0, 255/255.0, 255/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = "press space to start\nor escape to quit"

	_, textH = text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/4.0+textH+op.LineSpacing)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	op.LineSpacing = 20
	u.font.Size = 16
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 255/255.0, 255/255.0, 255/255.0)
	op.PrimaryAlign = text.AlignStart

	str =
		`Controls
W: Up
S: Down
A: Left
D: Right
Space/Left Click: Shoot
Escape: Pause/Unpause`

	_, textH = text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(WINDOW_PADDING, wsY-WINDOW_PADDING-textH)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()
}

func (u *Ui) drawPauseScreen(screen *ebiten.Image) {
	wsX, wsY := GetWindowSize()

	vector.DrawFilledRect(screen, 0, 0, float32(wsX), float32(wsY), color.RGBA{0, 0, 0, uint8(math.Floor(255 * 0.5))}, true)

	op := &text.DrawOptions{}
	op.LineSpacing = 30
	u.font.Size = 100
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255*0.5/255.0)
	op.PrimaryAlign = text.AlignCenter

	str := "PAUSED"

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/2.0-textH/2.0)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 255/255.0, 255/255.0, 255*0.5/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = "press escape to continue"

	_, textH = text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/2.0+textH+op.LineSpacing)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

}

func (u *Ui) drawDeathScreen(screen *ebiten.Image) {
	wsX, wsY := GetWindowSize()

	op := &text.DrawOptions{}
	op.LineSpacing = 0
	u.font.Size = 100
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 0/255.0, 0/255.0, 255/255.0)
	op.PrimaryAlign = text.AlignCenter

	str := "YOU ARE DEAD"

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/2.0-textH/2.0-20)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	op.LineSpacing = 20
	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 0/255.0, 0/255.0, 255/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = "nice try though"

	_, textH2 := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY/2.0+textH2+op.LineSpacing-20)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	op.LineSpacing = 30
	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(255/255.0, 255/255.0, 255/255.0, 255/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = "press space to restart\nor escape to quit"

	// _, textH3 := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(wsX/2.0, wsY-wsY/4.0-op.LineSpacing)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()
}

func (u *Ui) drawPlayerHpBar(screen *ebiten.Image) {
	_, wsY := GetWindowSize()

	op := &text.DrawOptions{}
	u.font.Size = 20
	op.ColorScale.Reset()

	if u.game.player.character.hp.current <= u.game.player.character.hp.max*0.25 {
		// Low Health
		op.ColorScale.Scale(255/255.0, 0/255.0, 0/255.0, 255/255.0)

	} else if u.game.player.character.hp.current <= u.game.player.character.hp.max*0.5 {
		// Median Health
		op.ColorScale.Scale(255/255.0, 255/255.0, 0/255.0, 255/255.0)

	} else {
		// Default
		op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255/255.0)
	}

	op.PrimaryAlign = text.AlignStart

	strs := []string{
		strconv.FormatFloat(u.game.player.character.hp.current, 'f', -1, 64),
		strconv.FormatFloat(u.game.player.character.hp.max, 'f', -1, 64),
	}
	str := "HP: " + strings.Join(strs, "/")

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(WINDOW_PADDING, wsY-WINDOW_PADDING-textH)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()
}

func (u *Ui) drawPlayerStats(screen *ebiten.Image) {
	wsX, wsY := GetWindowSize()

	op := &text.DrawOptions{}
	op.LineSpacing = 20
	u.font.Size = 16
	op.ColorScale.Reset()

	op.ColorScale.Scale(255/255.0, 255/255.0, 255/255.0, 255/255.0)

	op.PrimaryAlign = text.AlignStart

	strs := []string{
		TrimTrailingZeros(strconv.FormatFloat(u.game.player.attack.damage, 'f', 2, 64)),
		TrimTrailingZeros(strconv.FormatFloat(u.game.player.attack.criticalChance, 'f', 2, 64)) + "%",
		"x" + TrimTrailingZeros(strconv.FormatFloat(u.game.player.attack.criticalModifier, 'f', 2, 64)),
	}
	str := strings.Join(strs, "\n")

	textW, textH := text.Measure(str, u.font, op.LineSpacing)

	if textW < 70 {
		textW = 70
	}

	op.GeoM.Translate(wsX-WINDOW_PADDING-textW, wsY-WINDOW_PADDING-textH)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	op.PrimaryAlign = text.AlignEnd

	strs = []string{
		"Damage: ",
		"Crit. Chance: ",
		"Crit. Modifier: ",
	}
	str = strings.Join(strs, "\n")

	op.GeoM.Translate(wsX-WINDOW_PADDING-textW, wsY-WINDOW_PADDING-textH)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()
}

func (u *Ui) drawEnemiesHpBar(screen *ebiten.Image) {

	// Loop through enemies
	if len(u.game.enemies) > 0 {
		for _, enemy := range u.game.enemies {
			if enemy.character.position.collision != nil && !enemy.disabled {

				posX := float32(enemy.character.position.collision.x0)
				posY := float32(enemy.character.position.collision.y0 - 20)
				posW := float32(enemy.character.position.collision.x1 - enemy.character.position.collision.x0)
				posH := float32(8.0)

				vector.StrokeRect(screen, posX, posY, posW, posH, 1.0, color.RGBA{255, 255, 255, 1}, true)

				posW *= float32(enemy.character.hp.current*100/enemy.character.hp.max) / 100

				vector.DrawFilledRect(screen, posX+1, posY+1, posW-1, posH-1, color.RGBA{255, 0, 0, 1}, true)
			}
		}
	}
}

func (u *Ui) drawScore(screen *ebiten.Image) {

	wsX, _ := GetWindowSize()

	op := &text.DrawOptions{}
	u.font.Size = 100
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255*0.5/255.0)
	op.PrimaryAlign = text.AlignCenter

	str := strconv.FormatInt(u.game.score.GetScore(), 10)

	_, textH := text.Measure(str, u.font, op.LineSpacing)

	op.GeoM.Translate(float64(wsX)/2.0, 0)
	text.Draw(screen, str, u.font, op)
	op.GeoM.Reset()

	u.font.Size = 24
	op.ColorScale.Reset()
	op.ColorScale.Scale(10/255.0, 191/255.0, 245/255.0, 255*0.5/255.0)
	op.PrimaryAlign = text.AlignCenter

	str = strconv.FormatInt(u.game.score.GetHighScore(), 10)

	op.GeoM.Translate(float64(wsX)/2.0, textH-20)
	text.Draw(screen, "best: "+str, u.font, op)
	op.GeoM.Reset()
}

// Draw Damage Numbers
func (u *Ui) drawDamageNumbers(screen *ebiten.Image) {

	const MAX_TICKS = 200

	if len(u.game.damageNumbers) > 0 {
		var newDamageNumbers []DamageNumber

		for _, damageNumber := range u.game.damageNumbers {

			op := &text.DrawOptions{}
			op.LineSpacing = 16
			u.font.Size = 16
			op.ColorScale.Reset()

			switch damageNumber.effect {
			case "golden":
				op.ColorScale.Scale(255/255.0, 223/255.0, 0/255.0, 255.0)
			case "hurt":
				op.ColorScale.Scale(255/255.0, 0/255.0, 0/255.0, 255.0)
			default:
				op.ColorScale.Scale(255/255.0, 255/255.0, 255/255.0, 255.0)
			}

			alpha := float32(MAX_TICKS-damageNumber.ticksPassed) * 1 / 100
			if alpha < 0 {
				alpha = 0
			} else if alpha > 1 {
				alpha = 1
			}

			op.ColorScale.ScaleAlpha(alpha)

			op.PrimaryAlign = text.AlignCenter

			str := TrimTrailingZeros(strconv.FormatFloat(damageNumber.damage, 'f', 2, 64))

			op.GeoM.Translate(damageNumber.x, damageNumber.y-float64(damageNumber.ticksPassed))
			text.Draw(screen, str, u.font, op)
			op.GeoM.Reset()

			// Live for only 100 ticks
			if damageNumber.ticksPassed <= MAX_TICKS {
				damageNumber.ticksPassed += 1

				newDamageNumbers = append(newDamageNumbers, damageNumber)
			}
		}

		u.game.damageNumbers = newDamageNumbers
	}

}
