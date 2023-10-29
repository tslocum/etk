package etk

import (
	"image"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

// Input is a text input widget. The Input widget is simply a Text widget that
// also accepts user input.
type Input struct {
	*Box
	Field  *messeji.InputField
	Cursor string
	focus  bool
}

// NewInput returns a new Input widget.
func NewInput(prefix string, text string, onSelected func(text string) (handled bool)) *Input {
	textColor := Style.TextColorDark
	/*if TextColor == nil {
		textColor = Style.InputColor
	}*/

	i := messeji.NewInputField(Style.TextFont)
	i.SetPrefix(prefix)
	i.SetSuffix("")
	i.SetText(text)
	i.SetForegroundColor(textColor)
	i.SetBackgroundColor(Style.InputBgColor)
	i.SetHandleKeyboard(true)
	i.SetSelectedFunc(func() (accept bool) {
		return onSelected(i.Text())
	})

	return &Input{
		Box:    NewBox(),
		Field:  i,
		Cursor: "_",
	}
}

// SetRect sets the position and size of the widget.
func (i *Input) SetRect(r image.Rectangle) {
	i.Box.rect = r

	i.Field.SetRect(r)
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
		cursor = i.Cursor
	}
	i.Field.SetSuffix(cursor)
	return true
}

// Clear clears the textbuffer.
func (i *Input) Clear() {
	i.Field.SetText("")
}

// Write writes to the text buffer.
func (i *Input) Write(p []byte) (n int, err error) {
	return i.Field.Write(p)
}

// Text returns the content of the text buffer.
func (i *Input) Text() string {
	return i.Field.Text()
}

// HandleKeyboard is called when a keyboard event occurs.
func (i *Input) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	if !i.focus {
		return false, nil
	}

	return i.Field.HandleKeyboardEvent(key, r)
}

// HandleMouse is called when a mouse event occurs.
func (i *Input) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return i.Field.HandleMouseEvent(cursor, pressed, clicked)
}

// Draw draws the widget on the screen.
func (i *Input) Draw(screen *ebiten.Image) error {
	i.Field.Draw(screen)
	return nil
}
