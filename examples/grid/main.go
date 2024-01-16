//go:build example

package main

import (
	"fmt"
	"log"

	"code.rocket9labs.com/tslocum/etk"
	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("etk grid example")

	newButton := func(i int) *etk.Button {
		return etk.NewButton(fmt.Sprintf("Button %d", i), func() error {
			log.Printf("Pressed button %d", i)
			return nil
		})
	}

	newText := func(size int) *etk.Text {
		t := etk.NewText(fmt.Sprintf("%dpx Text", size))
		t.SetHorizontal(messeji.AlignCenter)
		t.SetVertical(messeji.AlignCenter)
		return t
	}

	g := newGame()

	b1 := newButton(1)
	b2 := newButton(2)
	b3 := newButton(3)

	grid := etk.NewGrid()
	grid.SetRowSizes(75, -1, 75, -1, 150)

	// First row.
	grid.AddChildAt(newText(75), 0, 0, 2, 1)

	// Second row.
	grid.AddChildAt(b1, 0, 1, 2, 1)

	// Third row.
	grid.AddChildAt(newText(75), 0, 2, 2, 1)

	// Fourth row.
	grid.AddChildAt(b2, 0, 3, 1, 1)
	grid.AddChildAt(b3, 1, 3, 1, 1)

	// Fifth row.
	grid.AddChildAt(newText(150), 0, 4, 2, 1)

	etk.SetRoot(grid)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
