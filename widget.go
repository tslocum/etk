package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Rect() image.Rectangle
	SetRect(r image.Rectangle)
	HandleMouse() (handled bool, err error)
	HandleKeyboard() (handled bool, err error)
	Draw(screen *ebiten.Image) error
	Children() []Widget
}
