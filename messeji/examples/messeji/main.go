//go:build example

package main

import (
	"log"

	"code.rocket9labs.com/tslocum/etk/messeji/examples/messeji/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("メッセージ")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetMaxTPS(144)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMinimum)

	g := game.NewDemoGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
