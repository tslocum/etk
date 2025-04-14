//go:build example

package main

import (
	"codeberg.org/tslocum/etk"
)

func newWindowExample() (string, etk.Widget, etk.Widget) {
	text := etk.NewText(windowLabel)
	text.SetPadding(etk.Scale(10))

	return "window", text, nil
}

const windowLabel = `This widget showcase utilizes a Window.

The Window is currently showing this text as its content.

The tab list to the right allows navigating to the other widgets within the window.`
