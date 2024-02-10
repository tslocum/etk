package etk

import (
	"image"
	"image/color"
	"sync"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Input is a text input widget. The Input widget is simply a Text widget that
// also accepts user input.
type Input struct {
	*Box
	field  *messeji.InputField
	cursor string
	focus  bool
}

// NewInput returns a new Input widget.
func NewInput(text string, onSelected func(text string) (handled bool)) *Input {
	textColor := Style.TextColorDark
	/*if TextColor == nil {
		textColor = Style.InputColor
	}*/

	f := messeji.NewInputField(Style.TextFont, Style.TextFontMutex)
	f.SetForegroundColor(textColor)
	f.SetBackgroundColor(transparent)
	f.SetScrollBarColors(Style.ScrollAreaColor, Style.ScrollHandleColor)
	f.SetScrollBorderSize(Scale(Style.ScrollBorderSize))
	f.SetScrollBorderColors(Style.ScrollBorderColorTop, Style.ScrollBorderColorRight, Style.ScrollBorderColorBottom, Style.ScrollBorderColorLeft)
	f.SetPrefix("")
	f.SetSuffix("")
	f.SetText(text)
	f.SetHandleKeyboard(true)
	f.SetSelectedFunc(func() (accept bool) {
		return onSelected(f.Text())
	})

	i := &Input{
		Box:    NewBox(),
		field:  f,
		cursor: "_",
	}
	i.SetBackground(Style.InputBgColor)
	return i
}

// SetRect sets the position and size of the widget.
func (i *Input) SetRect(r image.Rectangle) {
	i.Box.rect = r

	i.field.SetRect(r)

	for _, w := range i.children {
		w.SetRect(r)
	}
}

// Foreground return the color of the text within the field.
func (i *Input) Foreground() color.RGBA {
	i.Lock()
	defer i.Unlock()

	return i.field.ForegroundColor()
}

// SetForegroundColor sets the color of the text within the field.
func (i *Input) SetForeground(c color.RGBA) {
	i.Lock()
	defer i.Unlock()

	i.field.SetForegroundColor(c)
}

// SetPrefix sets the text shown before the input text.
func (i *Input) SetPrefix(prefix string) {
	i.Lock()
	defer i.Unlock()

	i.field.SetPrefix(prefix)
}

// SetSuffix sets the text shown after the input text.
func (i *Input) SetSuffix(suffix string) {
	i.Lock()
	defer i.Unlock()

	i.field.SetSuffix(suffix)
}

// SetCursor sets the cursor appended to the text buffer when focused.
func (i *Input) SetCursor(cursor string) {
	i.Lock()
	defer i.Unlock()

	i.cursor = cursor
	if i.focus {
		i.field.SetSuffix(cursor)
	}
}

// Focus returns the focus state of the widget.
func (i *Input) Focus() bool {
	return i.focus
}

// SetFocus sets the focus state of the widget.
func (i *Input) SetFocus(focus bool) bool {
	i.focus = focus

	var cursor string
	if focus {
		cursor = i.cursor
	}
	i.field.SetSuffix(cursor)
	return true
}

// Text returns the content of the text buffer.
func (i *Input) Text() string {
	i.Lock()
	defer i.Unlock()

	return i.field.Text()
}

// SetText sets the text in the field.
func (i *Input) SetText(text string) {
	i.Lock()
	defer i.Unlock()

	i.field.SetText(text)
}

// SetScrollBarWidth sets the width of the scroll bar.
func (i *Input) SetScrollBarWidth(width int) {
	i.Lock()
	defer i.Unlock()

	i.field.SetScrollBarWidth(width)
}

// SetScrollBarColors sets the color of the scroll bar area and handle.
func (i *Input) SetScrollBarColors(area color.RGBA, handle color.RGBA) {
	i.Lock()
	defer i.Unlock()

	i.field.SetScrollBarColors(Style.ScrollAreaColor, Style.ScrollHandleColor)
}

// SetScrollBarVisible sets whether the scroll bar is visible on the screen.
func (i *Input) SetScrollBarVisible(scrollVisible bool) {
	i.Lock()
	defer i.Unlock()

	i.field.SetScrollBarVisible(scrollVisible)
}

// SetAutoHideScrollBar sets whether the scroll bar is automatically hidden
// when the entire text buffer is visible.
func (i *Input) SetAutoHideScrollBar(autoHide bool) {
	i.Lock()
	defer i.Unlock()

	i.field.SetAutoHideScrollBar(autoHide)
}

// SetFont sets the font face of the text within the field.
func (i *Input) SetFont(face font.Face, mutex *sync.Mutex) {
	i.Lock()
	defer i.Unlock()

	i.field.SetFont(face, mutex)
}

// Padding returns the amount of padding around the text within the field.
func (i *Input) Padding() int {
	i.Lock()
	defer i.Unlock()

	return i.field.Padding()
}

// SetPadding sets the amount of padding around the text within the field.
func (i *Input) SetPadding(padding int) {
	i.Lock()
	defer i.Unlock()

	i.field.SetPadding(padding)
}

// SetWordWrap sets a flag which, when enabled, causes text to wrap without breaking words.
func (i *Input) SetWordWrap(wrap bool) {
	i.Lock()
	defer i.Unlock()

	i.field.SetWordWrap(wrap)
}

// SetHorizontal sets the horizontal alignment of the text within the field.
func (i *Input) SetHorizontal(h Alignment) {
	i.Lock()
	defer i.Unlock()

	i.field.SetHorizontal(messeji.Alignment(h))
}

// SetVertical sets the vertical alignment of the text within the field.
func (i *Input) SetVertical(h Alignment) {
	i.Lock()
	defer i.Unlock()

	i.field.SetVertical(messeji.Alignment(h))
}

// SetMask sets the rune used to mask the text buffer contents. Set to 0 to disable.
func (i *Input) SetMask(r rune) {
	i.Lock()
	defer i.Unlock()

	i.field.SetMask(r)
}

// Write writes to the text buffer.
func (i *Input) Write(p []byte) (n int, err error) {
	return i.field.Write(p)
}

// HandleKeyboard is called when a keyboard event occurs.
func (i *Input) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	if !i.focus {
		return false, nil
	}

	return i.field.HandleKeyboardEvent(key, r)
}

// HandleMouse is called when a mouse event occurs.
func (i *Input) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return i.field.HandleMouseEvent(cursor, pressed, clicked)
}

// Draw draws the widget on the screen.
func (i *Input) Draw(screen *ebiten.Image) error {
	i.field.Draw(screen)
	return nil
}
