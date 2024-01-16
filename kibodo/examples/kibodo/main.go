//go:build example

package main

import (
	"log"

	"code.rocket9labs.com/tslocum/etk/kibodo/examples/kibodo/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("キーボード")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetTPS(60)
	ebiten.SetVsyncEnabled(true)

	g := game.NewDemoGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
