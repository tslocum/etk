//go:build example

package main

import (
	"codeberg.org/tslocum/etk"
)

func newSelectExample() (etk.Widget, etk.Widget) {
	s := etk.NewSelect(etk.Scale(int(float64(etk.Style.TextSize)*1.5)), nil)
	s.AddOption("Option 1")
	s.AddOption("Option 2")
	s.AddOption("Option 3")
	s.AddOption("Option 4")

	selectList := s.Children()[0]

	f := etk.NewFrame()
	f.SetPositionChildren(true)
	f.SetMaxWidth(etk.Scale(200))
	f.SetMaxHeight(etk.Scale(50))
	f.AddChild(s)

	selectDemo := etk.NewGrid()
	selectDemo.SetColumnPadding(etk.Scale(50))
	selectDemo.SetRowPadding(etk.Scale(50))
	selectDemo.AddChildAt(f, 0, 0, 1, 1)

	outer := etk.NewFrame()
	outer.SetPositionChildren(true)
	outer.AddChild(selectDemo, etk.NewFrame(selectList))

	return outer, nil
}
