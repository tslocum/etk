//go:build example

package main

import (
	"fmt"
	"log"

	"code.rocket9labs.com/tslocum/etk"
)

func newGridExample() (etk.Widget, etk.Widget) {
	newButton := func(i int) *etk.Button {
		return etk.NewButton("Flexible", func() error {
			log.Printf("Clicked button %d", i)
			return nil
		})
	}

	newText := func(size int) *etk.Text {
		t := etk.NewText(fmt.Sprintf("Fixed (%dpx)", size))
		t.SetHorizontal(etk.AlignCenter)
		t.SetVertical(etk.AlignCenter)
		return t
	}

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

	return grid, nil
}
