//go:build example

package main

import (
	"code.rocket9labs.com/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
)

type game struct {
}

func newGame() *game {
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
