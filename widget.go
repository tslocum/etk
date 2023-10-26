package etk

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Rect() image.Rectangle
	SetRect(r image.Rectangle)
	Background() color.RGBA
	SetBackground(background color.RGBA)
	SetFocus(focus bool) (accept bool)
	SetVisible(visible bool)
	Visible() bool
	HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error)
	HandleKeyboard() (handled bool, err error)
	HandleKeyboardEvent(ebiten.Key, rune) (handled bool, err error)
	Draw(screen *ebiten.Image) error
	Children() []Widget
}
