package etk

import (
	"image"
	"image/color"

	"codeberg.org/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Input is a text input widget. The Input widget is simply a Text widget that
// also accepts user input.
type Input struct {
	*Box
	field           *messeji.InputField
	onChange        func(text string, r rune) (accept bool)
	onConfirm       func(text string) (handled bool)
	cursor          string
	borderSize      int
	borderFocused   color.RGBA
	borderUnfocused color.RGBA
	focus           bool
}

// NewInput returns a new Input widget.
func NewInput(text string, onChange func(text string, r rune) (accept bool), onConfirm func(text string) (handled bool)) *Input {
	f := messeji.NewInputField(Style.TextFont, Scale(Style.TextSize), fontMutex)
	f.SetForegroundColor(Style.TextColorLight)
	f.SetBackgroundColor(transparent)
	f.SetScrollBarColors(Style.ScrollAreaColor, Style.ScrollHandleColor)
	f.SetScrollBorderSize(Scale(Style.ScrollBorderSize))
	f.SetScrollBorderColors(Style.ScrollBorderTop, Style.ScrollBorderRight, Style.ScrollBorderBottom, Style.ScrollBorderLeft)
	f.SetPrefix("")
	f.SetSuffix("")
	f.SetText(text)
	f.SetHandleKeyboard(true)

	i := &Input{
		Box:             NewBox(),
		field:           f,
		onChange:        onChange,
		onConfirm:       onConfirm,
		cursor:          "_",
		borderSize:      Scale(Style.InputBorderSize),
		borderFocused:   Style.InputBorderFocused,
		borderUnfocused: Style.InputBorderUnfocused,
	}
	i.SetBackground(Style.InputBgColor)
	f.SetChangedFunc(func(r rune) (accept bool) {
		if i.onChange != nil {
			text := f.Text()
			if r == 0 && len(text) > 0 {
				return i.onChange(text[:len(text)-1], 0)
			}
			return i.onChange(text+string(r), r)
		}
		return true
	})
	f.SetSelectedFunc(func() (accept bool) {
		if i.onConfirm != nil {
			return i.onConfirm(f.Text())
		}
		return true
	})
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

// SetBorderSize sets the size of the border around the field.
func (i *Input) SetBorderSize(size int) {
	i.Lock()
	defer i.Unlock()

	i.borderSize = size
}

// SetBorderColors sets the border colors of the field when focused and unfocused.
func (i *Input) SetBorderColors(focused color.RGBA, unfocused color.RGBA) {
	i.Lock()
	defer i.Unlock()

	i.borderFocused = focused
	i.borderUnfocused = unfocused
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

// SetFont sets the font and text size of the field. Scaling is not applied.
func (t *Input) SetFont(fnt *text.GoTextFaceSource, size int) {
	t.Lock()
	defer t.Unlock()

	t.field.SetFont(fnt, size, fontMutex)
}

// SetAutoResize sets whether the font is automatically scaled down when it is
// too large to fit the entire text buffer on one line.
func (t *Input) SetAutoResize(resize bool) {
	t.Lock()
	defer t.Unlock()

	t.field.SetAutoResize(resize)
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
func (i *Input) SetVertical(v Alignment) {
	i.Lock()
	defer i.Unlock()

	i.field.SetVertical(messeji.Alignment(v))
}

// SetMask sets the rune used to mask the text buffer contents. Set to 0 to disable.
func (i *Input) SetMask(r rune) {
	i.Lock()
	defer i.Unlock()

	i.field.SetMask(r)
}

// SetChangeFunc sets the handler called when the text input changes. When the
// backspace key is pressed, the current text and a rune value of 0 is passed.
func (i *Input) SetChangeFunc(onChange func(text string, r rune) (accept bool)) {
	i.Lock()
	defer i.Unlock()

	i.onChange = onChange
}

// SetConfirmFunc sets the handler called when the text input is confirmed.
func (i *Input) SetConfirmFunc(onConfirm func(text string) (handled bool)) {
	i.Lock()
	defer i.Unlock()

	i.onConfirm = onConfirm
}

// Cursor returns the cursor shape shown when a mouse cursor hovers over the
// widget, or -1 to let widgets beneath determine the cursor shape.
func (i *Input) Cursor() ebiten.CursorShapeType {
	return ebiten.CursorShapeText
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

	// Draw border.
	if i.borderSize == 0 {
		return nil
	}
	r := i.rect
	c := i.borderUnfocused
	if i.focus {
		c = i.borderFocused
	}
	screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+i.borderSize, r.Max.Y)).(*ebiten.Image).Fill(c)
	screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+i.borderSize)).(*ebiten.Image).Fill(c)
	screen.SubImage(image.Rect(r.Max.X-i.borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(c)
	screen.SubImage(image.Rect(r.Min.X, r.Max.Y-i.borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(c)
	return nil
}
