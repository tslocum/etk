//go:build example

package main

import (
	"codeberg.org/tslocum/etk"
)

func newInputExample() (string, etk.Widget, etk.Widget) {
	buffer := etk.NewText("Press enter to append input below to this buffer.")
	buffer.SetFollow(true)
	buffer.SetPadding(etk.Scale(10))

	onselected := func(text string) (handled bool) {
		buffer.Write([]byte("\nInput: " + text))
		return true
	}
	input := etk.NewInput("", onselected)
	input.SetPadding(etk.Scale(10))

	inputFlex := etk.NewFlex()
	inputFlex.SetVertical(true)
	inputFlex.AddChild(buffer, input)

	return "input", inputFlex, input
}
