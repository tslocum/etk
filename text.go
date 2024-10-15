package etk

import (
	"image"
	"image/color"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
)

// Text is a text display widget.
type Text struct {
	*Box
	field         *messeji.TextField
	textFont      *sfnt.Font
	textSize      int
	textResize    bool
	textAutoSize  int
	scrollVisible bool
	children      []Widget
}

// NewText returns a new Text widget.
func NewText(text string) *Text {
	f := newText()
	f.SetText(text)
	f.SetForegroundColor(Style.TextColorLight)
	f.SetHandleKeyboard(true)

	t := &Text{
		Box:           NewBox(),
		field:         f,
		textFont:      Style.TextFont,
		textSize:      Scale(Style.TextSize),
		scrollVisible: true,
	}
	t.resizeFont()
	return t
}

// SetRect sets the position and size of the widget.
func (t *Text) SetRect(r image.Rectangle) {
	t.Lock()
	defer t.Unlock()

	t.rect = r
	t.field.SetRect(r)
	t.resizeFont()
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
func (t *Text) SetVertical(v Alignment) {
	t.Lock()
	defer t.Unlock()

	t.field.SetVertical(messeji.Alignment(v))
}

// Cursor returns the cursor shape shown when a mouse cursor hovers over the
// widget, or -1 to let widgets beneath determine the cursor shape.
func (t *Text) Cursor() ebiten.CursorShapeType {
	return ebiten.CursorShapeDefault
}

// Write writes to the text buffer.
func (t *Text) Write(p []byte) (n int, err error) {
	t.Lock()
	defer t.Unlock()

	n, err = t.field.Write(p)
	if err != nil {
		return n, err
	}
	t.resizeFont()
	return n, err
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
	t.resizeFont()
}

// SetLast sets the text of the last line of the field.
func (t *Text) SetLast(text string) {
	t.Lock()
	defer t.Unlock()

	t.field.SetLast(text)
	t.resizeFont()
}

func (t *Text) resizeFont() {
	if !t.textResize {
		if t.textAutoSize == t.textSize {
			return
		}
		t.textAutoSize = t.textSize
		ff := FontFace(t.textFont, t.textSize)
		t.field.SetFont(ff, fontMutex)
		return
	}

	w, h := t.rect.Dx()-t.field.Padding()*2, t.rect.Dy()-t.field.Padding()*2
	if w == 0 || h == 0 {
		if t.textAutoSize == t.textSize {
			return
		}
		t.textAutoSize = t.textSize
		ff := FontFace(t.textFont, t.textSize)
		t.field.SetFont(ff, fontMutex)
		return
	}

	var autoSize int
	var ff font.Face
	for autoSize = t.textSize; autoSize > 0; autoSize-- {
		ff = FontFace(t.textFont, autoSize)
		bounds := BoundString(ff, t.field.Text())
		if bounds.Dx() <= w && bounds.Dy() <= h {
			break
		}
	}
	if t.textAutoSize == autoSize {
		return
	}

	t.field.SetFont(ff, fontMutex)
	t.textAutoSize = autoSize
}

func (t *Text) scrollBarVisible() bool {
	if t.textResize {
		return false
	}
	return t.scrollVisible
}

// SetScrollBarVisible sets whether the scroll bar is visible on the screen.
func (t *Text) SetScrollBarVisible(scrollVisible bool) {
	t.Lock()
	defer t.Unlock()

	t.scrollVisible = scrollVisible
	t.field.SetScrollBarVisible(t.scrollBarVisible())
}

// SetAutoHideScrollBar sets whether the scroll bar is automatically hidden
// when the entire text buffer is visible.
func (t *Text) SetAutoHideScrollBar(autoHide bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetAutoHideScrollBar(autoHide)
}

// FontSize returns the font size of the field.
func (t *Text) FontSize() int {
	t.Lock()
	defer t.Unlock()

	return t.textAutoSize
}

// SetFont sets the font and text size of the field. Scaling is not applied.
func (t *Text) SetFont(fnt *sfnt.Font, size int) {
	t.Lock()
	defer t.Unlock()

	t.textFont, t.textSize = fnt, size
	t.resizeFont()
}

// SetAutoResize sets whether the font is automatically scaled down when it is
// too large to fit the entire text buffer on one line.
func (t *Text) SetAutoResize(resize bool) {
	t.Lock()
	defer t.Unlock()

	t.textResize = resize
	t.resizeFont()
	t.field.SetScrollBarVisible(t.scrollBarVisible())
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
