package etk

import (
	"image"
	"image/color"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

type Text struct {
	*messeji.TextField

	background color.RGBA
	children   []Widget
}

func NewText(text string) *Text {
	textColor := Style.TextColorLight

	l := messeji.NewTextField(Style.TextFont)
	l.SetText(text)
	l.SetForegroundColor(textColor)
	l.SetBackgroundColor(Style.TextBgColor)
	l.SetHandleKeyboard(true)

	return &Text{
		TextField: l,
	}
}

func (t *Text) Background() color.RGBA {
	t.Lock()
	defer t.Unlock()

	return t.background
}

func (t *Text) SetBackground(background color.RGBA) {
	t.Lock()
	defer t.Unlock()

	t.background = background
}

func (t *Text) SetFocus(focus bool) bool {
	return false
}

func (t *Text) Focus() bool {
	return false
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
	t.TextField.SetText("")
}

// Write writes to the field's buffer.
func (t *Text) Write(p []byte) (n int, err error) {
	return t.TextField.Write(p)
}

func (t *Text) Text() string {
	return t.TextField.Text()
}

func (t *Text) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	return t.TextField.HandleKeyboardEvent(key, r)
}

func (t *Text) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return t.TextField.HandleMouseEvent(cursor, pressed, clicked)
}

func (t *Text) Draw(screen *ebiten.Image) error {
	// Draw label.
	t.TextField.Draw(screen)
	return nil
}

var _ Widget = &Text{}
