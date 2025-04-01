//go:build example

package main

import (
	"image/color"

	"codeberg.org/tslocum/etk"
)

func newCheckboxExample() (etk.Widget, etk.Widget) {
	var chk *etk.Checkbox
	var label *etk.Button
	onSelectChk := func() error {
		if chk.Selected() {
			label.SetText("Checked")
		} else {
			label.SetText("Unchecked")
		}
		return nil
	}
	onSelectLabel := func() error {
		chk.SetSelected(!chk.Selected())
		onSelectChk()
		return nil
	}
	chk = etk.NewCheckbox(onSelectChk)
	label = etk.NewButton("Unchecked", onSelectLabel)
	label.SetHorizontal(etk.AlignStart)
	label.SetVertical(etk.AlignCenter)
	label.SetForeground(color.RGBA{255, 255, 255, 255})
	label.SetBackground(color.RGBA{0, 0, 0, 0})
	label.SetBorderColors(color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0})

	grid := etk.NewGrid()
	grid.SetColumnSizes(etk.Scale(50), -1)
	grid.AddChildAt(chk, 0, 0, 1, 1)
	grid.AddChildAt(label, 1, 0, 1, 1)

	f := etk.NewFrame()
	f.SetPositionChildren(true)
	f.SetMaxHeight(etk.Scale(50))
	f.AddChild(grid)

	btnDemo := etk.NewGrid()
	btnDemo.SetColumnPadding(etk.Scale(50))
	btnDemo.SetRowPadding(etk.Scale(50))
	btnDemo.AddChildAt(f, 0, 0, 1, 1)

	return btnDemo, nil
}
