//go:build example

package main

import (
	"code.rocket9labs.com/tslocum/etk"
)

func newWindowExample() (etk.Widget, etk.Widget) {
	text := etk.NewText(windowLabel)
	text.SetPadding(etk.Scale(10))
	text.SetFollow(false)

	return text, nil
}

const windowLabel = `This widget showcase utilizes a Window.

The Window is currently showing this text as its content.

The tab list to the right allows navigating to the other widgets within the window.`
