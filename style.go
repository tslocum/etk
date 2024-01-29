package etk

import (
	"image/color"
	"log"
	"sync"

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

// Attributes represents a default attribute configuration. Integer values will be scaled.
type Attributes struct {
	TextFont      font.Face
	TextFontMutex *sync.Mutex

	TextColorLight color.RGBA
	TextColorDark  color.RGBA

	TextBgColor color.RGBA

	BorderSize        int
	BorderColorTop    color.RGBA
	BorderColorRight  color.RGBA
	BorderColorBottom color.RGBA
	BorderColorLeft   color.RGBA

	ScrollAreaColor   color.RGBA
	ScrollHandleColor color.RGBA

	ScrollBorderSize        int
	ScrollBorderColorTop    color.RGBA
	ScrollBorderColorRight  color.RGBA
	ScrollBorderColorBottom color.RGBA
	ScrollBorderColorLeft   color.RGBA

	InputBgColor color.RGBA

	ButtonTextColor       color.RGBA
	ButtonBgColor         color.RGBA
	ButtonBgColorDisabled color.RGBA
}

// Style is the current default attribute configuration. Integer values will be scaled.
var Style = &Attributes{
	TextFont:      defaultFont(),
	TextFontMutex: &sync.Mutex{},

	TextColorLight: color.RGBA{255, 255, 255, 255},
	TextColorDark:  color.RGBA{0, 0, 0, 255},

	TextBgColor: transparent,

	BorderSize:        4,
	BorderColorTop:    color.RGBA{220, 220, 220, 255},
	BorderColorRight:  color.RGBA{0, 0, 0, 255},
	BorderColorBottom: color.RGBA{0, 0, 0, 255},
	BorderColorLeft:   color.RGBA{220, 220, 220, 255},

	ScrollAreaColor:   color.RGBA{200, 200, 200, 255},
	ScrollHandleColor: color.RGBA{108, 108, 108, 255},

	ScrollBorderSize:        2,
	ScrollBorderColorTop:    color.RGBA{240, 240, 240, 255},
	ScrollBorderColorRight:  color.RGBA{0, 0, 0, 255},
	ScrollBorderColorBottom: color.RGBA{0, 0, 0, 255},
	ScrollBorderColorLeft:   color.RGBA{240, 240, 240, 255},

	InputBgColor: color.RGBA{0, 128, 0, 255},

	ButtonBgColor:         color.RGBA{255, 255, 255, 255},
	ButtonBgColorDisabled: color.RGBA{110, 110, 110, 255},
}
