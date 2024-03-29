//go:build example

package main

import (
	"flag"
	"fmt"
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

	// Input demo.
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

		w.AddChildWithTitle(inputDemo, "Input")
	}

	// Flex demo.
	{
		newButton := func(i int) *etk.Button {
			return etk.NewButton(fmt.Sprintf("Button %d", i), func() error {
				log.Printf("Pressed button %d", i)
				return nil
			})
		}

		b1 := newButton(1)
		b2 := newButton(2)

		topFlex := etk.NewFlex()
		topFlex.AddChild(b1, b2)

		b3 := newButton(3)
		b4 := newButton(4)
		b5 := newButton(5)

		bottomFlex := etk.NewFlex()
		bottomFlex.AddChild(b3, b4, b5)

		flexDemo := etk.NewFlex()
		flexDemo.SetVertical(true)
		flexDemo.AddChild(topFlex, bottomFlex)

		w.AddChildWithTitle(flexDemo, "Flex")
	}

	etk.SetRoot(w)
	etk.SetFocus(input)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
