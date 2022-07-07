package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Rect() image.Rectangle
	SetRect(r image.Rectangle)
	HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error)
	HandleKeyboard() (handled bool, err error)
	Draw(screen *ebiten.Image) error
	Children() []Widget
}
