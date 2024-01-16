package kibodo

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Keyboard is an on-screen keyboard widget.
type Keyboard struct {
	x, y int
	w, h int

	visible       bool
	alpha         float64
	passPhysical  bool
	allowUserHide bool

	incomingBuffer []rune

	inputEvents []*Input

	keys         [][]*Key
	normalKeys   [][]*Key
	extendedKeys [][]*Key
	showExtended bool

	backgroundLower *ebiten.Image
	backgroundUpper *ebiten.Image
	backgroundDirty bool

	op *ebiten.DrawImageOptions

	backgroundColor     color.Color
	lastBackgroundColor color.Color

	shift bool

	touchIDs    []ebiten.TouchID
	holdTouchID ebiten.TouchID
	holdKey     *Key
	wasPressed  bool

	hideShortcuts []ebiten.Key

	labelFont  font.Face
	lineHeight int
	lineOffset int

	backspaceDelay  time.Duration
	backspaceRepeat time.Duration
	backspaceLast   time.Time

	scheduleFrameFunc func()
}

// NewKeyboard returns a new Keyboard widget.
func NewKeyboard() *Keyboard {
	fontFace, err := defaultFontFace(64)
	if err != nil {
		log.Fatal(err)
	}

	k := &Keyboard{
		alpha: 1.0,
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
		keys:            KeysQWERTY,
		normalKeys:      KeysQWERTY,
		backgroundLower: ebiten.NewImage(1, 1),
		backgroundUpper: ebiten.NewImage(1, 1),
		backgroundColor: color.Black,
		holdTouchID:     -1,
		labelFont:       fontFace,
		backspaceDelay:  500 * time.Millisecond,
		backspaceRepeat: 75 * time.Millisecond,
	}
	k.fontUpdated()
	return k
}

func defaultFont() (*opentype.Font, error) {
	return opentype.Parse(fonts.MPlus1pRegular_ttf)
}

