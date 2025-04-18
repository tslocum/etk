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

	// Cursor returns the cursor shape shown when a mouse cursor hovers over
	// the widget, or -1 to let widgets beneath determine the cursor shape.
	Cursor() ebiten.CursorShapeType

	// HandleKeyboard is called when a keyboard event occurs. Either a key or a
	// rune is set, specifying the pressed key. When a key is set, its value is
	// greater than or equal to 0. When a rune is set, the value of key is -1.
	HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error)

	// HandleMouse is called when a mouse event occurs. Only mouse events that
	// are on top of the widget are passed to the widget, except after clicking
	// within the widget and then dragging the cursor outside of the widget.
	HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error)

	// Clip returns whether the widget and its children are restricted to drawing
	// within the widget's rect area of the screen. For best performance, Clip
	// should return false unless clipping is actually needed.
	Clip() bool

	// Draw draws the widget on the screen.
	Draw(screen *ebiten.Image) error

	// Children returns the children of the widget. Children are drawn in the
	// order they are returned. Keyboard and mouse events are passed to children
	// in reverse order.
	Children() []Widget
}

// WithoutFocus wraps a widget to ignore focus.
type WithoutFocus struct {
	Widget
}

// Focus returns the focus state of the widget.
func (w *WithoutFocus) Focus() bool {
	return false
}

// SetFocus sets the focus state of the widget.
func (w *WithoutFocus) SetFocus(focus bool) (accept bool) {
	return false
}

// WithoutMouse wraps a widget to ignore all mouse events.
type WithoutMouse struct {
	Widget
}

// HandleMouse is called when a mouse event occurs. Only mouse events that are
// on top of the widget are passed to the widget.
func (w *WithoutMouse) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

// WithoutMouseExceptScroll wraps a widget to ignore all mouse events except
// scroll events.
type WithoutMouseExceptScroll struct {
	Widget
	handleOnce bool
}

// HandleMouse is called when a mouse event occurs. Only mouse events that are
// on top of the widget are passed to the widget.
func (w *WithoutMouseExceptScroll) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	if pressed || clicked {
		w.handleOnce = true
		return w.Widget.HandleMouse(cursor, pressed, clicked)
	} else if w.handleOnce {
		w.handleOnce = false
		return w.Widget.HandleMouse(cursor, pressed, clicked)
	}
	return false, nil
}
