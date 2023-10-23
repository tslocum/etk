package etk

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	rect image.Rectangle

	children []Widget

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

func (b *Box) SetFocus(focus bool) bool {
	return false
}

func (b *Box) Focus() bool {
	return false
}

func (b *Box) HandleKeyboardEvent(key ebiten.Key, r rune) (handled bool, err error) {
	return false, nil
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
