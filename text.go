package etk

import (
	"image"
	"image/color"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

// Text is a text display widget.
type Text struct {
	*messeji.TextField

	background color.RGBA
	children   []Widget
}

// NewText returns a new Text widget.
func NewText(text string) *Text {
	textColor := Style.TextColorLight

	l := messeji.NewTextField(Style.TextFont, Style.TextFontMutex)
	l.SetText(text)
	l.SetForegroundColor(textColor)
	l.SetBackgroundColor(Style.TextBgColor)
	l.SetScrollBarColors(Style.ScrollAreaColor, Style.ScrollHandleColor)
	l.SetHandleKeyboard(true)

	return &Text{
		TextField: l,
	}
}

// Background returns the background color of the widget.
func (t *Text) Background() color.RGBA {
	t.Lock()
	defer t.Unlock()

	return t.background
}

// SetBackground sets the background color of the widget.
func (t *Text) SetBackground(background color.RGBA) {
	t.Lock()
	defer t.Unlock()

	t.background = background
}

// Focus returns the focus state of the widget.
func (t *Text) Focus() bool {
	return false
}

// SetFocus sets the focus state of the widget.
func (t *Text) SetFocus(focus bool) bool {
	return false
}

// Clear clears the text buffer.
func (t *Text) Clear() {
	t.TextField.SetText("")
}

// Write writes to the text buffer.
func (t *Text) Write(p []byte) (n int, err error) {
	return t.TextField.Write(p)
}

// Text returns the content of the text buffer.
func (t *Text) Text() string {
	return t.TextField.Text()
}

// HandleKeyboard is called when a keyboard event occurs.
func (t *Text) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	return t.TextField.HandleKeyboardEvent(key, r)
}

// HandleMouse is called when a mouse event occurs.
func (t *Text) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return t.TextField.HandleMouseEvent(cursor, pressed, clicked)
}

// Draw draws the widget on the screen.
func (t *Text) Draw(screen *ebiten.Image) error {
	t.TextField.Draw(screen)
	return nil
}

// Children returns the children of the widget.
func (t *Text) Children() []Widget {
	t.Lock()
	defer t.Unlock()

	return t.children
}

// AddChild adds a child to the widget.
func (t *Text) AddChild(w ...Widget) {
	t.Lock()
	defer t.Unlock()

	t.children = append(t.children, w...)
}

var _ Widget = &Text{}
