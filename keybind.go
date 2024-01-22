package etk

import "github.com/hajimehoshi/ebiten/v2"

// Shortcuts represents a keyboard shortcut configuration.
type Shortcuts struct {
	ConfirmKeyboard []ebiten.Key
	ConfirmMouse    []ebiten.MouseButton
	ConfirmGamepad  []ebiten.GamepadButton

	// A sentinel rune value may be set for the confirm and back actions.
	// This allows working around on-screen keyboard issues on Android.
	ConfirmRune rune
	BackRune    rune
}

// Bindings is the current keyboard shortcut configuration.
var Bindings = &Shortcuts{
	ConfirmKeyboard: []ebiten.Key{ebiten.KeyEnter, ebiten.KeyKPEnter},
	ConfirmMouse:    []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight},
	ConfirmGamepad:  []ebiten.GamepadButton{ebiten.GamepadButton0},
}