func defaultFontFace(size float64) (font.Face, error) {
	f, err := defaultFont()
	if err != nil {
		return nil, err
	}
	const dpi = 72 // TODO
	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

// SetRect sets the position and size of the widget.
func (k *Keyboard) SetRect(x, y, w, h int) {
	if k.x == x && k.y == y && k.w == w && k.h == h {
		return
	}
	k.x, k.y, k.w, k.h = x, y, w, h

	k.updateKeyRects()
	k.backgroundDirty = true
}

// Rect returns the position and size of the widget.
func (k *Keyboard) Rect() image.Rectangle {
	return image.Rect(k.x, k.y, k.x+k.w, k.y+k.h)
}

// GetKeys returns the keys of the keyboard.
func (k *Keyboard) GetKeys() [][]*Key {
	return k.keys
}

// SetKeys sets the keys of the keyboard.
func (k *Keyboard) SetKeys(keys [][]*Key) {
	k.normalKeys = keys

	if !k.showExtended && !keysEqual(keys, k.keys) {
		k.keys = keys
		k.updateKeyRects()
		k.backgroundDirty = true
	}
}

// SetExtendedKeys sets the keys of the keyboard when the .
func (k *Keyboard) SetExtendedKeys(keys [][]*Key) {
	k.extendedKeys = keys

	if k.showExtended && !keysEqual(keys, k.keys) {
		k.keys = keys
		k.updateKeyRects()
		k.backgroundDirty = true
	}
}

// SetShowExtended sets whether the normal or extended keyboard is shown.
func (k *Keyboard) SetShowExtended(show bool) {
	if k.showExtended == show {
		return
	}
	k.showExtended = show
	if k.showExtended {
		k.keys = k.extendedKeys
	} else {
		k.keys = k.normalKeys
	}
	k.updateKeyRects()
	k.backgroundDirty = true
}

// SetLabelFont sets the key label font.
func (k *Keyboard) SetLabelFont(face font.Face) {
	k.labelFont = face
	k.fontUpdated()

	k.backgroundDirty = true
}

func (k *Keyboard) fontUpdated() {
	m := k.labelFont.Metrics()
	k.lineHeight = m.Height.Round()
	k.lineOffset = m.Ascent.Round()
}

// SetHideShortcuts sets the key shortcuts which, when pressed, will hide the
// keyboard.
func (k *Keyboard) SetHideShortcuts(shortcuts []ebiten.Key) {
	k.hideShortcuts = shortcuts
}

func (k *Keyboard) updateKeyRects() {
	if len(k.keys) == 0 {
		return
	}

	maxCells := 0
	for _, rowKeys := range k.keys {
		if len(rowKeys) > maxCells {
			maxCells = len(rowKeys)
		}
	}

	// TODO user configurable
	cellPaddingW := 1
	cellPaddingH := 1

	cellH := (k.h - (cellPaddingH * (len(k.keys) - 1))) / len(k.keys)

	row := 0
	x, y := 0, 0
	for _, rowKeys := range k.keys {
		if len(rowKeys) == 0 {
			continue
		}

		availableWidth := k.w
		for _, key := range rowKeys {
			if key.Wide {
				availableWidth = availableWidth / 2
				break
			}
		}

		cellW := (availableWidth - (cellPaddingW * (len(rowKeys) - 1))) / len(rowKeys)

		x = 0
		for i, key := range rowKeys {
			key.w, key.h = cellW, cellH
			key.x, key.y = x, y

			if i == len(rowKeys)-1 {
				key.w = k.w - key.x
			}

			if key.Wide {
				key.w = k.w - k.w/2 + (cellW)
			}

			x += key.w
		}

		// Count non-empty rows only
		row++
		y += (cellH + cellPaddingH)
	}
}

func (k *Keyboard) at(x, y int) *Key {
	if !k.visible {
		return nil
	}
	if x >= k.x && x <= k.x+k.w && y >= k.y && y <= k.y+k.h {
		x, y = x-k.x, y-k.y // Offset
		for _, rowKeys := range k.keys {
			for _, key := range rowKeys {
				if x >= key.x && x <= key.x+key.w && y >= key.y && y <= key.y+key.h {
					return key
				}
			}
		}
	}
	return nil
}

// KeyAt returns the key located at the specified position, or nil if no key is found.
func (k *Keyboard) KeyAt(x, y int) *Key {
	return k.at(x, y)
}

func (k *Keyboard) handleToggleExtendedKey(inputKey ebiten.Key) bool {
	if inputKey != KeyToggleExtended {
		return false
	}
	k.showExtended = !k.showExtended
	if k.showExtended {
		k.keys = k.extendedKeys
	} else {
		k.keys = k.normalKeys
	}
	k.updateKeyRects()
	k.backgroundDirty = true
	return true
}

func (k *Keyboard) handleHideKey(inputKey ebiten.Key) bool {
	if !k.allowUserHide {
		return false
	}

	for _, key := range k.hideShortcuts {
		if key == inputKey {
			k.Hide()
			return true
		}
	}
	return false
}

// Hit handles a key press.
func (k *Keyboard) Hit(key *Key) {
	input := key.LowerInput
	if k.shift {
		input = key.UpperInput
	}

	if input.Key == ebiten.KeyShift {
		k.shift = !k.shift
		if k.scheduleFrameFunc != nil {
			k.scheduleFrameFunc()
		}
		return
	} else if k.handleToggleExtendedKey(input.Key) || k.handleHideKey(input.Key) {
		return
	}

	k.inputEvents = append(k.inputEvents, input)
}

// HandleMouse passes the specified mouse event to the on-screen keyboard.
func (k *Keyboard) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	if k.backgroundDirty {
		k.drawBackground()
		k.backgroundDirty = false
	}

	pressDuration := 50 * time.Millisecond
	if k.wasPressed && !pressed && !clicked {
		var key *Key
		if cursor.X != 0 || cursor.Y != 0 {
			key = k.at(cursor.X, cursor.Y)
		} else {
		PRESSKEY:
			for _, rowKeys := range k.keys {
				for _, rowKey := range rowKeys {
					if rowKey.pressed {
						key = rowKey
						break PRESSKEY
					}
				}
			}
		}
		for _, rowKeys := range k.keys {
			for _, rowKey := range rowKeys {
				if key != nil && rowKey == key {
					continue
				}
				rowKey.pressed = false
			}
		}
		if key != nil {
			key.pressed = true

			k.Hit(key)

			go func() {
				time.Sleep(pressDuration)

				key.pressed = false
				if k.scheduleFrameFunc != nil {
					k.scheduleFrameFunc()
				}
			}()
		}
		k.wasPressed = false
	} else if pressed {
		key := k.at(cursor.X, cursor.Y)
		if key != nil {
			if !key.pressed {
				input := key.LowerInput
				if k.shift {
					input = key.UpperInput
				}

				// Repeat backspace and delete operations.
				if input.Key == ebiten.KeyBackspace || input.Key == ebiten.KeyDelete {
					k.backspaceLast = time.Now().Add(k.backspaceDelay)
				}
				go func() {
					t := time.NewTicker(k.backspaceRepeat)
					for {
						<-t.C

						if !key.pressed {
							t.Stop()
							return
						}

						if (input.Key == ebiten.KeyBackspace || input.Key == ebiten.KeyDelete) && time.Since(k.backspaceLast) >= k.backspaceRepeat {
							k.backspaceLast = time.Now()
							k.inputEvents = append(k.inputEvents, &Input{Key: input.Key})
						}
					}

				}()
			}
			key.pressed = true
			k.wasPressed = true

			for _, rowKeys := range k.keys {
				for _, rowKey := range rowKeys {
					if rowKey == key || !rowKey.pressed {
						continue
					}
					rowKey.pressed = false
				}
			}
		}
	}
	return true, nil
}

