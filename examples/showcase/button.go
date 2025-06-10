//go:build example

package main

import (
	"fmt"

	"codeberg.org/tslocum/etk"
)

func newButtonExample() (string, etk.Widget, etk.Widget) {
	var button *etk.Button
	var i int
	onSelect := func() error {
		i++
		label := "Clicked 1 time"
		if i > 1 {
			label = fmt.Sprintf("Clicked %d times", i)
		}
		button.SetText(label)
		return nil
	}
	button = etk.NewButton("Click here", onSelect)

	frame := etk.NewFrame()
	frame.SetPositionChildren(true)
	frame.SetMaxHeight(etk.Scale(100))
	frame.SetMaxWidth(etk.Scale(300))
	frame.AddChild(button)

	buttonGrid := etk.NewGrid()
	buttonGrid.SetColumnPadding(etk.Scale(50))
	buttonGrid.SetRowPadding(etk.Scale(50))
	buttonGrid.AddChildAt(frame, 0, 0, 1, 1)

	return "button", buttonGrid, nil
}
