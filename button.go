package etk

import (
	"image"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

// Button is a clickable button.
type Button struct {
	*Box

	Label *messeji.TextField

	onSelected func() error
}

// NewButton returns a new Button widget.
func NewButton(label string, onSelected func() error) *Button {
	textColor := Style.ButtonTextColor
	if textColor.A == 0 {
		textColor = Style.TextColorDark
	}

	l := messeji.NewTextField(Style.TextFont, Style.TextFontMutex)
	l.SetText(label)
	l.SetForegroundColor(textColor)
	l.SetBackgroundColor(transparent)
	l.SetHorizontal(messeji.AlignCenter)
	l.SetVertical(messeji.AlignCenter)
	l.SetScrollBarVisible(false)

	return &Button{
		Box:        NewBox(),
		Label:      l,
		onSelected: onSelected,
	}
}

// SetRect sets the position and size of the Button.
func (b *Button) SetRect(r image.Rectangle) {
	b.Box.rect = r

	b.Label.SetRect(r)
}

// HandleKeyboard is called when a keyboard event occurs.
func (b *Button) HandleKeyboard(ebiten.Key, rune) (handled bool, err error) {
	return false, nil
}

// HandleMouse is called when a mouse event occurs.
func (b *Button) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	if !clicked {
		return true, nil
	}

	b.Lock()
	onSelected := b.onSelected
	if onSelected == nil {
		b.Unlock()
		return true, nil
	}
	b.Unlock()

	return true, onSelected()
}

// Draw draws the button on the screen.
func (b *Button) Draw(screen *ebiten.Image) error {
	// TODO background color
	// Draw background.
	screen.SubImage(b.rect).(*ebiten.Image).Fill(Style.ButtonBgColor)

	// Draw label.
	b.Label.Draw(screen)

	// Draw border.
	const borderSize = 4
	screen.SubImage(image.Rect(b.rect.Min.X, b.rect.Min.Y, b.rect.Max.X, b.rect.Min.Y+borderSize)).(*ebiten.Image).Fill(Style.BorderColor)
	screen.SubImage(image.Rect(b.rect.Min.X, b.rect.Max.Y-borderSize, b.rect.Max.X, b.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
	screen.SubImage(image.Rect(b.rect.Min.X, b.rect.Min.Y, b.rect.Min.X+borderSize, b.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
	screen.SubImage(image.Rect(b.rect.Max.X-borderSize, b.rect.Min.Y, b.rect.Max.X, b.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)

	return nil
}
