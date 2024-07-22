package etk

import (
	"image/color"

	"golang.org/x/image/font/sfnt"
)

var transparent = color.RGBA{0, 0, 0, 0}

// Attributes represents a default attribute configuration. Integer values will be scaled.
type Attributes struct {
	TextFont *sfnt.Font
	TextSize int

	TextColorLight color.RGBA
	TextColorDark  color.RGBA

	TextBgColor color.RGBA

	BorderSize int

	BorderColorTop    color.RGBA
	BorderColorRight  color.RGBA
	BorderColorBottom color.RGBA
	BorderColorLeft   color.RGBA

	ScrollAreaColor   color.RGBA
	ScrollHandleColor color.RGBA

	ScrollBorderSize int

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
	TextSize: 32,

	TextColorLight: color.RGBA{255, 255, 255, 255},
	TextColorDark:  color.RGBA{0, 0, 0, 255},

	TextBgColor: transparent,

	BorderSize: 4,

	BorderColorTop:    color.RGBA{220, 220, 220, 255},
	BorderColorRight:  color.RGBA{0, 0, 0, 255},
	BorderColorBottom: color.RGBA{0, 0, 0, 255},
	BorderColorLeft:   color.RGBA{220, 220, 220, 255},

	ScrollAreaColor:   color.RGBA{200, 200, 200, 255},
	ScrollHandleColor: color.RGBA{108, 108, 108, 255},

	ScrollBorderSize: 2,

	ScrollBorderColorTop:    color.RGBA{240, 240, 240, 255},
	ScrollBorderColorRight:  color.RGBA{0, 0, 0, 255},
	ScrollBorderColorBottom: color.RGBA{0, 0, 0, 255},
	ScrollBorderColorLeft:   color.RGBA{240, 240, 240, 255},

	InputBgColor: color.RGBA{0, 128, 0, 255},

	ButtonBgColor:         color.RGBA{255, 255, 255, 255},
	ButtonBgColorDisabled: color.RGBA{110, 110, 110, 255},
}
