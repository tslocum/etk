package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Flex is a flexible stack-based layout which may be oriented horizontally or
// vertically. Children are positioned with equal spacing by default. A minimum
// size may instead be specified via SetChildSize, causing children to be
// positioned similar to a flexbox, where each child either has the minimum
// size or the child stretches to fill the remaining row or column.
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
		Box:       NewBox(),
		columnGap: 5,
		rowGap:    5,
	}
}

// SetRect sets the position and size of the widget.
func (f *Flex) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	f.Box.rect = r
	f.modified = true
}

// SetGaps sets the gaps between each child in the Flex.
func (f *Flex) SetGaps(columnGap int, rowGap int) {
	f.Lock()
	defer f.Unlock()

	if f.columnGap == columnGap && f.rowGap == rowGap {
		return
	}

	f.columnGap, f.rowGap = columnGap, rowGap
	f.modified = true
}

// SetChildSize sets the minimum size of each child in the Flex.
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
	r := f.rect
	childWidth := f.childWidth
	if childWidth == 0 {
		if f.vertical {
			childWidth = r.Dx()
		} else if len(f.children) > 0 {
			var gapSpace int
			if len(f.children) > 1 {
				gapSpace = f.columnGap * (len(f.children) - 1)
			}
			childWidth = (r.Dx() - gapSpace) / len(f.children)
		}
	} else if childWidth > r.Dx() {
		childWidth = r.Dx()
	}
	childHeight := f.childHeight
	if childHeight == 0 {
		if f.vertical && len(f.children) > 0 {
			var gapSpace int
			if len(f.children) > 1 {
				gapSpace = f.rowGap * (len(f.children) - 1)
			}
			childHeight = (r.Dy() - gapSpace) / len(f.children)
		} else {
			childHeight = r.Dy()
		}
	} else if childHeight > r.Dy() {
		childHeight = r.Dy()
	}

	rects := make([]image.Rectangle, len(f.children))
	x1, y1 := r.Min.X, r.Min.Y
	if f.vertical {
		for i := range f.children {
			x2, y2 := x1+childWidth, y1+childHeight
			if y2 > r.Max.Y {
				return
			}
			rects[i] = image.Rect(x1, y1, x2, y2)

			y1 += childHeight + f.rowGap
			if y1 > r.Max.Y-childHeight {
				rects[i].Max.Y = r.Max.Y
				x1 += childWidth + f.columnGap
				y1 = r.Min.Y
			}
		}
	} else {
		for i := range f.children {
			x2, y2 := x1+childWidth, y1+childHeight
			if x2 > r.Max.X {
				return
			}
			rects[i] = image.Rect(x1, y1, x2, y2)

			x1 += childWidth + f.columnGap
			if x1 > r.Max.X-childWidth {
				rects[i].Max.X = r.Max.X
				y1 += childHeight + f.rowGap
				x1 = r.Min.X
			}
		}
	}
	for i, child := range f.children {
		child.SetRect(rects[i])
	}
}