// Update handles user input. This function is called by Ebitengine.
func (k *Keyboard) Update() error {
	if !k.visible {
		return nil
	}

	if k.backgroundDirty {
		k.drawBackground()
		k.backgroundDirty = false
	}

	// Pass through physical keyboard input
	if k.passPhysical {
		// Read input characters
		k.incomingBuffer = ebiten.AppendInputChars(k.incomingBuffer[:0])
		if len(k.incomingBuffer) > 0 {
			for _, r := range k.incomingBuffer {
				k.inputEvents = append(k.inputEvents, &Input{Rune: r}) // Pass through
			}
		} else {
			// Read keys
			for _, key := range allKeys {
				if inpututil.IsKeyJustPressed(key) {
					if k.handleHideKey(key) {
						// Hidden
						return nil
					}
					k.inputEvents = append(k.inputEvents, &Input{Key: key}) // Pass through
				}
			}
		}
	}
	// Handle mouse input
	pressDuration := 50 * time.Millisecond
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		key := k.at(x, y)
		if key != nil {
			for _, rowKeys := range k.keys {
				for _, rowKey := range rowKeys {
					rowKey.pressed = false
				}
			}
			key.pressed = true

			k.Hit(key)

			go func() {
				time.Sleep(pressDuration)

				key.pressed = false
				if k.scheduleFrameFunc != nil {
					k.scheduleFrameFunc()
				}
			}()
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		key := k.at(x, y)
		if key != nil {
			if !key.pressed {
				input := key.LowerInput
				if k.shift {
					input = key.UpperInput
				}

				// Repeat backspace and delete operations.
				if input.Key == ebiten.KeyBackspace || input.Key == ebiten.KeyDelete {
					k.backspaceLast = time.Now().Add(k.backspaceDelay)
				}
				go func() {
					t := time.NewTicker(k.backspaceRepeat)
					for {
						<-t.C

						if !key.pressed {
							t.Stop()
							return
						}

						if (input.Key == ebiten.KeyBackspace || input.Key == ebiten.KeyDelete) && time.Since(k.backspaceLast) >= k.backspaceRepeat {
							k.backspaceLast = time.Now()
							k.inputEvents = append(k.inputEvents, &Input{Key: input.Key})
						}
					}

				}()
			}
			key.pressed = true

			for _, rowKeys := range k.keys {
				for _, rowKey := range rowKeys {
					if rowKey == key || !rowKey.pressed {
						continue
					}
					rowKey.pressed = false
				}
			}
		}
	}
	// Handle touch input
	if k.holdTouchID != -1 {
		x, y := ebiten.TouchPosition(k.holdTouchID)
		if x == 0 && y == 0 {
			k.holdTouchID = -1
		} else {
			key := k.at(x, y)
			if key != k.holdKey {
				k.holdTouchID = -1
				return nil
			}
			//k.Hold(key)
			k.holdKey = key
		}
	}
	if k.holdTouchID == -1 {
		k.touchIDs = inpututil.AppendJustPressedTouchIDs(k.touchIDs[:0])
		for _, id := range k.touchIDs {
			x, y := ebiten.TouchPosition(id)

			key := k.at(x, y)
			if key != nil {
				input := key.LowerInput
				if k.shift {
					input = key.UpperInput
				}

				if !key.pressed {
					key.pressed = true
					key.pressedTouchID = id

					for _, rowKeys := range k.keys {
						for _, rowKey := range rowKeys {
							if rowKey != key && rowKey.pressed {
								rowKey.pressed = false
							}
						}
					}

					k.Hit(key)
					k.holdTouchID = id
					k.holdKey = key

					// Repeat backspace and delete operations.
					if input.Key == ebiten.KeyBackspace || input.Key == ebiten.KeyDelete {
						k.backspaceLast = time.Now().Add(k.backspaceDelay)
					}

					go func() {
						var touchIDs []ebiten.TouchID
						t := time.NewTicker(pressDuration)
						for range t.C {
							touchIDs = ebiten.AppendTouchIDs(touchIDs[:0])

							var found bool
							for _, touchID := range touchIDs {
								if id == touchID {
									found = true
									break
								}
							}

							if found {
								tx, ty := ebiten.TouchPosition(id)
								if tx != 0 || ty != 0 {
									x, y = tx, ty
								}
							}

							if !found {
								key.pressed = false
								if k.scheduleFrameFunc != nil {
									k.scheduleFrameFunc()
								}
								t.Stop()
								return
							}

							// Repeat backspace and delete operations.
							if (input.Key == ebiten.KeyBackspace || input.Key == ebiten.KeyDelete) && time.Since(k.backspaceLast) >= k.backspaceRepeat {
								k.backspaceLast = time.Now()
								k.inputEvents = append(k.inputEvents, &Input{Key: input.Key})
							}
						}
					}()
				}
			}
		}
	}
	return nil
}

