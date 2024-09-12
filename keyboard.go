package etk

import (
	"image"

	"code.rocket9labs.com/tslocum/etk/kibodo"
	"github.com/hajimehoshi/ebiten/v2"
)

// Keyboard is an on-screen keyboard widget. User input is automatically passed
// to the focused widget.
type Keyboard struct {
	*Box
	k        *kibodo.Keyboard
	incoming []*kibodo.Input
}

// NewKeyboard returns a new Keyboard widget.
func NewKeyboard() *Keyboard {
	k := kibodo.NewKeyboard(Style.TextFont)
	k.Show()
	return &Keyboard{
		Box: NewBox(),
		k:   k,
	}
}

// SetRect sets the position and size of the keyboard.
func (k *Keyboard) SetRect(r image.Rectangle) {
	k.Lock()
	defer k.Unlock()
	k.rect = r
	k.k.SetRect(r.Min.X, r.Min.Y, r.Dx(), r.Dy())
}

// Visible returns the visibility of the keyboard.
func (k *Keyboard) Visible() bool {
	k.Lock()
	defer k.Unlock()
	return k.visible && k.k.Visible()
}

// SetVisible sets the visibility of the keyboard.
func (k *Keyboard) SetVisible(visible bool) {
	k.Lock()
	defer k.Unlock()
	k.visible = visible
	if visible {
		k.k.Show()
	} else {
		k.k.Hide()
	}
}

// Keys returns the keys of the keyboard.
func (k *Keyboard) Keys() [][]*kibodo.Key {
	k.Lock()
	defer k.Unlock()
	return k.k.GetKeys()
}

// SetKeys sets the keys of the keyboard.
func (k *Keyboard) SetKeys(keys [][]*kibodo.Key) {
	k.Lock()
	defer k.Unlock()
	k.k.SetKeys(keys)
}

// SetExtendedKeys sets the keys of the keyboard when the .
func (k *Keyboard) SetExtendedKeys(keys [][]*kibodo.Key) {
	k.Lock()
	defer k.Unlock()
	k.k.SetExtendedKeys(keys)
}

// SetShowExtended sets whether the normal or extended keyboard is shown.
func (k *Keyboard) SetShowExtended(show bool) {
	k.Lock()
	defer k.Unlock()
	k.k.SetShowExtended(show)
}

// SetScheduleFrameFunc sets the function called whenever the screen should be redrawn.
func (k *Keyboard) SetScheduleFrameFunc(f func()) {
	k.Lock()
	defer k.Unlock()
	k.k.SetScheduleFrameFunc(f)
}

// Cursor returns the cursor shape shown when a mouse cursor hovers over the
// widget, or -1 to let widgets beneath determine the cursor shape.
func (k *Keyboard) Cursor() ebiten.CursorShapeType {
	return ebiten.CursorShapePointer
}

// HandleMouse is called when a mouse event occurs.
func (k *Keyboard) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	k.Lock()
	defer k.Unlock()
	return k.k.HandleMouse(cursor, pressed, clicked)
}

// Draw draws the keyboard on the screen.
func (k *Keyboard) Draw(screen *ebiten.Image) error {
	k.Lock()
	defer k.Unlock()
	k.incoming = k.k.AppendInput(k.incoming[:0])
	w := Focused()
	if w != nil {
		for _, key := range k.incoming {
			_, err := w.HandleKeyboard(key.Key, key.Rune)
			if err != nil {
				return err
			}
		}
	}
	k.k.Draw(screen)
	return nil
}
