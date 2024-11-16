//go:build example

package game

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"strings"
	"sync"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const initialText = `
Welcome to the メッセージ (messeji) text widgets demo.
This is a TextField, which can be used to display text.
Below is an InputField, which accepts keyboard input.
<Tab> to cycle horizontal alignment.
<Enter> to append input text to buffer.
<Ctrl+Tab> to toggle word wrap.
<Ctrl+Enter> to toggle multi-line input.
`

var (
	fontSource *text.GoTextFaceSource
	fontSize   = 24
	fontMutex  *sync.Mutex
)

func init() {
	var err error
	fontSource, err = text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
}

type game struct {
	w, h int

	buffer *messeji.TextField

	input *messeji.InputField

	singleLine bool

	op *ebiten.DrawImageOptions

	spinnerIndex int

	horizontal messeji.Alignment
}

// NewDemoGame returns a new messeji demo game.
func NewDemoGame() *game {
	g := &game{
		buffer: messeji.NewTextField(fontSource, fontSize, fontMutex),
		input:  messeji.NewInputField(fontSource, fontSize, fontMutex),
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
	}

	g.buffer.SetText(strings.TrimSpace(initialText))
	g.buffer.SetPadding(7)

	g.input.SetHandleKeyboard(true)
	g.input.SetSelectedFunc(func() (accept bool) {
		log.Printf("Input: %s", g.input.Text())

		g.buffer.Write([]byte(fmt.Sprintf("\nInput: %s", g.input.Text())))

		return true
	})
	g.input.SetPadding(7)

	return g
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth == g.w && outsideHeight == g.h {
		return outsideWidth, outsideHeight
	}

	padding := 10

	w, h := outsideWidth-padding*2, g.input.LineHeight()*3+g.input.Padding()*2
	if h > outsideHeight-padding {
		h = outsideHeight - padding
	}

	x, y := outsideWidth/2-w/2, outsideHeight-h-padding

	g.buffer.SetRect(image.Rect(padding, padding, outsideWidth-padding, y-padding))

	g.input.SetRect(image.Rect(x, y, x+w, y+h))

	g.w, g.h = outsideWidth, outsideHeight
	return outsideWidth, outsideHeight
}

func (g *game) Update() error {
	if (inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter)) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		g.singleLine = !g.singleLine
		g.input.SetSingleLine(g.singleLine)
		return nil
	} else if (inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyKPEnter)) && !g.input.Visible() {
		g.input.SetVisible(true)
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		wrap := g.buffer.WordWrap()
		g.buffer.SetWordWrap(!wrap)
		g.input.SetWordWrap(!wrap)
		return nil
	} else if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.horizontal++
		if g.horizontal > messeji.AlignEnd {
			g.horizontal = messeji.AlignStart
		}
		g.buffer.SetHorizontal(g.horizontal)
		return nil
	}

	err := g.buffer.Update()
	if err != nil {
		return fmt.Errorf("failed to update buffer: %s", err)
	}

	err = g.input.Update()
	if err != nil {
		return fmt.Errorf("failed to update input field: %s", err)
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Draw display field.
	g.buffer.Draw(screen)

	// Draw input field.
	g.input.Draw(screen)
}
