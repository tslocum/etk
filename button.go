package etk

import (
	"image"
	"image/color"
	"log"

	"code.rocketnine.space/tslocum/messeji"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO
var mplusNormalFont font.Face

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Button struct {
	*Box

	label *messeji.TextField
}

func NewButton(label string, onSelected func()) *Button {
	textColor := Style.ButtonTextColor
	if textColor == nil {
		textColor = Style.TextColor
	}

	l := messeji.NewTextField(mplusNormalFont)
	l.SetText(label)
	l.SetForegroundColor(textColor)
	l.SetBackgroundColor(color.RGBA{0, 0, 0, 0})
	l.SetHorizontal(messeji.AlignCenter)
	l.SetVertical(messeji.AlignCenter)

	return &Button{
		Box:   NewBox(),
		label: l, // TODO
	}
}

func (b *Button) SetRect(r image.Rectangle) {
	b.Box.rect = r

	b.label.SetRect(r)
}

func (b *Button) HandleMouse() (handled bool, err error) {
	return false, nil
}

func (b *Button) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

func (b *Button) Draw(screen *ebiten.Image) error {
	// TODO background color
	// Draw background.
	screen.SubImage(b.rect).(*ebiten.Image).Fill(Style.ButtonBgColor)

	// Draw label.
	b.label.Draw(screen)

	// Draw border.
	const borderSize = 4
	screen.SubImage(image.Rect(b.rect.Min.X, b.rect.Min.Y, b.rect.Max.X, b.rect.Min.Y+borderSize)).(*ebiten.Image).Fill(Style.BorderColor)
	screen.SubImage(image.Rect(b.rect.Min.X, b.rect.Max.Y-borderSize, b.rect.Max.X, b.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
	screen.SubImage(image.Rect(b.rect.Min.X, b.rect.Min.Y, b.rect.Min.X+borderSize, b.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
	screen.SubImage(image.Rect(b.rect.Max.X-borderSize, b.rect.Min.Y, b.rect.Max.X, b.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)

	return nil
}
