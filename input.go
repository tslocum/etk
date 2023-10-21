package etk

import (
	"image"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	*Box
	Field  *messeji.InputField
	Cursor string
}

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

// Clear clears the field's buffer.
func (i *Input) Clear() {
	i.Field.SetText("")
}

// Write writes to the field's buffer.
func (i *Input) Write(p []byte) (n int, err error) {
	return i.Field.Write(p)
}

func (i *Input) Text() string {
	return i.Field.Text()
}

func (i *Input) SetRect(r image.Rectangle) {
	i.Box.rect = r

	i.Field.SetRect(r)
}

func (i *Input) SetFocus(focus bool) {
	i.focus = focus

	var cursor string
	if focus {
		cursor = i.Cursor
	}
	i.Field.SetSuffix(cursor)
}

func (i *Input) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

func (i *Input) HandleKeyboard() (handled bool, err error) {
	if !i.focus {
		return false, nil
	}

	err = i.Field.Update()
	return true, err
}

func (i *Input) Draw(screen *ebiten.Image) error {
	// Draw label.
	i.Field.Draw(screen)
	return nil
}
