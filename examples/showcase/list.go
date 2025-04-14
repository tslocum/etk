//go:build example

package main

import (
	"fmt"
	"log"

	"codeberg.org/tslocum/etk"
)

func newListExample() (string, etk.Widget, etk.Widget) {
	const fontSize = 32
	onSelected := func(index int) (accept bool) {
		log.Printf("Selected item at index %d", index)
		return true
	}

	ff := etk.FontFace(etk.Style.TextFont, fontSize)
	m := ff.Metrics()
	list := etk.NewList(etk.Scale(int(m.HAscent+m.HDescent)), onSelected)

	for i := 0; i < 100; i++ {
		t := etk.NewText(fmt.Sprintf("Item #%d", i+1))
		t.SetVertical(etk.AlignCenter)
		t.SetFont(etk.Style.TextFont, fontSize)
		t.SetAutoResize(true)
		list.AddChildAt(t, 0, i)
	}

	list.SetSelectedItem(0, 0)

	return "list", list, list
}

const listLabel = "The tab navigation to the right is a List."
