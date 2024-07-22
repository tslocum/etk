//go:build example

package main

import (
	"log"

	"code.rocket9labs.com/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type game struct {
}

func newGame() *game {
	etk.Style.TextFont = defaultFont()
	g := &game{}
	return g
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	etk.Layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
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

func defaultFont() *sfnt.Font {
	f, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
