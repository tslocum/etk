package kibodo

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Key represents a virtual key.
type Key struct {
	LowerLabel string
	UpperLabel string
	LowerInput *Input
	UpperInput *Input
	Wide       bool

	x, y int
	w, h int

	pressed        bool
	pressedTime    time.Time
	pressedTouchID ebiten.TouchID
	repeatTime     time.Time
}

// Input represents the input event from a key press.
type Input struct {
	Rune rune
	Key  ebiten.Key
}

func (i *Input) String() string {
	if i.Rune > 0 {
		return string(i.Rune)
	}
	return i.Key.String()
}
