package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Flex struct {
	*Box

	vertical bool
}

func NewFlex() *Flex {
	return &Flex{
		Box: NewBox(),
	}
}

func (f *Flex) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	f.Box.rect = r
	f.reposition()
}

func (f *Flex) SetVertical(v bool) {
	f.Lock()
	defer f.Unlock()

	if f.vertical == v {
		return
	}

	f.vertical = v
	f.reposition()
}

func (f *Flex) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

func (f *Flex) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

func (f *Flex) Draw(screen *ebiten.Image) error {
	f.Lock()
	defer f.Unlock()

	for _, child := range f.children {
		err := child.Draw(screen)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Flex) reposition() {
	l := len(f.children)
	r := f.rect

	if f.vertical {
		childHeight := float64(r.Dy()) / float64(l)

		minY := r.Min.Y
		for i, child := range f.children {
			maxY := r.Min.Y + int(childHeight*float64(i+1))
			if i == l-1 {
				maxY = r.Max.Y
			}
			child.SetRect(image.Rect(r.Min.X, minY, r.Max.X, maxY))

			minY = maxY
		}

		return
	}

	childWidth := float64(r.Dx()) / float64(l)

	minX := r.Min.X
	for i, child := range f.children {
		maxX := r.Min.X + int(childWidth*float64(i+1))
		if i == l-1 {
			maxX = r.Max.X
		}
		child.SetRect(image.Rect(minX, r.Min.Y, maxX, r.Max.Y))

		minX = maxX
	}
}