func (k *Keyboard) drawBackground() {
	if k.w == 0 || k.h == 0 {
		return
	}

	if !k.backgroundLower.Bounds().Eq(image.Rect(0, 0, k.w, k.h)) || !k.backgroundUpper.Bounds().Eq(image.Rect(0, 0, k.w, k.h)) || k.backgroundColor != k.lastBackgroundColor {
		k.backgroundLower = ebiten.NewImage(k.w, k.h)
		k.backgroundUpper = ebiten.NewImage(k.w, k.h)
		k.lastBackgroundColor = k.backgroundColor
	}
	k.backgroundLower.Fill(k.backgroundColor)
	k.backgroundUpper.Fill(k.backgroundColor)

	halfLineHeight := k.lineHeight / 2

	lightShade := color.RGBA{150, 150, 150, 255}
	darkShade := color.RGBA{30, 30, 30, 255}

	var keyImage *ebiten.Image
	for i := 0; i < 2; i++ {
		shift := i == 1
		img := k.backgroundLower
		if shift {
			img = k.backgroundUpper
		}
		for _, rowKeys := range k.keys {
			for _, key := range rowKeys {
				r := image.Rect(key.x, key.y, key.x+key.w, key.y+key.h)
				keyImage = img.SubImage(r).(*ebiten.Image)

				// Draw key background
				// TODO configurable
				keyImage.Fill(color.RGBA{90, 90, 90, 255})

				// Draw key label
				label := key.LowerLabel
				if shift {
					label = key.UpperLabel
				}

				bounds := text.BoundString(k.labelFont, label)
				x := (key.w - bounds.Dx()) / 2
				if x < 0 {
					x = 0
				}
				y := halfLineHeight + (key.h-halfLineHeight)/2
				text.Draw(keyImage, label, k.labelFont, key.x+x, key.y+y, color.White)

				// Draw border
				keyImage.SubImage(image.Rect(key.x, key.y, key.x+key.w, key.y+1)).(*ebiten.Image).Fill(lightShade)
				keyImage.SubImage(image.Rect(key.x, key.y, key.x+1, key.y+key.h)).(*ebiten.Image).Fill(lightShade)
				keyImage.SubImage(image.Rect(key.x, key.y+key.h-1, key.x+key.w, key.y+key.h)).(*ebiten.Image).Fill(darkShade)
				keyImage.SubImage(image.Rect(key.x+key.w-1, key.y, key.x+key.w, key.y+key.h)).(*ebiten.Image).Fill(darkShade)
			}
		}
	}
}

