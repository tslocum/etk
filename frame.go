package etk

import "image"

// Frame is a widget container. All children are displayed at once. Children are
// not repositioned by default. Repositioning may be enabled via SetPositionChildren.
type Frame struct {
	*Box
	positionChildren bool
}

// NewFrame returns a new Frame widget.
func NewFrame() *Frame {
	return &Frame{
		Box: NewBox(),
	}
}

// SetRect sets the position and size of the widget.
func (f *Frame) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	f.rect = r

	if f.positionChildren {
		for _, w := range f.children {
			w.SetRect(f.rect)
		}
	}
}

// SetPositionChildren sets a flag that determines whether child widgets are
// repositioned when the Frame is repositioned.
func (f *Frame) SetPositionChildren(position bool) {
	f.Lock()
	defer f.Unlock()

	f.positionChildren = position

	if f.positionChildren {
		for _, w := range f.children {
			w.SetRect(f.rect)
		}
	}
}

// AddChild adds a child to the widget.
func (f *Frame) AddChild(w ...Widget) {
	f.Lock()
	defer f.Unlock()

	f.children = append(f.children, w...)

	if f.positionChildren {
		for _, wgt := range w {
			wgt.SetRect(f.rect)
		}
	}
}
