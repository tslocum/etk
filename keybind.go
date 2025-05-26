package etk

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Shortcuts represents the keyboard, mouse and gamepad input configurations.
type Shortcuts struct {
	DoubleClickThreshold time.Duration

	MoveLeftKeyboard  []ebiten.Key
	MoveRightKeyboard []ebiten.Key
	MoveDownKeyboard  []ebiten.Key
	MoveUpKeyboard    []ebiten.Key

	MoveLeftGamepad  []ebiten.StandardGamepadButton
	MoveRightGamepad []ebiten.StandardGamepadButton
	MoveDownGamepad  []ebiten.StandardGamepadButton
	MoveUpGamepad    []ebiten.StandardGamepadButton

	ConfirmKeyboard []ebiten.Key
	ConfirmMouse    []ebiten.MouseButton
	ConfirmGamepad  []ebiten.StandardGamepadButton

	// A sentinel rune value may be set for the confirm and back actions.
	// This allows working around on-screen keyboard issues on Android.
	ConfirmRune rune
	BackRune    rune
}

// Bindings is the current keyboard, mouse and gamepad input configurations.
var Bindings = &Shortcuts{
	DoubleClickThreshold: 500 * time.Millisecond,

	MoveLeftKeyboard:  []ebiten.Key{ebiten.KeyLeft},
	MoveRightKeyboard: []ebiten.Key{ebiten.KeyRight},
	MoveDownKeyboard:  []ebiten.Key{ebiten.KeyDown},
	MoveUpKeyboard:    []ebiten.Key{ebiten.KeyUp},

	MoveLeftGamepad:  []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftLeft},
	MoveRightGamepad: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftRight},
	MoveDownGamepad:  []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftBottom},
	MoveUpGamepad:    []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftTop},

	ConfirmKeyboard: []ebiten.Key{ebiten.KeyEnter, ebiten.KeyKPEnter},
	ConfirmMouse:    []ebiten.MouseButton{ebiten.MouseButtonLeft, ebiten.MouseButtonRight},
	ConfirmGamepad:  []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonRightBottom},
}
