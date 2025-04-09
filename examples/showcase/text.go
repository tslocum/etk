//go:build example

package main

import (
	"codeberg.org/tslocum/etk"
)

func newTextExample() (etk.Widget, etk.Widget) {
	text := etk.NewText(loremIpsum)
	text.SetPadding(etk.Scale(10))

	return text, nil
}
