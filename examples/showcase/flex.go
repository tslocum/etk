//go:build example

package main

import (
	"fmt"

	"code.rocket9labs.com/tslocum/etk"
)

func newFlexExample() (etk.Widget, etk.Widget) {
	newLabel := func(i int) *etk.Text {
		t := etk.NewText(fmt.Sprintf("Item #%d", i))
		t.SetPadding(etk.Scale(10))
		return t
	}

	l1 := newLabel(1)
	l2 := newLabel(2)
	l3 := newLabel(3)
	l4 := newLabel(4)
	l5 := newLabel(5)

	flexDemo := etk.NewFlex()
	flexDemo.SetChildSize(etk.Scale(300), etk.Scale(75))
	flexDemo.AddChild(l1, l2, l3, l4, l5)

	return flexDemo, nil
}
