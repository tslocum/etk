//go:build example

package main

import (
	"bytes"

	"codeberg.org/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type game struct {
}

func newGame() *game {
	etk.Style.TextFont = defaultFont()
	text.CacheGlyphs(loremIpsum, etk.FontFace(etk.Style.TextFont, etk.Scale(etk.Style.TextSize)))
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
		panic(err)
	}
}

func defaultFont() *text.GoTextFaceSource {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
	return source
}

func monoFont() *text.GoTextFaceSource {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		panic(err)
	}
	return source
}
