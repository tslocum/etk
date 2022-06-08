package etk

import "image/color"

type Attributes struct {
	TextColor color.Color

	BorderColor color.Color

	ButtonTextColor       color.Color
	ButtonBgColor         color.Color
	ButtonBgColorDisabled color.Color
}

var Style = &Attributes{
	TextColor: color.RGBA{0, 0, 0, 255},

	BorderColor: color.RGBA{0, 0, 0, 255},

	ButtonBgColor:         color.RGBA{255, 255, 255, 255},
	ButtonBgColorDisabled: color.RGBA{110, 110, 110, 255},
}
