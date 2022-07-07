package etk

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var transparent = color.RGBA{0, 0, 0, 0}

func defaultFont() font.Face {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	defaultFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	return defaultFont
}

type Attributes struct {
	TextFont font.Face

	TextColorLight color.Color
	TextColorDark  color.Color

	TextBgColor color.Color

	BorderColor color.Color

	InputBgColor color.Color

	ButtonTextColor       color.Color
	ButtonBgColor         color.Color
	ButtonBgColorDisabled color.Color
}

var Style = &Attributes{
	TextFont: defaultFont(),

	TextColorLight: color.RGBA{255, 255, 255, 255},
	TextColorDark:  color.RGBA{0, 0, 0, 255},

	TextBgColor: transparent,

	BorderColor: color.RGBA{0, 0, 0, 255},

	InputBgColor: color.RGBA{0, 128, 0, 255},

	ButtonBgColor:         color.RGBA{255, 255, 255, 255},
	ButtonBgColorDisabled: color.RGBA{110, 110, 110, 255},
}
