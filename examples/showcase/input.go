//go:build example

package main

import (
	"code.rocket9labs.com/tslocum/etk"
)

func newInputExample() (etk.Widget, etk.Widget) {
	buffer := etk.NewText("Press enter to append input below to this buffer.")
	onselected := func(text string) (handled bool) {
		buffer.Write([]byte("\nInput: " + text))
		return true
	}
	input := etk.NewInput("", onselected)
	input.SetPadding(etk.Scale(10))
	inputDemo := etk.NewFlex()
	inputDemo.SetVertical(true)
	inputDemo.AddChild(buffer, input)

	return inputDemo, input
}
