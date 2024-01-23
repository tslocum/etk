package etk

import (
	"image"
	"image/color"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
)

// Button is a clickable button.
type Button struct {
	*Box

	Label *messeji.TextField

	borderTop    color.RGBA
	borderRight  color.RGBA
	borderBottom color.RGBA
	borderLeft   color.RGBA
	onSelected   func() error
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
		Box:          NewBox(),
		Label:        l,
		onSelected:   onSelected,
		borderTop:    Style.BorderColorTop,
		borderRight:  Style.BorderColorRight,
		borderBottom: Style.BorderColorBottom,
		borderLeft:   Style.BorderColorLeft,
	}
}

// SetRect sets the position and size of the Button.
func (b *Button) SetRect(r image.Rectangle) {
	b.Box.rect = r

	b.Label.SetRect(r)

	for _, w := range b.children {
		w.SetRect(r)
	}
}

// SetBorderColor sets the color of the top, right, bottom and left border.
func (b *Button) SetBorderColor(top color.RGBA, right color.RGBA, bottom color.RGBA, left color.RGBA) {
	b.Lock()
	defer b.Unlock()

	b.borderTop = top
	b.borderRight = right
	b.borderBottom = bottom
	b.borderLeft = left
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
	r := b.rect

	// Draw background.
	screen.SubImage(r).(*ebiten.Image).Fill(Style.ButtonBgColor)

	// Draw label.
	b.Label.Draw(screen)

	// Draw border.
	const borderSize = 4
	screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+borderSize, r.Max.Y)).(*ebiten.Image).Fill(b.borderLeft)
	screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+borderSize)).(*ebiten.Image).Fill(b.borderTop)
	screen.SubImage(image.Rect(r.Max.X-borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(b.borderRight)
	screen.SubImage(image.Rect(r.Min.X, r.Max.Y-borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(b.borderBottom)

	return nil
}
