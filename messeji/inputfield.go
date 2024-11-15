package messeji

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// InputField is a text input field. Call Update and Draw when your Game's
// Update and Draw methods are called.
//
// Note: A position and size must be set via SetRect before the field will appear.
// Keyboard events are not handled by default, and may be enabled via SetHandleKeyboard.
type InputField struct {
	*TextField

	// changedFunc is a function which is called when the text buffer is changed.
	// The function may return false to skip adding the rune to the text buffer.
	changedFunc func(r rune) (accept bool)

	// selectedFunc is a function which is called when the enter key is pressed. The
	// function may return true to clear the text buffer.
	selectedFunc func() (accept bool)

	// readBuffer is where incoming runes are stored before being added to the input buffer.
	readBuffer []rune

	// keyBuffer is where incoming keys are stored before being added to the input buffer.
	keyBuffer []ebiten.Key

	// rawRuneBuffer is where incoming raw runes are stored before being added to the input buffer.
	rawRuneBuffer []rune

	// rawKeyBuffer is where incoming raw keys are stored before being added to the input buffer.
	rawKeyBuffer []ebiten.Key

	sync.Mutex
}

// NewInputField returns a new InputField. See type documentation for more info.
func NewInputField(face *text.GoTextFace, faceMutex *sync.Mutex) *InputField {
	f := &InputField{
		TextField: NewTextField(face, faceMutex),
	}
	f.TextField.suffix = "_"
	return f
}

// SetHandleKeyboard sets a flag controlling whether keyboard input should be handled
// by the field. This can be used to facilitate focus changes between multiple inputs.
func (f *InputField) SetHandleKeyboard(handle bool) {
	f.Lock()
	defer f.Unlock()

	f.handleKeyboard = handle
}

// SetChangedFunc sets a handler which is called when the text buffer is changed.
// The handler may return true to add the rune to the text buffer.
func (f *InputField) SetChangedFunc(changedFunc func(r rune) (accept bool)) {
	f.changedFunc = changedFunc
}

// SetSelectedFunc sets a handler which is called when the enter key is pressed.
// Providing a nil function value will remove the existing handler (if set).
// The handler may return true to clear the text buffer.
func (f *InputField) SetSelectedFunc(selectedFunc func() (accept bool)) {
	f.selectedFunc = selectedFunc
}

// HandleKeyboardEvent passes the provided key or rune to the Inputfield.
func (f *InputField) HandleKeyboardEvent(key ebiten.Key, r rune) (handled bool, err error) {
	f.Lock()
	defer f.Unlock()

	if !f.visible || rectIsZero(f.r) {
		return
	}

	if !f.handleKeyboard {
		return
	}

	// Handle rune event.
	if r > 0 {
		f.handleRunes([]rune{r})
		return true, nil
	}

	// Handle key event.
	f.handleKeys([]ebiten.Key{key})
	return true, nil
}

func (f *InputField) handleRunes(runes []rune) bool {
	var redraw bool
	for _, r := range runes {
		if f.changedFunc != nil {
			f.Unlock()
			accept := f.changedFunc(r)
			f.Lock()

			if !accept {
				continue
			}
		}

		f.TextField._write([]byte(string(r)))
		redraw = true
	}

	return redraw
}

func (f *InputField) handleKeys(keys []ebiten.Key) bool {
	var redraw bool
	for _, key := range keys {
		switch key {
		case ebiten.KeyBackspace:
			l := len(f.buffer)
			if l > 0 {
				var rewrap bool
				if len(f.incoming) != 0 {
					line := string(f.incoming)
					f.incoming = append(f.incoming, []byte(line[:len(line)-1])...)
				} else if len(f.buffer[l-1]) == 0 {
					f.buffer = f.buffer[:l-1]
					rewrap = true
				} else {
					line := string(f.buffer[l-1])
					f.buffer[l-1] = []byte(line[:len(line)-1])
					rewrap = true
				}
				if rewrap && (f.needWrap == -1 || f.needWrap > l-1) {
					f.needWrap = l - 1
				}
				redraw = true
				f.modified = true
				f.redraw = true
			}
		case ebiten.KeyEnter, ebiten.KeyKPEnter:
			if f.selectedFunc != nil {
				f.Unlock()
				accept := f.selectedFunc()
				f.Lock()

				// Clear input buffer.
				if accept {
					f.incoming = f.incoming[:0]
					f.buffer = f.buffer[:0]
					f.bufferWrapped = f.bufferWrapped[:0]
					f.lineWidths = f.lineWidths[:0]
					f.needWrap = 0
					f.wrapStart = 0
					f.modified = true
					f.redraw = true
					redraw = true
				}
			} else if !f.singleLine {
				// Append newline.
				f.incoming = append(f.incoming, '\n')
				f.modified = true
				f.redraw = true
				redraw = true
			}
		}
	}
	return redraw
}

// Update updates the input field. This function should be called when
// Game.Update is called.
func (f *InputField) Update() error {
	f.Lock()
	defer f.Unlock()

	if !f.visible || rectIsZero(f.r) {
		return nil
	}

	if !f.handleKeyboard {
		return f.TextField.Update()
	}

	var redraw bool

	// Handler rune input.
	f.readBuffer = ebiten.AppendInputChars(f.readBuffer[:0])
	if f.handleRunes(f.readBuffer) {
		redraw = true
	}
	if f.handleRunes(f.rawRuneBuffer) {
		redraw = true
	}
	f.rawRuneBuffer = f.rawRuneBuffer[:0]

	// Handle key input.
	f.keyBuffer = inpututil.AppendJustPressedKeys(f.keyBuffer[:0])
	if f.handleKeys(f.keyBuffer) {
		redraw = true
	}
	if f.handleKeys(f.rawKeyBuffer) {
		redraw = true
	}
	f.rawKeyBuffer = f.rawKeyBuffer[:0]

	if redraw {
		f.bufferModified()
	}

	return f.TextField.Update()
}
