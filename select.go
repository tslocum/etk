package etk

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Select is a dropdown selection widget.
type Select struct {
	*Box
	label    *Text
	list     *List
	onSelect func(index int) (accept bool)
	items    []string
	open     bool
}

// NewSelect returns a new Select widget.
func NewSelect(itemHeight int, onSelect func(index int) (accept bool)) *Select {
	s := &Select{
		Box:      NewBox(),
		label:    NewText(""),
		onSelect: onSelect,
	}
	s.label.SetVertical(AlignCenter)
	s.label.SetForeground(Style.ButtonTextColor)
	s.SetBackground(Style.ButtonBgColor)
	s.list = NewList(itemHeight, s.selectList)
	s.list.SetBackground(Style.ButtonBgColor)
	s.list.SetDrawBorder(true)
	s.list.SetVisible(false)
	s.list.SetSelectionMode(SelectRow)
	s.AddChild(s.list)
	s.updateLabel()
	return s
}

// SetRect sets the position and size of the widget.
func (s *Select) SetRect(r image.Rectangle) {
	s.Lock()
	defer s.Unlock()
	s.rect = r
	s.label.SetRect(r)
	listRect := r.Add(image.Point{X: 0, Y: r.Dy()})
	itemCount := len(s.items)
	listRect.Max.Y = listRect.Min.Y + itemCount*s.list.itemHeight
	_, height := ScreenSize()
	if listRect.Max.Y > height {
		listRect.Max.Y = height
	}
	s.list.SetRect(listRect)
}

// SetHighlightColor sets the color used to highlight the currently selected item.
func (s *Select) SetHighlightColor(c color.RGBA) {
	s.list.SetHighlightColor(c)
}

// SetSelectedItem sets the currently selected item.
func (s *Select) SetSelectedItem(index int) {
	s.Lock()
	defer s.Unlock()
	if index < 0 || index >= len(s.items) {
		return
	}
	s.list.SetSelectedItem(0, index)
	s.updateLabel()
}

// Children returns the children of the widget.
func (s *Select) Children() []Widget {
	s.Lock()
	defer s.Unlock()

	return s.children
}

// AddChild adds a child to the widget. Selection options are added via AddOption.
func (s *Select) AddChild(w ...Widget) {
	s.Lock()
	defer s.Unlock()

	s.children = append(s.children, w...)
}

// Clear removes all children from the widget.
func (s *Select) Clear() {
	s.Lock()
	defer s.Unlock()

	s.items = nil
	s.list.Clear()
	s.updateLabel()
}

// AddOption adds an option to the widget.
func (s *Select) AddOption(label string) {
	s.Lock()
	defer s.Unlock()

	s.items = append(s.items, label)
	if len(s.items) == 1 {
		s.list.selectedY = 0
		s.updateLabel()
	}

	t := NewText(label)
	t.SetVertical(AlignCenter)
	t.SetForeground(Style.ButtonTextColor)
	s.list.AddChildAt(t, 0, len(s.items)-1)
}

func (s *Select) updateLabel() {
	var text string
	if len(s.items) > 0 && s.list.selectedY >= 0 && s.list.selectedY < len(s.items) {
		text = s.items[s.list.selectedY]
	}
	if s.open {
		text = "▼ " + text
	} else {
		text = "▶ " + text
	}
	s.label.SetText(text)
}

func (s *Select) selectList(index int) (accept bool) {
	s.Lock()
	s.list.grid.visible = false
	s.open = false
	onSelect := s.onSelect
	s.Unlock()

	if onSelect != nil {
		if !onSelect(index) {
			return false
		}
	}

	s.list.selectedY = index
	s.updateLabel()
	return true
}

// SetMenuVisible sets the visibility of the dropdown menu.
func (s *Select) SetMenuVisible(visible bool) {
	s.Lock()
	defer s.Unlock()

	s._setMenuVisible(visible)
}

func (s *Select) _setMenuVisible(visible bool) {
	s.open = visible
	s.list.SetVisible(visible)
	s.updateLabel()
}

// HandleKeyboard is called when a keyboard event occurs.
func (s *Select) HandleKeyboard(ebiten.Key, rune) (handled bool, err error) {
	return false, nil
}

// HandleMouse is called when a mouse event occurs.
func (s *Select) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	s.Lock()
	defer s.Unlock()

	if clicked {
		s._setMenuVisible(!s.open)
	}
	return true, nil
}

// Draw draws the widget on the screen.
func (s *Select) Draw(screen *ebiten.Image) error {
	s.Lock()
	defer s.Unlock()

	// Draw label.
	s.label.Draw(screen)

	// Draw border.
	r := s.rect
	borderSize := Scale(Style.BorderSize)
	screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+borderSize, r.Max.Y)).(*ebiten.Image).Fill(Style.BorderColorLeft)
	screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+borderSize)).(*ebiten.Image).Fill(Style.BorderColorTop)
	screen.SubImage(image.Rect(r.Max.X-borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(Style.BorderColorRight)
	screen.SubImage(image.Rect(r.Min.X, r.Max.Y-borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(Style.BorderColorBottom)

	return nil
}