// Draw draws the widget on the provided image.  This function is called by Ebitengine.
func (k *Keyboard) Draw(target *ebiten.Image) {
	if !k.visible {
		return
	}

	if k.backgroundDirty {
		k.drawBackground()
		k.backgroundDirty = false
	}

	var background *ebiten.Image
	if !k.shift {
		background = k.backgroundLower
	} else {
		background = k.backgroundUpper
	}

	k.op.GeoM.Reset()
	k.op.GeoM.Translate(float64(k.x), float64(k.y))
	k.op.ColorM.Scale(1, 1, 1, k.alpha)
	target.DrawImage(background, k.op)
	k.op.ColorM.Reset()

	// Draw pressed keys
	for _, rowKeys := range k.keys {
		for _, key := range rowKeys {
			if !key.pressed {
				continue
			}

			// TODO buffer to prevent issues with alpha channel
			k.op.GeoM.Reset()
			k.op.GeoM.Translate(float64(k.x+key.x), float64(k.y+key.y))
			k.op.ColorM.Scale(0.75, 0.75, 0.75, k.alpha)

			target.DrawImage(background.SubImage(image.Rect(key.x, key.y, key.x+key.w, key.y+key.h)).(*ebiten.Image), k.op)
			k.op.ColorM.Reset()

			// Draw shadow.
			darkShade := color.RGBA{60, 60, 60, 255}
			subImg := target.SubImage(image.Rect(k.x+key.x, k.y+key.y, k.x+key.x+key.w, k.y+key.y+1)).(*ebiten.Image)
			subImg.Fill(darkShade)
			subImg = target.SubImage(image.Rect(k.x+key.x, k.y+key.y, k.x+key.x+1, k.y+key.y+key.h)).(*ebiten.Image)
			subImg.Fill(darkShade)
		}
	}
}

// SetAllowUserHide sets a flag that controls whether the widget may be hidden
// by the user.
func (k *Keyboard) SetAllowUserHide(allow bool) {
	k.allowUserHide = allow
}

// SetPassThroughPhysicalInput sets a flag that controls whether physical
// keyboard input is passed through to the widget's input buffer. This is not
// enabled by default.
func (k *Keyboard) SetPassThroughPhysicalInput(pass bool) {
	k.passPhysical = pass
}

// SetAlpha sets the transparency level of the widget on a scale of 0 to 1.0.
func (k *Keyboard) SetAlpha(alpha float64) {
	k.alpha = alpha
}

// Show shows the widget.
func (k *Keyboard) Show() {
	k.visible = true
}

// Visible returns whether the widget is currently shown.
func (k *Keyboard) Visible() bool {
	return k.visible
}

// Hide hides the widget.
func (k *Keyboard) Hide() {
	k.visible = false
	if k.showExtended {
		k.showExtended = false
		k.keys = k.normalKeys
		k.updateKeyRects()
		k.backgroundDirty = true
	}
}

// AppendInput appends user input that was received since the function was last called.
func (k *Keyboard) AppendInput(events []*Input) []*Input {
	events = append(events, k.inputEvents...)
	k.inputEvents = nil
	return events
}

// SetScheduleFrameFunc sets the function called whenever the screen should be redrawn.
func (k *Keyboard) SetScheduleFrameFunc(f func()) {
	k.scheduleFrameFunc = f
}

func keysEqual(a [][]*Key, b [][]*Key) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range b[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
