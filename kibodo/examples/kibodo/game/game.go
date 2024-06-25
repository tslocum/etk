//go:build example

package game

import (
	"fmt"

	"code.rocket9labs.com/tslocum/etk/kibodo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type game struct {
	w, h int

	k *kibodo.Keyboard

	userInput []byte

	incomingInput []*kibodo.Input

	op *ebiten.DrawImageOptions

	buffer *ebiten.Image
}

var spinner = []byte(`-\|/`)

// NewDemoGame returns a new kibodo demo game.
func NewDemoGame() *game {
	k := kibodo.NewKeyboard()
	k.SetPassThroughPhysicalInput(true)
	k.SetKeys(kibodo.KeysQWERTY)

	g := &game{
		k: k,
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
	}

	go g.showKeyboard()

	return g
}

func (g *game) showKeyboard() {
	if g.k.Visible() {
		return
	}

	// Clear current input
	g.userInput = nil

	// Show keyboard
	g.k.Show()
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	s := ebiten.DeviceScaleFactor()
	outsideWidth, outsideHeight = int(float64(outsideWidth)*s), int(float64(outsideHeight)*s)
	if g.w == outsideWidth && g.h == outsideHeight {
		return outsideWidth, outsideHeight
	}

	g.w, g.h = outsideWidth, outsideHeight

	g.buffer = ebiten.NewImage(g.w, g.h)

	y := 200
	if g.h > g.w && (g.h-g.w) > 200 {
		y = g.h - g.w
	}
	g.k.SetRect(0, y, g.w, g.h-y)

	return outsideWidth, outsideHeight
}

func (g *game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.k.Visible() {
		g.showKeyboard()
	}

	g.incomingInput = g.k.AppendInput(g.incomingInput[:0])
	for _, input := range g.incomingInput {
		if input.Rune > 0 {
			g.userInput = append(g.userInput, []byte(string(input.Rune))...)
			continue
		}
		if input.Key == ebiten.KeyBackspace {
			s := string(g.userInput)
			if len(s) > 0 {
				g.userInput = []byte(s[:len(s)-1])
			}
			continue
		} else if input.Key == ebiten.KeyEnter {
			g.userInput = nil
			continue
		} else if input.Key < 0 {
			continue
		}
		g.userInput = append(g.userInput, []byte("<"+input.Key.String()+">")...)
	}

	return g.k.Update()
}

func (g *game) Draw(screen *ebiten.Image) {
	g.k.Draw(screen)

	g.buffer.Clear()
	ebitenutil.DebugPrint(g.buffer, fmt.Sprintf("FPS %0.0f\nTPS %0.0f\n\n%s", ebiten.ActualFPS(), ebiten.ActualTPS(), g.userInput))
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(3, 0)
	g.op.GeoM.Scale(2, 2)
	screen.DrawImage(g.buffer, g.op)
}
