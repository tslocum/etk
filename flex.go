package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Flex is a flexbox layou which may be oriented horizontally or vertically.
type Flex struct {
	*Box
	vertical                bool
	childWidth, childHeight int
	columnGap, rowGap       int
	modified                bool
}

// NewFlex returns a new Flex widget.
func NewFlex() *Flex {
	return &Flex{
		Box: NewBox(),
	}
}

// SetRect sets the position and size of the widget.
func (f *Flex) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	f.Box.rect = r
	f.modified = true
}

// SetGapSize sets the gap between child in the Flex.
func (f *Flex) SetGapSize(columnGap int, rowGap int) {
	f.Lock()
	defer f.Unlock()

	if f.columnGap == columnGap && f.rowGap == rowGap {
		return
	}

	f.columnGap, f.rowGap = columnGap, rowGap
	f.modified = true
}

// SetChildSize sets the size of each child in the Flex.
func (f *Flex) SetChildSize(width int, height int) {
	f.Lock()
	defer f.Unlock()

	if f.childWidth == width && f.childHeight == height {
		return
	}

	f.childWidth, f.childHeight = width, height
	f.modified = true
}

// SetVertical sets the orientation of the child widget stacking.
func (f *Flex) SetVertical(v bool) {
	f.Lock()
	defer f.Unlock()

	if f.vertical == v {
		return
	}

	f.vertical = v
	f.modified = true
}

// AddChild adds a child to the widget.
func (f *Flex) AddChild(w ...Widget) {
	f.Lock()
	defer f.Unlock()

	f.children = append(f.children, w...)
	f.modified = true
}

// HandleKeyboard is called when a keyboard event occurs.
func (f *Flex) HandleKeyboard(ebiten.Key, rune) (handled bool, err error) {
	return false, nil
}

// HandleMouse is called when a mouse event occurs.
func (f *Flex) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

// Draw draws the widget on the screen.
func (f *Flex) Draw(screen *ebiten.Image) error {
	f.Lock()
	defer f.Unlock()

	if f.modified {
		f.reposition()
		f.modified = false
	}

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

	// flexbox
	// gap
	// orientation

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
