package etk

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"sync"
	"time"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Alignment specifies how text is aligned within the field.
type Alignment int

const (
	// AlignStart aligns text at the start of the field.
	AlignStart Alignment = 0

	// AlignCenter aligns text at the center of the field.
	AlignCenter Alignment = 1

	// AlignEnd aligns text at the end of the field.
	AlignEnd Alignment = 2
)

var root Widget

var drawDebug bool

var (
	lastWidth, lastHeight int

	lastX, lastY = -math.MaxInt, -math.MaxInt

	touchIDs      []ebiten.TouchID
	activeTouchID = ebiten.TouchID(-1)

	focusedWidget Widget

	pressedWidget Widget

	cursorShape ebiten.CursorShapeType

	lastBackspaceRepeat time.Time

	keyBuffer  []ebiten.Key
	runeBuffer []rune

	fontMutex = &sync.Mutex{}
)

var debugColor = color.RGBA{0, 0, 255, 255}

const (
	backspaceRepeatWait = 500 * time.Millisecond
	backspaceRepeatTime = 75 * time.Millisecond
)

var deviceScale float64

// ScaleFactor returns the device scale factor. When running on Android, this function
// may only be called during or after the first Layout call made by Ebitengine.
func ScaleFactor() float64 {
	if deviceScale == 0 {
		monitor := ebiten.Monitor()
		if monitor != nil {
			deviceScale = monitor.DeviceScaleFactor()
		}
		if deviceScale <= 0 {
			deviceScale = ebiten.DeviceScaleFactor()
		}

	}
	return deviceScale
}

// Scale applies the device scale factor to the provided value and returns the result.
// When running on Android, this function may only be called during or after the first
// Layout call made by Ebitengine.
func Scale(v int) int {
	if deviceScale == 0 {
		monitor := ebiten.Monitor()
		if monitor != nil {
			deviceScale = monitor.DeviceScaleFactor()
		}
		if deviceScale <= 0 {
			deviceScale = ebiten.DeviceScaleFactor()
		}

	}
	return int(float64(v) * deviceScale)
}

var (
	fontCache     = make(map[string]font.Face)
	fontCacheLock sync.Mutex
)

// FontFace returns a face for the provided font and size. Scaling is not applied.
func FontFace(fnt *sfnt.Font, size int) font.Face {
	id := fmt.Sprintf("%p/%d", fnt, size)

	fontCacheLock.Lock()
	defer fontCacheLock.Unlock()

	f := fontCache[id]
	if f != nil {
		return f
	}

	const dpi = 72
	f, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	fontCache[id] = f
	return f
}

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

func int26ToRect(r fixed.Rectangle26_6) image.Rectangle {
	x, y := r.Min.X, r.Min.Y
	w, h := r.Max.X-r.Min.X, r.Max.Y-r.Min.Y
	return image.Rect(x.Round(), y.Round(), (x + w).Round(), (y + h).Round())
}

// BoundString returns the bounds of the provided string.
func BoundString(f font.Face, s string) image.Rectangle {
	fontMutex.Lock()
	defer fontMutex.Unlock()

	bounds, _ := font.BoundString(f, s)
	return int26ToRect(bounds)
}

// SetDebug sets whether debug information is drawn on screen. When enabled,
// all visible widgets are outlined.
func SetDebug(debug bool) {
	drawDebug = debug
}

// ScreenSize returns the current screen size.
func ScreenSize() (width int, height int) {
	return lastWidth, lastHeight
}

