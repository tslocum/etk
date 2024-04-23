package etk

import (
	"image"
	"image/color"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/sfnt"
)

// Text is a text display widget.
type Text struct {
	*Box
	field    *messeji.TextField
	children []Widget
}

// NewText returns a new Text widget.
func NewText(text string) *Text {
	f := newText()
	f.SetText(text)
	f.SetForegroundColor(Style.TextColorLight)
	f.SetHandleKeyboard(true)

	return &Text{
		Box:   NewBox(),
		field: f,
	}
}

// SetRect sets the position and size of the widget.
func (t *Text) SetRect(r image.Rectangle) {
	t.Lock()
	defer t.Unlock()

	t.rect = r
	t.field.SetRect(r)
}

// Foreground return the color of the text within the field.
func (t *Text) Foreground() color.RGBA {
	t.Lock()
	defer t.Unlock()

	return t.field.ForegroundColor()
}

// SetForegroundColor sets the color of the text within the field.
func (t *Text) SetForeground(c color.RGBA) {
	t.Lock()
	defer t.Unlock()

	t.field.SetForegroundColor(c)
}

// Focus returns the focus state of the widget.
func (t *Text) Focus() bool {
	return false
}

// SetFocus sets the focus state of the widget.
func (t *Text) SetFocus(focus bool) bool {
	return false
}

// SetScrollBarWidth sets the width of the scroll bar.
func (t *Text) SetScrollBarWidth(width int) {
	t.Lock()
	defer t.Unlock()

	t.field.SetScrollBarWidth(width)
}

// SetScrollBarColors sets the color of the scroll bar area and handle.
func (t *Text) SetScrollBarColors(area color.RGBA, handle color.RGBA) {
	t.Lock()
	defer t.Unlock()

	t.field.SetScrollBarColors(Style.ScrollAreaColor, Style.ScrollHandleColor)
}

// SetScrollBorderColor sets the color of the top, right, bottom and left border
// of the scroll bar handle.
func (t *Text) SetScrollBorderColors(top color.RGBA, right color.RGBA, bottom color.RGBA, left color.RGBA) {
	t.Lock()
	defer t.Unlock()

	t.field.SetScrollBorderColors(top, right, bottom, left)
}

// SetWordWrap sets a flag which, when enabled, causes text to wrap without breaking words.
func (t *Text) SetWordWrap(wrap bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetWordWrap(wrap)
}

// SetHorizontal sets the horizontal alignment of the text within the field.
func (t *Text) SetHorizontal(h Alignment) {
	t.Lock()
	defer t.Unlock()

	t.field.SetHorizontal(messeji.Alignment(h))
}

// SetVertical sets the vertical alignment of the text within the field.
func (t *Text) SetVertical(h Alignment) {
	t.Lock()
	defer t.Unlock()

	t.field.SetVertical(messeji.Alignment(h))
}

// Write writes to the text buffer.
func (t *Text) Write(p []byte) (n int, err error) {
	t.Lock()
	defer t.Unlock()

	return t.field.Write(p)
}

// Text returns the content of the text buffer.
func (t *Text) Text() string {
	t.Lock()
	defer t.Unlock()

	return t.field.Text()
}

// SetText sets the text in the field.
func (t *Text) SetText(text string) {
	t.Lock()
	defer t.Unlock()

	t.field.SetText(text)
}

// SetScrollBarVisible sets whether the scroll bar is visible on the screen.
func (t *Text) SetScrollBarVisible(scrollVisible bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetScrollBarVisible(scrollVisible)
}

// SetAutoHideScrollBar sets whether the scroll bar is automatically hidden
// when the entire text buffer is visible.
func (t *Text) SetAutoHideScrollBar(autoHide bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetAutoHideScrollBar(autoHide)
}

// SetFont sets the font and text size of the field. Scaling is not applied.
func (t *Text) SetFont(fnt *sfnt.Font, size int) {
	t.Lock()
	defer t.Unlock()

	t.field.SetFont(FontFace(fnt, size), fontMutex)
}

// Padding returns the amount of padding around the text within the field.
func (t *Text) Padding() int {
	t.Lock()
	defer t.Unlock()

	return t.field.Padding()
}

// SetPadding sets the amount of padding around the text within the field.
func (t *Text) SetPadding(padding int) {
	t.Lock()
	defer t.Unlock()

	t.field.SetPadding(padding)
}

// SetFollow sets whether the field should automatically scroll to the end when
// content is added to the buffer.
func (t *Text) SetFollow(follow bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetFollow(follow)
}

// SetSingleLine sets whether the field displays all text on a single line.
// When enabled, the field scrolls horizontally. Otherwise, it scrolls vertically.
func (t *Text) SetSingleLine(single bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetSingleLine(single)
}

// SetMask sets the rune used to mask the text buffer contents. Set to 0 to disable.
func (t *Text) SetMask(r rune) {
	t.Lock()
	defer t.Unlock()

	t.field.SetMask(r)
}

// HandleKeyboard is called when a keyboard event occurs.
func (t *Text) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	return t.field.HandleKeyboardEvent(key, r)
}

// HandleMouse is called when a mouse event occurs.
func (t *Text) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return t.field.HandleMouseEvent(cursor, pressed, clicked)
}

// Draw draws the widget on the screen.
func (t *Text) Draw(screen *ebiten.Image) error {
	t.field.Draw(screen)
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
