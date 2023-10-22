package etk

import (
	"fmt"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var root Widget

var (
	lastWidth, lastHeight int

	lastX, lastY = -math.MaxInt, -math.MaxInt

	touchIDs []ebiten.TouchID

	focusedWidget Widget
)

func SetRoot(w Widget) {
	root = w
	if lastWidth != 0 || lastHeight != 0 {
		root.SetRect(image.Rect(0, 0, lastWidth, lastHeight))
	}
	SetFocus(root)
}

func SetFocus(w Widget) {
	if focusedWidget != nil {
		focusedWidget.SetFocus(false)
	}
	focusedWidget = w
	if w != nil {
		w.SetFocus(true)
	}
}

func Layout(outsideWidth, outsideHeight int) {
	if root == nil {
		panic("no root widget specified")
	}

	if outsideWidth != lastWidth || outsideHeight != lastHeight {
		root.SetRect(image.Rect(0, 0, outsideWidth, outsideHeight))
		lastWidth, lastHeight = outsideWidth, outsideHeight
	}
}

func Update() error {
	if root == nil {
		panic("no root widget specified")
	}

	var cursor image.Point

	// Handle touch input.

	var touchInput bool

	var clicked bool
	touchIDs = inpututil.AppendJustPressedTouchIDs(touchIDs[:0])
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		if x != 0 || y != 0 {
			cursor = image.Point{x, y}
			clicked = true
			touchInput = true
		}
	}

	// Handle mouse input.

	var pressed bool
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

	_, _, err := update(root, cursor, pressed, clicked, false, false)
	return err
}

func getWidgetAt(w Widget, cursor image.Point) Widget {
	if !cursor.In(w.Rect()) {
		return nil
	}
	for _, child := range w.Children() {
		if cursor.In(child.Rect()) {
			result := getWidgetAt(child, cursor)
			if result != nil {
				return result
			}
		}
	}
	return w
}

func update(w Widget, cursor image.Point, pressed bool, clicked bool, mouseHandled bool, keyboardHandled bool) (bool, bool, error) {
	var err error
	children := w.Children()
	for _, child := range children {
		mouseHandled, keyboardHandled, err = update(child, cursor, pressed, clicked, mouseHandled, keyboardHandled)
		if err != nil {
			return false, false, err
		} else if mouseHandled && keyboardHandled {
			return true, true, nil
		}
	}
	if !mouseHandled && cursor.In(w.Rect()) {
		mouseHandled, err = w.HandleMouse(cursor, pressed, clicked)
		if err != nil {
			return false, false, fmt.Errorf("failed to handle widget mouse input: %s", err)
		}
		if clicked && mouseHandled {
			SetFocus(w)
		}
	}
	if !keyboardHandled && w == focusedWidget {
		keyboardHandled, err = w.HandleKeyboard()
		if err != nil {
			return false, false, fmt.Errorf("failed to handle widget keyboard input: %s", err)
		}
	}
	return mouseHandled, keyboardHandled, nil
}

func Draw(screen *ebiten.Image) error {
	if root == nil {
		panic("no root widget specified")
	}

	return draw(root, screen)
}

func draw(w Widget, screen *ebiten.Image) error {
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
