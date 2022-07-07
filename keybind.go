package etk

import "github.com/hajimehoshi/ebiten/v2"

type Shortcuts struct {
	ConfirmKeyboard []ebiten.Key
	ConfirmMouse    []ebiten.MouseButton
	ConfirmGamepad  []ebiten.GamepadButton
}

var Bindings = &Shortcuts{
	ConfirmKeyboard: []ebiten.Key{ebiten.KeyEnter, ebiten.KeyKPEnter},
	ConfirmMouse:    []ebiten.MouseButton{ebiten.MouseButtonLeft},
	ConfirmGamepad:  []ebiten.GamepadButton{ebiten.GamepadButton0},
}
