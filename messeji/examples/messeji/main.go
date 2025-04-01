//go:build example

package main

import (
	"log"

	"codeberg.org/tslocum/etk/messeji/examples/messeji/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("メッセージ")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetTPS(144)
	ebiten.SetVsyncEnabled(true)

	g := game.NewDemoGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
