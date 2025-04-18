//go:build example

package main

import (
	"bytes"
	"embed"
	"flag"
	"image/color"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"codeberg.org/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.go
var embedFS embed.FS

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

	mono := monoFont()

	g := newGame()

	w := etk.NewWindow()

	addExample := func(f func() (name string, root etk.Widget, focus etk.Widget)) {
		name, root, focus := f()

		buf, err := embedFS.ReadFile(name + ".go")
		if err != nil {
			panic(err)
		}

		t := etk.NewText("")
		t.SetWordWrap(false)
		enableMono := true
		if enableMono {
			t.SetFont(mono, etk.Scale(14))
			t.SetLineHeight(etk.Scale(20))
		}

		buf = bytes.TrimPrefix(buf, []byte("//go:build example\n\n"))
		if buf[len(buf)-1] == '\n' {
			buf = buf[:len(buf)-1]
		}
		_, err = t.Write(bytes.ReplaceAll(buf, []byte("\t"), []byte("  ")))
		if err != nil {
			panic(err)
		}

		borderSize := etk.Scale(7)
		borderShade := uint8(60)
		borderColor := color.RGBA{borderShade, borderShade, borderShade, 255}

		boxA := etk.NewBox()
		boxA.SetBackground(borderColor)
		boxB := etk.NewBox()
		boxB.SetBackground(borderColor)

		grid := etk.NewGrid()
		grid.SetRowSizes(-1, -1, borderSize, -1)
		grid.SetColumnSizes(-1, borderSize)
		grid.AddChildAt(root, 0, 0, 1, 2)
		grid.AddChildAt(boxA, 0, 2, 1, 1)
		grid.AddChildAt(t, 0, 3, 1, 1)
		grid.AddChildAt(boxB, 1, 0, 1, 4)

		w.AddChildWithLabel(grid, focus, strings.Title(name))
	}
	addExample(newButtonExample)
	addExample(newCheckboxExample)
	addExample(newFlexExample)
	addExample(newGridExample)
	addExample(newInputExample)
	addExample(newListExample)
	addExample(newSelectExample)
	addExample(newSpriteExample)
	addExample(newTextExample)
	addExample(newWindowExample)

	w.Show(0)

	etk.SetRoot(w)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
