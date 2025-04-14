//go:build example

package main

import (
	"codeberg.org/tslocum/etk"
)

func newSelectExample() (string, etk.Widget, etk.Widget) {
	s := etk.NewSelect(etk.Scale(int(float64(etk.Style.TextSize)*1.5)), nil)
	s.AddOption("Option 1")
	s.AddOption("Option 2")
	s.AddOption("Option 3")
	s.AddOption("Option 4")

	frame := etk.NewFrame()
	frame.SetPositionChildren(true)
	frame.SetMaxWidth(etk.Scale(200))
	frame.SetMaxHeight(etk.Scale(50))
	frame.AddChild(s)

	selectDemo := etk.NewGrid()
	selectDemo.SetColumnPadding(etk.Scale(50))
	selectDemo.SetRowPadding(etk.Scale(50))
	selectDemo.AddChildAt(frame, 0, 0, 1, 1)

	selectList := s.Children()[0]

	outer := etk.NewFrame()
	outer.SetPositionChildren(true)
	outer.AddChild(selectDemo, etk.NewFrame(selectList))

	return "select", outer, nil
}
