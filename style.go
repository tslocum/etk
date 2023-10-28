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

	TextColorLight color.RGBA
	TextColorDark  color.RGBA

	TextBgColor color.RGBA

	ScrollAreaColor   color.RGBA
	ScrollHandleColor color.RGBA

	BorderColor color.RGBA

	InputBgColor color.RGBA

	ButtonTextColor       color.RGBA
	ButtonBgColor         color.RGBA
	ButtonBgColorDisabled color.RGBA
}

var Style = &Attributes{
	TextFont: defaultFont(),

	TextColorLight: color.RGBA{255, 255, 255, 255},
	TextColorDark:  color.RGBA{0, 0, 0, 255},

	TextBgColor: transparent,

	ScrollAreaColor:   color.RGBA{200, 200, 200, 255},
	ScrollHandleColor: color.RGBA{108, 108, 108, 255},

	BorderColor: color.RGBA{0, 0, 0, 255},

	InputBgColor: color.RGBA{0, 128, 0, 255},

	ButtonBgColor:         color.RGBA{255, 255, 255, 255},
	ButtonBgColorDisabled: color.RGBA{110, 110, 110, 255},
}
