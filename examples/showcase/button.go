//go:build example

package main

import (
	"fmt"

	"code.rocket9labs.com/tslocum/etk"
)

func newButtonExample() (etk.Widget, etk.Widget) {
	var btn *etk.Button
	var clicked int
	onClick := func() error {
		clicked++
		label := "Clicked 1 time"
		if clicked > 1 {
			label = fmt.Sprintf("Clicked %d times", clicked)
		}
		btn.SetText(label)
		return nil
	}
	btn = etk.NewButton("Click here", onClick)

	f := etk.NewFrame()
	f.SetPositionChildren(true)
	f.SetMaxHeight(etk.Scale(100))
	f.SetMaxWidth(etk.Scale(300))
	f.AddChild(btn)

	btnDemo := etk.NewGrid()
	btnDemo.SetColumnPadding(etk.Scale(50))
	btnDemo.SetRowPadding(etk.Scale(50))
	btnDemo.AddChildAt(f, 0, 0, 1, 1)

	return btnDemo, nil
}
