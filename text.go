package etk

import (
	"image"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

type Text struct {
	*messeji.TextField

	children []Widget
}

func NewText(text string) *Text {
	textColor := Style.TextColorLight

	l := messeji.NewTextField(Style.TextFont)
	l.SetText(text)
	l.SetForegroundColor(textColor)
	l.SetBackgroundColor(Style.TextBgColor)

	return &Text{
		TextField: l,
	}
}

func (t *Text) SetFocus(focus bool) {
	// Do nothing.
}

func (t *Text) Children() []Widget {
	t.Lock()
	defer t.Unlock()

	return t.children
}

func (t *Text) AddChild(w ...Widget) {
	t.Lock()
	defer t.Unlock()

	t.children = append(t.children, w...)
}

// Clear clears the field's buffer.
func (t *Text) Clear() {
	t.SetText("")
}

// Write writes to the field's buffer.
func (t *Text) Write(p []byte) (n int, err error) {
	return t.Write(p)
}

func (t *Text) Text() string {
	return t.Text()
}

func (t *Text) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

func (t *Text) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

func (t *Text) Draw(screen *ebiten.Image) error {
	// Draw label.
	t.TextField.Draw(screen)
	return nil
}
