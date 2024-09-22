package etk

import "image"

// Frame is a widget container. All children are displayed at once. Children are
// not repositioned by default. Repositioning may be enabled via SetPositionChildren.
type Frame struct {
	*Box
	padding          int
	positionChildren bool
	maxWidth         int
	maxHeight        int
}

// NewFrame returns a new Frame widget.
func NewFrame(w ...Widget) *Frame {
	f := &Frame{
		Box: NewBox(),
	}
	f.AddChild(w...)
	return f
}

// SetPadding sets the amount of padding around widgets in the frame.
func (f *Frame) SetPadding(padding int) {
	f.Lock()
	defer f.Unlock()

	f.padding = padding
	f.repositionAll()
}

// SetRect sets the position and size of the widget.
func (f *Frame) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	f.rect = r
	f.repositionAll()
}

// SetPositionChildren sets a flag that determines whether child widgets are
// repositioned when the Frame is repositioned.
func (f *Frame) SetPositionChildren(position bool) {
	f.Lock()
	defer f.Unlock()

	f.positionChildren = position
	f.repositionAll()
}

// SetMaxWidth sets the maximum width of widgets within the frame. This will
// only have an effect after SetPositionChildren(true) is called. 0 to disable.
func (f *Frame) SetMaxWidth(w int) {
	f.Lock()
	defer f.Unlock()

	f.maxWidth = w
	f.repositionAll()
}

// SetMaxHeight sets the maximum height of widgets within the frame. This will
// only have an effect after SetPositionChildren(true) is called. 0 to disable.
func (f *Frame) SetMaxHeight(h int) {
	f.Lock()
	defer f.Unlock()

	f.maxHeight = h
	f.repositionAll()
}

// AddChild adds a child to the widget.
func (f *Frame) AddChild(w ...Widget) {
	f.Lock()
	defer f.Unlock()

	f.children = append(f.children, w...)

	if !f.positionChildren {
		return
	}
	for _, child := range w {
		f.repositionChild(child)
	}
}

func (f *Frame) repositionChild(w Widget) {
	r := f.rect.Inset(f.padding)
	if f.maxWidth > 0 && r.Dx() > f.maxWidth {
		r.Max.X = r.Min.X + f.maxWidth
	}
	if f.maxHeight > 0 && r.Dy() > f.maxHeight {
		r.Max.Y = r.Min.Y + f.maxHeight
	}
	w.SetRect(r)
}

func (f *Frame) repositionAll() {
	if !f.positionChildren {
		return
	}
	for _, w := range f.children {
		f.repositionChild(w)
	}
}
