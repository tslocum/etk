//go:build example

package main

import (
	"fmt"

	"codeberg.org/tslocum/etk"
)

func newFlexExample() (string, etk.Widget, etk.Widget) {
	newLabel := func(i int) *etk.Text {
		t := etk.NewText(fmt.Sprintf("Item #%d", i))
		t.SetPadding(etk.Scale(10))
		return t
	}
	label1 := newLabel(1)
	label2 := newLabel(2)
	label3 := newLabel(3)
	label4 := newLabel(4)
	label5 := newLabel(5)

	flex := etk.NewFlex()
	flex.SetChildSize(etk.Scale(300), etk.Scale(75))
	flex.AddChild(label1, label2, label3, label4, label5)

	return "flex", flex, nil
}
