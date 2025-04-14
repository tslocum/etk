//go:build example

package main

import (
	"fmt"
	"log"

	"codeberg.org/tslocum/etk"
)

func newGridExample() (string, etk.Widget, etk.Widget) {
	newButton := func(i int) *etk.Button {
		return etk.NewButton("Flexible", func() error {
			log.Printf("Clicked button %d", i)
			return nil
		})
	}
	button1 := newButton(1)
	button2 := newButton(2)
	button3 := newButton(3)

	newText := func(size int) *etk.Text {
		t := etk.NewText(fmt.Sprintf("Fixed (%dpx)", size))
		t.SetHorizontal(etk.AlignCenter)
		t.SetVertical(etk.AlignCenter)
		return t
	}
	text1 := newText(75)
	text2 := newText(75)
	text3 := newText(150)

	grid := etk.NewGrid()
	grid.SetRowSizes(75, -1, 75, -1, 150)

	// First row.
	grid.AddChildAt(text1, 0, 0, 2, 1)

	// Second row.
	grid.AddChildAt(button1, 0, 1, 2, 1)

	// Third row.
	grid.AddChildAt(text2, 0, 2, 2, 1)

	// Fourth row.
	grid.AddChildAt(button2, 0, 3, 1, 1)
	grid.AddChildAt(button3, 1, 3, 1, 1)

	// Fifth row.
	grid.AddChildAt(text3, 0, 4, 2, 1)

	return "grid", grid, nil
}
