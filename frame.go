package etk

import "image"

// Frame is a widget container. All children are displayed at once. Children are
// not repositioned by default. Repositioning may be enabled via SetPositionChildren.
type Frame struct {
	*Box
	padding          int
	positionChildren bool
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
	f.reposition()
}

// SetRect sets the position and size of the widget.
func (f *Frame) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	f.rect = r
	f.reposition()
}

// SetPositionChildren sets a flag that determines whether child widgets are
// repositioned when the Frame is repositioned.
func (f *Frame) SetPositionChildren(position bool) {
	f.Lock()
	defer f.Unlock()

	f.positionChildren = position
	f.reposition()
}

// AddChild adds a child to the widget.
func (f *Frame) AddChild(w ...Widget) {
	f.Lock()
	defer f.Unlock()

	f.children = append(f.children, w...)

	if f.positionChildren {
		r := f.rect.Inset(f.padding)
		for _, wgt := range w {
			wgt.SetRect(r)
		}
	}
}

func (f *Frame) reposition() {
	if !f.positionChildren {
		return
	}
	r := f.rect.Inset(f.padding)
	for _, w := range f.children {
		w.SetRect(r)
	}
}
