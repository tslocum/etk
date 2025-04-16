package etk_test

import (
	"bytes"
	"log"

	"codeberg.org/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type game struct{}

func newGame() *game {
	// Load font.
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	etk.Style.TextFont = source

	g := &game{}
	return g
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return etk.Layout(outsideWidth, outsideHeight)
}

func (g *game) Update() error {
	return etk.Update()
}

func (g *game) Draw(screen *ebiten.Image) {
	err := etk.Draw(screen)
	if err != nil {
		log.Fatal(err)
	}
}

// A minimal example of how to use etk.
func Example() {
	// Initialize game.
	g := newGame()

	// Create text widget.
	t := etk.NewText("Hello, world!")

	// Set text widget as root widget.
	etk.SetRoot(t)

	// Run game.
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
