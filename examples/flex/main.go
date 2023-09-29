//go:build example

package main

import (
	"fmt"
	"log"

	"code.rocketnine.space/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("etk flex example")

	newButton := func(i int) *etk.Button {
		return etk.NewButton(fmt.Sprintf("Button %d", i), func() error {
			log.Printf("Pressed button %d", i)
			return nil
		})
	}

	g := newGame()

	b1 := newButton(1)
	b2 := newButton(2)

	topFlex := etk.NewFlex()
	topFlex.AddChild(b1, b2)

	b3 := newButton(3)
	b4 := newButton(4)
	b5 := newButton(5)

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
