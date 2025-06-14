package etk

import (
	"image"
	"image/color"
	"math"
	"sync"
	"time"

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

// List is a list of widgets.
type List struct {
	rect                 image.Rectangle
	grid                 *Grid
	focused              bool
	itemHeight           int
	highlightColor       color.RGBA
	maxY                 int
	selectionMode        SelectionMode
	selectedX, selectedY int
	selectedTime         time.Time
	onChange             func(index int) (accept bool)
	onConfirm            func(index int)
	items                [][]Widget
	offset               int
	recreateGrid         bool
	scrollRect           image.Rectangle
	scrollWidth          int
	scrollAreaColor      color.RGBA
	scrollHandleColor    color.RGBA
	scrollBorderSize     int
	scrollBorderTop      color.RGBA
	scrollBorderRight    color.RGBA
	scrollBorderBottom   color.RGBA
	scrollBorderLeft     color.RGBA
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
func NewList(itemHeight int, onChange func(index int) (accept bool), onConfirm func(index int)) *List {
	return &List{
		grid:               NewGrid(),
		itemHeight:         itemHeight,
		highlightColor:     color.RGBA{128, 128, 128, 255},
		maxY:               -1,
		selectionMode:      SelectRow,
		selectedX:          -1,
		selectedY:          -1,
		onChange:           onChange,
		onConfirm:          onConfirm,
		recreateGrid:       true,
		scrollWidth:        initialScrollWidth,
		scrollAreaColor:    initialScrollArea,
		scrollHandleColor:  initialScrollHandle,
		scrollBorderSize:   Style.ScrollBorderSize,
		scrollBorderTop:    Style.ScrollBorderTop,
		scrollBorderRight:  Style.ScrollBorderRight,
		scrollBorderBottom: Style.ScrollBorderBottom,
		scrollBorderLeft:   Style.ScrollBorderLeft,
	}
}

// Rect returns the position and size of the widget.
func (l *List) Rect() image.Rectangle {
	l.Lock()
	defer l.Unlock()

	return l.rect
}

// SetRect sets the position and size of the widget.
func (l *List) SetRect(r image.Rectangle) {
	l.Lock()
	defer l.Unlock()

	l.rect = r
	if l.showScrollBar() {
		r.Max.X -= l.scrollWidth
	}
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

	return l.focused
}

// SetFocus sets the focus state of the widget.
func (l *List) SetFocus(focus bool) (accept bool) {
	l.Lock()
	defer l.Unlock()

	l.focused = focus
	return true
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

// Clip returns whether the widget and its children are restricted to drawing
// within the widget's rect area of the screen. For best performance, Clip
// should return false unless clipping is actually needed.
func (l *List) Clip() bool {
	return true
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

// SetScrollBorderSize sets the size of the border around the scroll bar handle.
func (l *List) SetScrollBorderSize(size int) {
	l.Lock()
	defer l.Unlock()

	l.scrollBorderSize = size
}

// SetScrollBorderColor sets the color of the top, right, bottom and left border
// of the scroll bar handle.
func (l *List) SetScrollBorderColors(top color.RGBA, right color.RGBA, bottom color.RGBA, left color.RGBA) {
	l.Lock()
	defer l.Unlock()

	l.scrollBorderTop = top
	l.scrollBorderRight = right
	l.scrollBorderBottom = bottom
	l.scrollBorderLeft = left
}

// SetChangeFunc sets a handler which is called when the selected item changes.
// Providing a nil function value will remove the existing handler (if set).
// The handler may return false to return the selection to its original state.
func (l *List) SetChangeFunc(onChange func(index int) (accept bool)) {
	l.Lock()
	defer l.Unlock()

	l.onChange = onChange
}

// SetConfirmFunc sets a handler which is called when the list selection is confirmed.
// Providing a nil function value will remove the existing handler (if set).
func (l *List) SetConfirmFunc(onConfirm func(index int)) {
	l.Lock()
	defer l.Unlock()

	l.onConfirm = onConfirm
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
	if l.selectionMode == SelectNone {
		w = &WithoutMouseExceptScroll{Widget: w}
	} else {
		w = &WithoutMouse{Widget: w}
	}
	l.items[y] = append(l.items[y], w)
	if y > l.maxY {
		l.maxY = y
		l.recreateGrid = true
	}
}

// Rows returns the number of rows in the list.
func (l *List) Rows() int {
	l.Lock()
	defer l.Unlock()

	return l.maxY + 1
}

func (l *List) showScrollBar() bool {
	return len(l.items) > l.rect.Dy()/l.itemHeight
}

// clampOffset clamps the list offset.
func (l *List) clampOffset(offset int) int {
	if offset >= len(l.items)*l.itemHeight-l.rect.Dy() {
		offset = len(l.items)*l.itemHeight - l.rect.Dy()
	}
	if offset < 0 {
		offset = 0
	}
	return offset
}

// Cursor returns the cursor shape shown when a mouse cursor hovers over the
// widget, or -1 to let widgets beneath determine the cursor shape.
func (l *List) Cursor() ebiten.CursorShapeType {
	return ebiten.CursorShapeDefault
}

// HandleKeyboard is called when a keyboard event occurs.
func (l *List) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
	l.Lock()
	defer l.Unlock()

	if r == 0 {
		// Handle confirmation.
		for _, confirmKey := range Bindings.ConfirmKeyboard {
			if key == confirmKey {
				onConfirm := l.onConfirm
				if onConfirm != nil {
					l.Unlock()
					onConfirm(l.selectedY)
					l.Lock()
				}
				return true, nil
			}
		}

		// Handle movement.
		move := func(x int, y int) {
			y = l.selectedY + y
			if y >= 0 && y <= l.maxY {
				l.selectedY = y
			}
		}
		for _, leftKey := range Bindings.MoveLeftKeyboard {
			if key == leftKey {
				move(-1, 0)
				return true, nil
			}
		}
		for _, rightKey := range Bindings.MoveRightKeyboard {
			if key == rightKey {
				move(1, 0)
				return true, nil
			}
		}
		for _, downKey := range Bindings.MoveDownKeyboard {
			if key == downKey {
				move(0, 1)
				return true, nil
			}
		}
		for _, upKey := range Bindings.MoveUpKeyboard {
			if key == upKey {
				move(0, -1)
				return true, nil
			}
		}
	}

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
		if scroll < -maxScroll {
			scroll = -maxScroll
		} else if scroll > maxScroll {
			scroll = maxScroll
		}
		offset := l.clampOffset(l.offset - int(math.Round(scroll))*3*l.itemHeight)
		if offset != l.offset {
			l.offset = offset
			l.recreateGrid = true
		}
	}

	if l.showScrollBar() && (pressed || l.scrollDrag) {
		if pressed && cursor.In(l.scrollRect) {
			dragY := cursor.Y - l.rect.Min.Y
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
			offset := l.clampOffset(int(math.Round(float64(len(l.items)*l.itemHeight-l.rect.Dy()) * pct)))
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
	selected := (l.offset + cursor.Y - l.rect.Min.Y) / l.itemHeight
	if selected >= 0 && selected <= l.maxY {
		onChange := l.onChange
		if onChange != nil {
			l.Unlock()
			accept := onChange(selected)
			l.Lock()
			if !accept {
				return true, nil
			}
		}
		lastSelected := l.selectedY
		l.selectedY = selected

		if selected == lastSelected && time.Since(l.selectedTime) <= Bindings.DoubleClickThreshold {
			onConfirm := l.onConfirm
			if onConfirm != nil {
				l.Unlock()
				onConfirm(l.selectedY)
				l.Lock()
			}
			l.selectedTime = time.Time{}
			return true, nil
		}

		l.selectedTime = time.Now()
	}
	return true, nil
}

func (l *List) _recreateCrid(screen *ebiten.Image) {
	maxY := l.rect.Dy()/l.itemHeight + 1
	if maxY < 2 {
		maxY = 2
	}
	l.offset = l.clampOffset(l.offset)

	l.grid.Clear()
	rowSizes := make([]int, l.maxY+1)
	for i := range rowSizes {
		rowSizes[i] = l.itemHeight
	}
	l.grid.SetRowSizes(rowSizes...)
	var y int
	for i := range l.items {
		if i*l.itemHeight < l.offset-l.itemHeight+1 {
			continue
		} else if y > maxY {
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

	r := l.rect
	if l.showScrollBar() {
		r.Max.X -= l.scrollWidth
	}
	remainder := l.offset % l.itemHeight
	r.Min.Y = l.rect.Min.Y - remainder
	l.grid.SetRect(r)

	l.recreateGrid = false
}

// Draw draws the widget on the screen.
func (l *List) Draw(screen *ebiten.Image) error {
	l.Lock()
	defer l.Unlock()

	if l.recreateGrid {
		l._recreateCrid(screen)
	}

	// Draw grid.
	err := l.grid.Draw(screen)
	if err != nil {
		return err
	}

	// Highlight selection.
	drawHighlight := l.selectionMode != SelectNone && l.selectedY >= 0
	if drawHighlight {
		x, y := l.rect.Min.X, l.rect.Min.Y+l.selectedY*l.itemHeight-l.offset
		w, h := l.rect.Dx(), l.itemHeight
		r := clampRect(image.Rect(x, y, x+w, y+h), l.rect)
		if r.Dx() > 0 && r.Dy() > 0 {
			screen.SubImage(r).(*ebiten.Image).Fill(l.highlightColor)
		}
	}

	// Draw border.
	if l.drawBorder {
		const borderSize = 4
		screen.SubImage(image.Rect(l.grid.rect.Min.X, l.grid.rect.Min.Y, l.grid.rect.Max.X, l.grid.rect.Min.Y+borderSize)).(*ebiten.Image).Fill(Style.ButtonBorderBottom)
		screen.SubImage(image.Rect(l.grid.rect.Min.X, l.grid.rect.Max.Y-borderSize, l.grid.rect.Max.X, l.grid.rect.Max.Y)).(*ebiten.Image).Fill(Style.ButtonBorderBottom)
		screen.SubImage(image.Rect(l.grid.rect.Min.X, l.grid.rect.Min.Y, l.grid.rect.Min.X+borderSize, l.grid.rect.Max.Y)).(*ebiten.Image).Fill(Style.ButtonBorderBottom)
		screen.SubImage(image.Rect(l.grid.rect.Max.X-borderSize, l.grid.rect.Min.Y, l.grid.rect.Max.X, l.grid.rect.Max.Y)).(*ebiten.Image).Fill(Style.ButtonBorderBottom)
	}

	// Draw scroll bar.
	if !l.showScrollBar() {
		return nil
	}
	w, h := l.rect.Dx(), l.rect.Dy()
	scrollAreaX, scrollAreaY := l.rect.Min.X+w-l.scrollWidth, l.rect.Min.Y
	l.scrollRect = image.Rect(scrollAreaX, scrollAreaY, scrollAreaX+l.scrollWidth, scrollAreaY+h)

	scrollBarH := l.scrollWidth / 2
	if scrollBarH < 4 {
		scrollBarH = 4
	}

	scrollX, scrollY := l.rect.Min.X+w-l.scrollWidth, l.rect.Min.Y
	pct := float64(-l.offset) / float64(len(l.items)*l.itemHeight-l.rect.Dy())
	scrollY -= int(float64(h-scrollBarH) * pct)
	scrollBarRect := image.Rect(scrollX, scrollY, scrollX+l.scrollWidth, scrollY+scrollBarH)

	screen.SubImage(l.scrollRect).(*ebiten.Image).Fill(l.scrollAreaColor)
	screen.SubImage(scrollBarRect).(*ebiten.Image).Fill(l.scrollHandleColor)

	// Draw scroll handle border.
	if l.scrollBorderSize != 0 {
		r := scrollBarRect
		screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+l.scrollBorderSize, r.Max.Y)).(*ebiten.Image).Fill(l.scrollBorderLeft)
		screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+l.scrollBorderSize)).(*ebiten.Image).Fill(l.scrollBorderTop)
		screen.SubImage(image.Rect(r.Max.X-l.scrollBorderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(l.scrollBorderRight)
		screen.SubImage(image.Rect(r.Min.X, r.Max.Y-l.scrollBorderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(l.scrollBorderBottom)
	}
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
