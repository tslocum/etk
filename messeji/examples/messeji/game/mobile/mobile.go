//go:build example

package mobile

import (
	"codeberg.org/tslocum/etk/messeji/examples/messeji/game"
	"github.com/hajimehoshi/ebiten/v2/mobile"
)

func init() {
	mobile.SetGame(game.NewDemoGame())
}

// Dummy is a dummy exported function.
//
// gomobile will only compile packages that include at least one exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
