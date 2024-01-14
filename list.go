package etk

import (
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// SelectionMode represents a mode of selection.
type SelectionMode int

// Selection modes.
const (
	// SelectNone disables selection.
	SelectNone SelectionMode = iota

	// SelectRow enables selection by row.
	SelectRow

	// SelectColumn enables selection by column.
	SelectColumn
)

// List is a list of widgets. Rows or cells may optionally be selectable.
type List struct {
	grid                 *Grid
	itemHeight           int
	highlightColor       color.RGBA
	maxY                 int
	selectionMode        SelectionMode
	selectedX, selectedY int
	selectedFunc         func(index int) (accept bool)
	items                [][]Widget
	offset               int
	recreateGrid         bool
	scrollRect           image.Rectangle
	scrollWidth          int
	scrollAreaColor      color.RGBA
	scrollHandleColor    color.RGBA
	scrollDrag           bool
	drawBorder           bool
	sync.Mutex
}

const (
	initialPadding     = 5
	initialScrollWidth = 32
)

var (
	initialForeground   = color.RGBA{0, 0, 0, 255}
	initialBackground   = color.RGBA{255, 255, 255, 255}
	initialScrollArea   = color.RGBA{200, 200, 200, 255}
	initialScrollHandle = color.RGBA{108, 108, 108, 255}
)

// NewList returns a new List widget.
func NewList(itemHeight int, onSelected func(index int) (accept bool)) *List {
	return &List{
		grid:              NewGrid(),
		itemHeight:        itemHeight,
		highlightColor:    color.RGBA{255, 255, 255, 255},
		maxY:              -1,
		selectedY:         -1,
		selectedFunc:      onSelected,
		recreateGrid:      true,
		scrollWidth:       initialScrollWidth,
		scrollAreaColor:   initialScrollArea,
		scrollHandleColor: initialScrollHandle,
	}
}

// Rect returns the position and size of the widget.
func (l *List) Rect() image.Rectangle {
	l.Lock()
	defer l.Unlock()

	return l.grid.Rect()
}

// SetRect sets the position and size of the widget.
func (l *List) SetRect(r image.Rectangle) {
	l.Lock()
	defer l.Unlock()

	l.grid.SetRect(r)
	l.recreateGrid = true
}

// Background returns the background color of the widget.
func (l *List) Background() color.RGBA {
	l.Lock()
	defer l.Unlock()

	return l.grid.Background()
}

// SetBackground sets the background color of the widget.
func (l *List) SetBackground(background color.RGBA) {
	l.Lock()
	defer l.Unlock()

	l.grid.SetBackground(background)
}

// Focus returns the focus state of the widget.
func (l *List) Focus() bool {
	l.Lock()
	defer l.Unlock()

	return l.grid.Focus()
}

// SetFocus sets the focus state of the widget.
func (l *List) SetFocus(focus bool) (accept bool) {
	l.Lock()
	defer l.Unlock()

	return l.grid.SetFocus(focus)
}

// Visible returns the visibility of the widget.
func (l *List) Visible() bool {
	l.Lock()
	defer l.Unlock()

	return l.grid.Visible()
}

// SetVisible sets the visibility of the widget.
func (l *List) SetVisible(visible bool) {
	l.Lock()
	defer l.Unlock()

	l.grid.SetVisible(visible)
}

// SetColumnSizes sets the size of each column. A size of -1 represents an equal
// proportion of the available space.
func (l *List) SetColumnSizes(size ...int) {
	l.Lock()
	defer l.Unlock()

	l.grid.SetColumnSizes(size...)
}

// SetItemHeight sets the height of the list items.
func (l *List) SetItemHeight(itemHeight int) {
	l.Lock()
	defer l.Unlock()

	if l.itemHeight == itemHeight {
		return
	}
	l.itemHeight = itemHeight

	if l.maxY == -1 {
		return
	}
	rowSizes := make([]int, l.maxY+1)
	for i := range rowSizes {
		rowSizes[i] = l.itemHeight
	}
	l.grid.SetRowSizes(rowSizes...)
}

// SetSelectionMode sets the selection mode of the list.
func (l *List) SetSelectionMode(selectionMode SelectionMode) {
	l.Lock()
	defer l.Unlock()

	if l.selectionMode == selectionMode {
		return
	}
	l.selectionMode = selectionMode
}

// SetHighlightColor sets the color used to highlight the currently selected item.
func (l *List) SetHighlightColor(c color.RGBA) {
	l.Lock()
	defer l.Unlock()

	l.highlightColor = c
}

// SelectedItem returns the selected list item.
func (l *List) SelectedItem() (x int, y int) {
	l.Lock()
	defer l.Unlock()

	return l.selectedX, l.selectedY
}

// SetSelectedItem sets the selected list item.
func (l *List) SetSelectedItem(x int, y int) {
	l.Lock()
	defer l.Unlock()

	l.selectedX, l.selectedY = x, y
}

// SetScrollBarWidth sets the width of the scroll bar.
func (l *List) SetScrollBarWidth(width int) {
	l.Lock()
	defer l.Unlock()

	if l.scrollWidth == width {
		return
	}

	l.scrollWidth = width
}

// SetScrollBarColors sets the color of the scroll bar area and handle.
func (l *List) SetScrollBarColors(area color.RGBA, handle color.RGBA) {
	l.Lock()
	defer l.Unlock()

	l.scrollAreaColor, l.scrollHandleColor = area, handle
}

// SetSelectedFunc sets a handler which is called when a list item is selected.
// Providing a nil function value will remove the existing handler (if set).
// The handler may return false to return the selection to its original state.
func (l *List) SetSelectedFunc(f func(index int) (accept bool)) {
	l.Lock()
	defer l.Unlock()

	l.selectedFunc = f
}

// Children returns the children of the widget. Children are drawn in the
// order they are returned. Keyboard and mouse events are passed to children
// in reverse order.
func (l *List) Children() []Widget {
	l.Lock()
	defer l.Unlock()

	return l.grid.Children()
}

// AddChildAt adds a widget to the list at the specified position.
func (l *List) AddChildAt(w Widget, x int, y int) {
	l.Lock()
	defer l.Unlock()

	for i := y; i >= len(l.items); i-- {
		l.items = append(l.items, nil)
	}
	for i := x; i > len(l.items[y]); i-- {
		l.items[y] = append(l.items[y], nil)
	}
	l.items[y] = append(l.items[y], &ignoreMouse{w})
	if y > l.maxY {
		l.maxY = y
		l.recreateGrid = true
	}
}

func (l *List) showScrollBar() bool {
	return len(l.items) > l.grid.rect.Dy()/l.itemHeight
}

// clampOffset clamps the list offset.
func (l *List) clampOffset(offset int) int {
	if offset >= len(l.items)-(l.grid.rect.Dy()/l.itemHeight) {
		offset = len(l.items) - (l.grid.rect.Dy() / l.itemHeight)
	}
	if offset < 0 {
		offset = 0
	}
	return offset
}

// HandleKeyboard is called when a keyboard event occurs.
func (l *List) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	l.Lock()
	defer l.Unlock()

	return l.grid.HandleKeyboard(key, r)
}

