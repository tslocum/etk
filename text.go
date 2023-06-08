package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"code.rocketnine.space/tslocum/messeji"
)

type Text struct {
	*Box
	field *messeji.TextField
}

func NewText(text string) *Text {
	textColor := Style.TextColorLight

	l := messeji.NewTextField(Style.TextFont)
	l.SetText(text)
	l.SetForegroundColor(textColor)
	l.SetBackgroundColor(Style.TextBgColor)
	l.SetHorizontal(messeji.AlignCenter)
	l.SetVertical(messeji.AlignCenter)

	return &Text{
		Box:   NewBox(),
		field: l,
	}
}

// Clear clears the field's buffer.
func (t *Text) Clear() {
	t.field.SetText("")
}

// Write writes to the field's buffer.
func (t *Text) Write(p []byte) (n int, err error) {
	return t.field.Write(p)
}

func (t *Text) Text() string {
	return t.field.Text()
}

func (t *Text) SetRect(r image.Rectangle) {
	t.Box.rect = r

	t.field.SetRect(r)
}

func (t *Text) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

func (t *Text) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

func (t *Text) Draw(screen *ebiten.Image) error {
	// Draw label.
	t.field.Draw(screen)
	return nil
}
