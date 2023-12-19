package etk

import (
	"image"
	"image/color"
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
	sync.Mutex
}

// NewList returns a new List widget.
func NewList(itemHeight int, onSelected func(index int) (accept bool)) *List {
	return &List{
		grid:           NewGrid(),
		itemHeight:     itemHeight,
		highlightColor: color.RGBA{255, 255, 255, 255},
		maxY:           -1,
		selectedY:      -1,
		selectedFunc:   onSelected,
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

	l.grid.AddChildAt(&ignoreMouse{w}, x, y, 1, 1)
	if y > l.maxY {
		l.maxY = y
		rowSizes := make([]int, l.maxY+1)
		for i := range rowSizes {
			rowSizes[i] = l.itemHeight
		}
		l.grid.SetRowSizes(rowSizes...)
	}
}

// HandleKeyboard is called when a keyboard event occurs.
func (l *List) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	l.Lock()
	defer l.Unlock()

	return l.grid.HandleKeyboard(key, r)
}

// HandleMouse is called when a mouse event occurs. Only mouse events that
// are on top of the widget are passed to the widget.
func (l *List) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	l.Lock()
	defer l.Unlock()

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

	// Draw grid.
	err := l.grid.Draw(screen)
	if err != nil {
		return err
	}

	// Highlight selection.
	if l.selectionMode == SelectNone || l.selectedY < 0 {
		return nil
	}

	x, y := l.grid.rect.Min.X, l.grid.rect.Min.Y+l.selectedY*l.itemHeight
	w, h := l.grid.rect.Dx(), l.itemHeight
	r := image.Rect(x, y, x+w, y+h)
	screen.SubImage(r).(*ebiten.Image).Fill(l.highlightColor)
	return nil
}

// Clear clears all items from the list.
func (l *List) Clear() {
	l.Lock()
	defer l.Unlock()

	l.grid.Clear()
	l.maxY = -1
	l.selectedX, l.selectedY = 0, -1
}
