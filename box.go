package etk

import (
	"image"
	"sync"
)

type Box struct {
	rect image.Rectangle

	children []Widget

	focus bool

	sync.Mutex
}

func NewBox() *Box {
	return &Box{}
}

func (b *Box) Rect() image.Rectangle {
	b.Lock()
	defer b.Unlock()

	return b.rect
}

func (b *Box) SetRect(r image.Rectangle) {
	b.Lock()
	defer b.Unlock()

	b.rect = r
}

func (b *Box) SetFocus(focus bool) {
	b.focus = focus
}

func (b *Box) Focus() bool {
	return b.focus
}

func (b *Box) Children() []Widget {
	b.Lock()
	defer b.Unlock()

	return b.children
}

func (b *Box) AddChild(w ...Widget) {
	b.Lock()
	defer b.Unlock()

	b.children = append(b.children, w...)
}