// SetDrawBorder enables or disables borders being drawn around the list.
func (l *List) SetDrawBorder(drawBorder bool) {
	l.drawBorder = drawBorder
}

// HandleMouse is called when a mouse event occurs. Only mouse events that
// are on top of the widget are passed to the widget.
func (l *List) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	l.Lock()
	defer l.Unlock()

	_, scroll := ebiten.Wheel()
	if scroll != 0 {
		offset := l.clampOffset(l.offset - int(math.Round(scroll)))
		if offset != l.offset {
			l.offset = offset
			l.recreateGrid = true
		}
	}

	if l.showScrollBar() && (pressed || l.scrollDrag) {
		if pressed && cursor.In(l.scrollRect) {
			dragY := cursor.Y - l.grid.rect.Min.Y
			if dragY < 0 {
				dragY = 0
			} else if dragY > l.scrollRect.Dy() {
				dragY = l.scrollRect.Dy()
			}

			pct := float64(dragY) / float64(l.scrollRect.Dy())
			if pct < 0 {
				pct = 0
			} else if pct > 1 {
				pct = 1
			}

			lastOffset := l.offset
			offset := l.clampOffset(int(math.Round(float64(len(l.items)-(l.grid.rect.Dy()/l.itemHeight)) * pct)))
			if offset != lastOffset {
				l.offset = offset
				l.recreateGrid = true
			}
			l.scrollDrag = true
			return true, nil
		} else if !pressed {
			l.scrollDrag = false
		}
	}

	if !clicked || (cursor.X == 0 && cursor.Y == 0) {
		return true, nil
	}
	selected := (cursor.Y - l.grid.rect.Min.Y) / l.itemHeight
	if selected >= 0 && selected <= l.maxY {
		lastSelected := l.selectedY
		l.selectedY = selected
		if l.selectedFunc != nil {
			accept := l.selectedFunc(l.selectedY)
			if !accept {
				l.selectedY = lastSelected
			}
		}
	}
	return true, nil
}

