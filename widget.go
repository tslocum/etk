package etk

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Widget represents an interface element. Most widgets will embed Box and build
// on top of it.
type Widget interface {
	// Rect returns the position and size of the widget.
	Rect() image.Rectangle

	// SetRect sets the position and size of the widget.
	SetRect(r image.Rectangle)

	// Background returns the background color of the widget.
	Background() color.RGBA

	// SetBackground sets the background color of the widget.
	SetBackground(background color.RGBA)

	// Focus returns the focus state of the widget.
	Focus() bool

	// SetFocus sets the focus state of the widget.
	SetFocus(focus bool) (accept bool)

	// Visible returns the visibility of the widget.
	Visible() bool

	// SetVisible sets the visibility of the widget.
	SetVisible(visible bool)

	// HandleKeyboard is called when a keyboard event occurs.
	HandleKeyboard(ebiten.Key, rune) (handled bool, err error)

	// HandleMouse is called when a mouse event occurs. Only mouse events that
	// are on top of the widget are passed to the widget.
	HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error)

	// Draw draws the widget on the screen.
	Draw(screen *ebiten.Image) error

	// Children returns the children of the widget. Children are drawn in the
	// order they are returned. Keyboard and mouse events are passed to children
	// in reverse order.
	Children() []Widget
}

// ignoreMouse wraps a widget to ignore mouse events.
type ignoreMouse struct {
	Widget
}

// HandleMouse is called when a mouse event occurs. Only mouse events that
// are on top of the widget are passed to the widget.
func (i *ignoreMouse) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}
