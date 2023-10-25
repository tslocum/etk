package etk

import (
	"image"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	*Box

	Label *messeji.TextField

	onSelected func() error
}

func NewButton(label string, onSelected func() error) *Button {
	textColor := Style.ButtonTextColor
	if textColor == nil {
		textColor = Style.TextColorDark
	}

	l := messeji.NewTextField(Style.TextFont)
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

func (b *Button) SetRect(r image.Rectangle) {
	b.Box.rect = r

	b.Label.SetRect(r)
}

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

func (b *Button) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

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
