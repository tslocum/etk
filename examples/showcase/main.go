//go:build example

package main

import (
	"flag"
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
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(800, 600)

	g := newGame()

	w := etk.NewWindow()

	// Button.
	{
		example, defaultFocus := newButtonExample()
		w.AddChildWithLabel(example, defaultFocus, "Button")
	}

	// Checkbox.
	{
		example, defaultFocus := newCheckboxExample()
		w.AddChildWithLabel(example, defaultFocus, "Checkbox")
	}

	// Flex.
	{
		example, defaultFocus := newFlexExample()
		w.AddChildWithLabel(example, defaultFocus, "Flex")
	}

	// Grid.
	{
		example, defaultFocus := newGridExample()
		w.AddChildWithLabel(example, defaultFocus, "Grid")
	}

	// Input.
	{
		example, defaultFocus := newInputExample()
		w.AddChildWithLabel(example, defaultFocus, "Input")
	}

	// List.
	{
		example, defaultFocus := newListExample()
		w.AddChildWithLabel(example, defaultFocus, "List")
	}

	// Select.
	{
		example, defaultFocus := newSelectExample()
		w.AddChildWithLabel(example, defaultFocus, "Select")
	}

	// Text.
	{
		example, defaultFocus := newTextExample()
		w.AddChildWithLabel(example, defaultFocus, "Text")
	}

	// Window.
	{
		example, defaultFocus := newWindowExample()
		w.AddChildWithLabel(example, defaultFocus, "Window")
	}

	w.Show(0)

	etk.SetRoot(w)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
