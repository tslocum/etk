package etk

import (
	"fmt"
	"image"
	"math"
	"runtime/debug"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var root Widget

var (
	lastWidth, lastHeight int

	lastX, lastY = -math.MaxInt, -math.MaxInt

	touchIDs []ebiten.TouchID

	focusedWidget Widget

	lastBackspaceRepeat time.Time

	keyBuffer  []ebiten.Key
	runeBuffer []rune
)

const (
	backspaceRepeatWait = 500 * time.Millisecond
	backspaceRepeatTime = 75 * time.Millisecond
)

// SetRoot sets the root widget. The root widget and all of its children will
// be drawn on the screen and receive user input. The root widget will also be
// focused. To temporarily disable etk, set a nil root widget.
func SetRoot(w Widget) {
	root = w
	if root != nil && (lastWidth != 0 || lastHeight != 0) {
		root.SetRect(image.Rect(0, 0, lastWidth, lastHeight))
	}
	SetFocus(root)
}

// SetFocus focuses a widget.
func SetFocus(w Widget) {
	lastFocused := focusedWidget
	if w != nil && !w.SetFocus(true) {
		return
	}
	if lastFocused != nil && lastFocused != w {
		lastFocused.SetFocus(false)
	}
	focusedWidget = w
}

// Focused returns the currently focused widget. If no widget is focused, nil is returned.
func Focused() Widget {
	return focusedWidget
}

func boundString(f font.Face, s string) (bounds fixed.Rectangle26_6, advance fixed.Int26_6) {
	if strings.TrimSpace(s) == "" {
		return fixed.Rectangle26_6{}, 0
	}
	for i := 0; i < 100; i++ {
		bounds, advance = func() (fixed.Rectangle26_6, fixed.Int26_6) {
			defer func() {
				err := recover()
				if err != nil && i == 99 {
					debug.PrintStack()
					panic("failed to calculate bounds of string '" + s + "'")
				}
			}()
			bounds, advance = font.BoundString(f, s)
			return bounds, advance
		}()
		if !bounds.Empty() {
			return bounds, advance
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fixed.Rectangle26_6{}, 0
}

func int26ToRect(r fixed.Rectangle26_6) image.Rectangle {
	x, y := r.Min.X, r.Min.Y
	w, h := r.Max.X-r.Min.X, r.Max.Y-r.Min.Y
	return image.Rect(x.Round(), y.Round(), (x + w).Round(), (y + h).Round())
}

// BoundString returns the bounds of the provided string.
func BoundString(f font.Face, s string) image.Rectangle {
	bounds, _ := boundString(f, s)
	return int26ToRect(bounds)
}

// Layout sets the current screen size and resizes the root widget.
func Layout(outsideWidth, outsideHeight int) {
	if outsideWidth != lastWidth || outsideHeight != lastHeight {
		lastWidth, lastHeight = outsideWidth, outsideHeight
	}

	if root == nil {
		return
	}
	root.SetRect(image.Rect(0, 0, outsideWidth, outsideHeight))
}

// Update handles user input and passes it to the focused or clicked widget.
func Update() error {
	if root == nil {
		return nil
	}

	var cursor image.Point

	// Handle touch input.

	var pressed bool
	var clicked bool
	var touchInput bool

	touchIDs = inpututil.AppendJustPressedTouchIDs(touchIDs[:0])
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		if x != 0 || y != 0 {
			cursor = image.Point{x, y}

			pressed = true
			clicked = true
			touchInput = true
		}
	}

	// Handle mouse input.

	if !touchInput {
		x, y := ebiten.CursorPosition()
		cursor = image.Point{x, y}

		if lastX == -math.MaxInt && lastY == -math.MaxInt {
			lastX, lastY = x, y
		}
		for _, binding := range Bindings.ConfirmMouse {
			pressed = ebiten.IsMouseButtonPressed(binding)
			if pressed {
				break
			}
		}

		for _, binding := range Bindings.ConfirmMouse {
			clicked = inpututil.IsMouseButtonJustReleased(binding)
			if clicked {
				break
			}
		}
	}

	_, err := update(root, cursor, pressed, clicked, false)
	if err != nil {
		return fmt.Errorf("failed to handle widget mouse input: %s", err)
	}

	// Handle keyboard input.

	if focusedWidget != nil {
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
				lastBackspaceRepeat = time.Now().Add(backspaceRepeatWait)
			} else if time.Since(lastBackspaceRepeat) >= backspaceRepeatTime {
				lastBackspaceRepeat = time.Now()

				_, err := focusedWidget.HandleKeyboard(ebiten.KeyBackspace, 0)
				if err != nil {
					return err
				}
			}
		}

		keyBuffer = inpututil.AppendJustPressedKeys(keyBuffer[:0])
		for _, key := range keyBuffer {
			_, err := focusedWidget.HandleKeyboard(key, 0)
			if err != nil {
				return fmt.Errorf("failed to handle widget keyboard input: %s", err)
			}
		}

		runeBuffer = ebiten.AppendInputChars(runeBuffer[:0])
		for _, r := range runeBuffer {
			_, err := focusedWidget.HandleKeyboard(-1, r)
			if err != nil {
				return fmt.Errorf("failed to handle widget keyboard input: %s", err)
			}
		}
	}
	return nil
}

func at(w Widget, p image.Point) Widget {
	if w == nil {
		return nil
	}

	if !p.In(w.Rect()) {
		return nil
	}

	for _, child := range w.Children() {
		if !child.Visible() {
			continue
		}

		if p.In(child.Rect()) {
			result := at(child, p)
			if result != nil {
				return result
			}
		}
	}

	return w
}

// At returns the widget at the provided screen location.
func At(p image.Point) Widget {
	return at(root, p)
}

func update(w Widget, cursor image.Point, pressed bool, clicked bool, mouseHandled bool) (bool, error) {
	if w == nil {
		return false, nil
	}

	if !w.Visible() {
		return mouseHandled, nil
	}

	var err error
	children := w.Children()
	for i := len(children) - 1; i >= 0; i-- {
		mouseHandled, err = update(children[i], cursor, pressed, clicked, mouseHandled)
		if err != nil {
			return false, err
		} else if mouseHandled {
			return true, nil
		}
	}
	if !mouseHandled && cursor.In(w.Rect()) {
		mouseHandled, err = w.HandleMouse(cursor, pressed, clicked)
		if err != nil {
			return false, fmt.Errorf("failed to handle widget mouse input: %s", err)
		}
		if clicked && mouseHandled {
			SetFocus(w)
		}
	}
	return mouseHandled, nil
}

// Draw draws the root widget and its children to the screen.
func Draw(screen *ebiten.Image) error {
	return draw(root, screen)
}

func draw(w Widget, screen *ebiten.Image) error {
	if w == nil {
		return nil
	}

	if !w.Visible() {
		return nil
	}

	background := w.Background()
	if background.A > 0 {
		screen.SubImage(w.Rect()).(*ebiten.Image).Fill(background)
	}

	err := w.Draw(screen)
	if err != nil {
		return fmt.Errorf("failed to draw widget: %s", err)
	}

	children := w.Children()
	for _, child := range children {
		err = draw(child, screen)
		if err != nil {
			return fmt.Errorf("failed to draw widget: %s", err)
		}
	}

	return nil
}

func rectAtOrigin(r image.Rectangle) image.Rectangle {
	r.Max.X, r.Max.Y = r.Dx(), r.Dy()
	r.Min.X, r.Min.Y = 0, 0
	return r
}
