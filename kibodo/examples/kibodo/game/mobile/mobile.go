//go:build example

package mobile

import (
	"codeberg.org/tslocum/etk/kibodo/examples/kibodo/game"
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
