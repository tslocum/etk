//go:build example

package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"net/http"
	_ "net/http/pprof"

	"code.rocket9labs.com/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var debugAddress string
	flag.StringVar(&debugAddress, "debug", "", "serve debug information on address")
	flag.Parse()

	if debugAddress != "" {
		go func() {
			err := http.ListenAndServe(debugAddress, nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	ebiten.SetWindowTitle("etk widget showcase")

	g := newGame()

	w := etk.NewWindow()

	// Button.
	{
		var btn *etk.Button
		var clicked int
		onClick := func() error {
			clicked++
			label := "Clicked 1 time"
			if clicked > 1 {
				label = fmt.Sprintf("Clicked %d times", clicked)
			}
			btn.SetText(label)
			return nil
		}
		btn = etk.NewButton("Click here", onClick)

		f := etk.NewFrame()
		f.SetPositionChildren(true)
		f.SetMaxHeight(etk.Scale(100))
		f.SetMaxWidth(etk.Scale(300))
		f.AddChild(btn)

		btnDemo := etk.NewGrid()
		btnDemo.SetColumnPadding(etk.Scale(50))
		btnDemo.SetRowPadding(etk.Scale(50))
		btnDemo.AddChildAt(f, 0, 0, 1, 1)

		w.AddChildWithLabel(btnDemo, nil, "Button")
	}

	// Checkbox.
	{
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
		chk.SetBackground(color.RGBA{255, 255, 255, 255})
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

		w.AddChildWithLabel(btnDemo, nil, "Checkbox")
	}

	// Flex.
	{
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

		w.AddChildWithLabel(flexDemo, nil, "Flex")
	}

	// Input.
	buffer := etk.NewText("Press enter to append input below to this buffer.")
	onselected := func(text string) (handled bool) {
		buffer.Write([]byte("\nInput: " + text))
		return true
	}
	input := etk.NewInput("", onselected)
	input.SetPrefix(">")
	{
		inputDemo := etk.NewFlex()
		inputDemo.SetVertical(true)
		inputDemo.AddChild(buffer, input)

		w.AddChildWithLabel(inputDemo, input, "Input")
	}

	w.Show(0)

	etk.SetRoot(w)
	etk.SetFocus(input)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
