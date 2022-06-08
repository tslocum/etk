//go:build example
// +build example

package main

import (
	"code.rocketnine.space/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("etk showcase")

	g := newGame()

	b1 := etk.NewButton("Button 1", nil)
	b2 := etk.NewButton("Button 2", nil)

	topFlex := etk.NewFlex()
	topFlex.AddChild(b1, b2)

	b3 := etk.NewButton("Button 3", nil)
	b4 := etk.NewButton("Button 4", nil)
	b5 := etk.NewButton("Button 5", nil)

	bottomFlex := etk.NewFlex()
	bottomFlex.AddChild(b3, b4, b5)

	rootFlex := etk.NewFlex()
	rootFlex.SetVertical(true)
	rootFlex.AddChild(topFlex, bottomFlex)

	etk.SetRoot(rootFlex)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
