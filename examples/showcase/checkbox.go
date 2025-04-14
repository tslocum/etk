//go:build example

package main

import (
	"image/color"

	"codeberg.org/tslocum/etk"
)

func newCheckboxExample() (string, etk.Widget, etk.Widget) {
	var checkbox *etk.Checkbox
	var label *etk.Button
	onSelectCheckbox := func() error {
		if checkbox.Selected() {
			label.SetText("Checked")
		} else {
			label.SetText("Unchecked")
		}
		return nil
	}
	checkbox = etk.NewCheckbox(onSelectCheckbox)

	onSelectLabel := func() error {
		checkbox.SetSelected(!checkbox.Selected())
		onSelectCheckbox()
		return nil
	}
	label = etk.NewButton("Unchecked", onSelectLabel)
	label.SetHorizontal(etk.AlignStart)
	label.SetVertical(etk.AlignCenter)
	label.SetForeground(color.RGBA{255, 255, 255, 255})
	label.SetBackground(color.RGBA{0, 0, 0, 0})
	label.SetBorderColors(color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0})

	grid := etk.NewGrid()
	grid.SetColumnSizes(etk.Scale(50), -1)
	grid.AddChildAt(checkbox, 0, 0, 1, 1)
	grid.AddChildAt(label, 1, 0, 1, 1)

	frame := etk.NewFrame()
	frame.SetPositionChildren(true)
	frame.SetMaxHeight(etk.Scale(50))
	frame.AddChild(grid)

	checkboxGrid := etk.NewGrid()
	checkboxGrid.SetColumnPadding(etk.Scale(50))
	checkboxGrid.SetRowPadding(etk.Scale(50))
	checkboxGrid.AddChildAt(frame, 0, 0, 1, 1)

	return "checkbox", checkboxGrid, nil
}
