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

	ButtonBorderSize   int
	ButtonBorderTop    color.RGBA
	ButtonBorderRight  color.RGBA
	ButtonBorderBottom color.RGBA
	ButtonBorderLeft   color.RGBA

	InputBorderSize      int
	InputBorderFocused   color.RGBA
	InputBorderUnfocused color.RGBA

	ScrollAreaColor   color.RGBA
	ScrollHandleColor color.RGBA

	ScrollBorderSize   int
	ScrollBorderTop    color.RGBA
	ScrollBorderRight  color.RGBA
	ScrollBorderBottom color.RGBA
	ScrollBorderLeft   color.RGBA

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

	ButtonBorderSize:   4,
	ButtonBorderTop:    color.RGBA{220, 220, 220, 255},
	ButtonBorderRight:  color.RGBA{0, 0, 0, 255},
	ButtonBorderBottom: color.RGBA{0, 0, 0, 255},
	ButtonBorderLeft:   color.RGBA{220, 220, 220, 255},

	InputBorderSize:      2,
	InputBorderFocused:   color.RGBA{220, 220, 220, 255},
	InputBorderUnfocused: color.RGBA{0, 0, 0, 255},

	ScrollAreaColor:   color.RGBA{200, 200, 200, 255},
	ScrollHandleColor: color.RGBA{108, 108, 108, 255},

	ScrollBorderSize:   2,
	ScrollBorderTop:    color.RGBA{240, 240, 240, 255},
	ScrollBorderRight:  color.RGBA{0, 0, 0, 255},
	ScrollBorderBottom: color.RGBA{0, 0, 0, 255},
	ScrollBorderLeft:   color.RGBA{240, 240, 240, 255},

	InputBgColor: color.RGBA{0, 64, 0, 255},

	ButtonBgColor:         color.RGBA{255, 255, 255, 255},
	ButtonBgColorDisabled: color.RGBA{110, 110, 110, 255},
}
