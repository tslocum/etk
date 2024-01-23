package etk

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Box is a building block for other widgets. It may also be used as a spacer
// in layout widgets.
type Box struct {
	rect       image.Rectangle
	children   []Widget
	background color.RGBA
	visible    bool

	sync.Mutex
}

// NewBox returns a new Box widget.
func NewBox() *Box {
	return &Box{
		background: transparent,
		visible:    true,
	}
}

// Rect returns the position and size of the widget.
func (b *Box) Rect() image.Rectangle {
	b.Lock()
	defer b.Unlock()

	return b.rect
}

// SetRect sets the position and size of the widget.
func (b *Box) SetRect(r image.Rectangle) {
	b.Lock()
	b.rect = r
	b.Unlock()

	for _, w := range b.children {
		w.SetRect(r)
	}
}

// Background returns the background color of the widget.
func (b *Box) Background() color.RGBA {
	b.Lock()
	defer b.Unlock()

	return b.background
}

// SetBackground sets the background color of the widget.
func (b *Box) SetBackground(background color.RGBA) {
	b.Lock()
	defer b.Unlock()

	b.background = background
}

// Focus returns the focus state of the widget.
func (b *Box) Focus() bool {
	return false
}

// SetFocus sets the focus state of the widget.
func (b *Box) SetFocus(focus bool) bool {
	return false
}

// Visible returns the visibility of the widget.
func (b *Box) Visible() bool {
	return b.visible
}

// SetVisible sets the visibility of the widget.
func (b *Box) SetVisible(visible bool) {
	b.visible = visible
}

// HandleKeyboard is called when a keyboard event occurs.
func (b *Box) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	return false, nil
}

// HandleMouse is called when a mouse event occurs. Only mouse events that
// are on top of the widget are passed to the widget.
func (b *Box) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

// Draw draws the widget on the screen.
func (b *Box) Draw(screen *ebiten.Image) error {
	return nil
}

// Children returns the children of the widget. Children are drawn in the
// order they are returned. Keyboard and mouse events are passed to children
// in reverse order.
func (b *Box) Children() []Widget {
	b.Lock()
	defer b.Unlock()

	return b.children
}

// AddChild adds a child to the widget.
func (b *Box) AddChild(w ...Widget) {
	b.Lock()
	defer b.Unlock()

	b.children = append(b.children, w...)
}

// Clear removes all children from the widget.
func (b *Box) Clear() {
	b.Lock()
	defer b.Unlock()

	b.children = b.children[:0]
}

var _ Widget = &Box{}