// Draw draws the widget on the screen.
func (l *List) Draw(screen *ebiten.Image) error {
	l.Lock()
	defer l.Unlock()

	if l.recreateGrid {
		maxY := l.grid.rect.Dy() / l.itemHeight
		l.offset = l.clampOffset(l.offset)
		l.grid.Clear()
		rowSizes := make([]int, l.maxY+1)
		for i := range rowSizes {
			rowSizes[i] = l.itemHeight
		}
		l.grid.SetRowSizes(rowSizes...)
		var y int
		for i := range l.items {
			if i < l.offset {
				continue
			} else if y >= maxY {
				break
			}
			for x := range l.items[i] {
				w := l.items[i][x]
				if w == nil {
					continue
				}
				l.grid.AddChildAt(w, x, y, 1, 1)
			}
			y++
		}
		l.recreateGrid = false
	}

	// Draw grid.
	err := l.grid.Draw(screen)
	if err != nil {
		return err
	}

	// Highlight selection.
	drawHighlight := l.selectionMode != SelectNone && l.selectedY >= 0
	if drawHighlight {
		{
			x, y := l.grid.rect.Min.X, l.grid.rect.Min.Y+l.selectedY*l.itemHeight
			w, h := l.grid.rect.Dx(), l.itemHeight
			r := image.Rect(x, y, x+w, y+h)
			screen.SubImage(r).(*ebiten.Image).Fill(l.highlightColor)
		}
	}

	// Draw border.
	if l.drawBorder {
		const borderSize = 4
		screen.SubImage(image.Rect(l.grid.rect.Min.X, l.grid.rect.Min.Y, l.grid.rect.Max.X, l.grid.rect.Min.Y+borderSize)).(*ebiten.Image).Fill(Style.BorderColor)
		screen.SubImage(image.Rect(l.grid.rect.Min.X, l.grid.rect.Max.Y-borderSize, l.grid.rect.Max.X, l.grid.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
		screen.SubImage(image.Rect(l.grid.rect.Min.X, l.grid.rect.Min.Y, l.grid.rect.Min.X+borderSize, l.grid.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
		screen.SubImage(image.Rect(l.grid.rect.Max.X-borderSize, l.grid.rect.Min.Y, l.grid.rect.Max.X, l.grid.rect.Max.Y)).(*ebiten.Image).Fill(Style.BorderColor)
	}

	// Draw scroll bar.
	if !l.showScrollBar() {
		return nil
	}
	w, h := l.grid.rect.Dx(), l.grid.rect.Dy()
	scrollAreaX, scrollAreaY := l.grid.rect.Min.X+w-l.scrollWidth, l.grid.rect.Min.Y
	l.scrollRect = image.Rect(scrollAreaX, scrollAreaY, scrollAreaX+l.scrollWidth, scrollAreaY+h)

	scrollBarH := l.scrollWidth / 2
	if scrollBarH < 4 {
		scrollBarH = 4
	}

	scrollX, scrollY := l.grid.rect.Min.X+w-l.scrollWidth, l.grid.rect.Min.Y
	pct := float64(-l.offset) / float64(len(l.items)-(l.grid.rect.Dy()/l.itemHeight))
	scrollY -= int(float64(h-scrollBarH) * pct)
	scrollBarRect := image.Rect(scrollX, scrollY, scrollX+l.scrollWidth, scrollY+scrollBarH)

	screen.SubImage(l.scrollRect).(*ebiten.Image).Fill(l.scrollAreaColor)
	screen.SubImage(scrollBarRect).(*ebiten.Image).Fill(l.scrollHandleColor)
	return nil
}

// Clear clears all items from the list.
func (l *List) Clear() {
	l.Lock()
	defer l.Unlock()

	l.items = nil
	l.maxY = -1
	l.selectedX, l.selectedY = 0, -1
	l.offset = 0
	l.recreateGrid = true
}