// Layout sets the current screen size and resizes the root widget.
func Layout(outsideWidth int, outsideHeight int) {
	outsideWidth, outsideHeight = Scale(outsideWidth), Scale(outsideHeight)
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

	if activeTouchID != -1 {
		x, y := ebiten.TouchPosition(activeTouchID)
		if x != 0 || y != 0 {
			cursor = image.Point{x, y}

			pressed = true
			touchInput = true
		} else {
			activeTouchID = -1
		}
	}

	if activeTouchID == -1 {
		touchIDs = inpututil.AppendJustPressedTouchIDs(touchIDs[:0])
		for _, id := range touchIDs {
			x, y := ebiten.TouchPosition(id)
			if x != 0 || y != 0 {
				cursor = image.Point{x, y}

				pressed = true
				clicked = true
				touchInput = true

				activeTouchID = id
				break
			}
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
			clicked = inpututil.IsMouseButtonJustPressed(binding)
			if clicked {
				break
			}
		}
	}

	if !pressed && !clicked && pressedWidget != nil {
		_, err := pressedWidget.HandleMouse(cursor, false, false)
		if err != nil {
			return err
		}
		pressedWidget = nil
	}

	mouseHandled, err := update(root, cursor, pressed, clicked, false)
	if err != nil {
		return fmt.Errorf("failed to handle widget mouse input: %s", err)
	} else if !mouseHandled && cursorShape != ebiten.CursorShapeDefault {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		cursorShape = ebiten.CursorShapeDefault
	}

	// Handle keyboard input.

	if focusedWidget == nil {
		return nil
	} else if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
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
INPUTCHARS:
	for i, r := range runeBuffer {
		if i > 0 {
			for j, r2 := range runeBuffer {
				if j == i {
					break
				} else if r2 == r {
					continue INPUTCHARS
				}
			}
		}
		var err error
		switch r {
		case Bindings.ConfirmRune:
			_, err = focusedWidget.HandleKeyboard(ebiten.KeyEnter, 0)
		case Bindings.BackRune:
			_, err = focusedWidget.HandleKeyboard(ebiten.KeyBackspace, 0)
		default:
			_, err = focusedWidget.HandleKeyboard(-1, r)
		}
		if err != nil {
			return fmt.Errorf("failed to handle widget keyboard input: %s", err)
		}
	}
	return nil
}

func at(w Widget, p image.Point) Widget {
	if w == nil || !w.Visible() {
		return nil
	}

	for _, child := range w.Children() {
		result := at(child, p)
		if result != nil {
			return result
		}
	}

	if p.In(w.Rect()) {
		return w
	}

	return nil
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
		if pressed && !clicked && w != pressedWidget {
			return mouseHandled, nil
		}
		mouseHandled, err = w.HandleMouse(cursor, pressed, clicked)
		if err != nil {
			return false, fmt.Errorf("failed to handle widget mouse input: %s", err)
		} else if mouseHandled {
			if clicked {
				SetFocus(w)
				pressedWidget = w
			} else if pressedWidget != nil && (!pressed || pressedWidget != w) {
				pressedWidget = nil
			}
			shape := w.Cursor()
			if shape != -1 && shape != cursorShape {
				ebiten.SetCursorShape(shape)
				cursorShape = shape
			}
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

	if drawDebug {
		r := w.Rect()
		if !r.Empty() {
			x, y := r.Min.X, r.Min.Y
			w, h := r.Dx(), r.Dy()
			screen.SubImage(image.Rect(x, y, x+w, y+1)).(*ebiten.Image).Fill(debugColor)
			screen.SubImage(image.Rect(x, y+h-1, x+w, y+h)).(*ebiten.Image).Fill(debugColor)
			screen.SubImage(image.Rect(x, y, x+1, y+h)).(*ebiten.Image).Fill(debugColor)
			screen.SubImage(image.Rect(x+w-1, y, x+w, y+h)).(*ebiten.Image).Fill(debugColor)
		}
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

func newText() *messeji.TextField {
	f := messeji.NewTextField(FontFace(Style.TextFont, Scale(Style.TextSize)), fontMutex)
	f.SetForegroundColor(Style.TextColorLight)
	f.SetBackgroundColor(transparent)
	f.SetScrollBarColors(Style.ScrollAreaColor, Style.ScrollHandleColor)
	f.SetScrollBorderSize(Scale(Style.ScrollBorderSize))
	f.SetScrollBorderColors(Style.ScrollBorderColorTop, Style.ScrollBorderColorRight, Style.ScrollBorderColorBottom, Style.ScrollBorderColorLeft)
	return f
}

func rectAtOrigin(r image.Rectangle) image.Rectangle {
	r.Max.X, r.Max.Y = r.Dx(), r.Dy()
	r.Min.X, r.Min.Y = 0, 0
	return r
}
