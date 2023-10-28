package etk

import "image"

// Frame is a widget container. All children are displayed at once. Children are
// not repositioned by default. Repositioning may be enabled via SetPositionChildren.
type Frame struct {
	*Box
	positionChildren bool
}

func NewFrame() *Frame {
	return &Frame{
		Box: NewBox(),
	}
}

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
