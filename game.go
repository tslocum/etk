package etk

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var root Widget

var (
	lastWidth, lastHeight int
)

func SetRoot(w Widget) {
	root = w
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

	var mouseHandled bool
	var keyboardHandled bool
	var err error

	children := root.Children()
	for _, child := range children {
		if !mouseHandled {
			mouseHandled, err = child.HandleMouse()
			if err != nil {
				return fmt.Errorf("failed to handle widget mouse input: %s", err)
			}
		}
		if !keyboardHandled {
			keyboardHandled, err = child.HandleKeyboard()
			if err != nil {
				return fmt.Errorf("failed to handle widget keyboard input: %s", err)
			}
		}
		if mouseHandled && keyboardHandled {
			return nil
		}
	}
	if !mouseHandled {
		_, err = root.HandleMouse()
		if err != nil {
			return fmt.Errorf("failed to handle widget mouse input: %s", err)
		}
	}
	if !keyboardHandled {
		_, err = root.HandleKeyboard()
		if err != nil {
			return fmt.Errorf("failed to handle widget keyboard input: %s", err)
		}
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	if root == nil {
		panic("no root widget specified")
	}

	err := root.Draw(screen)
	if err != nil {
		return fmt.Errorf("failed to draw widget: %s", err)
	}

	children := root.Children()
	for _, child := range children {
		err = child.Draw(screen)
		if err != nil {
			return fmt.Errorf("failed to draw widget: %s", err)
		}
	}
	return nil
}
