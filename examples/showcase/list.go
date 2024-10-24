//go:build example

package main

import (
	"code.rocket9labs.com/tslocum/etk"
)

func newListExample() (etk.Widget, etk.Widget) {
	text := etk.NewText(listLabel)
	text.SetPadding(etk.Scale(10))
	text.SetFollow(false)

	return text, nil
}

const listLabel = "The tab navigation to the right is a List."
