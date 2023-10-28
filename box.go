package etk

import (
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	rect       image.Rectangle
	children   []Widget
	background color.RGBA
	visible    bool

	sync.Mutex
}

func NewBox() *Box {
	return &Box{
		visible: true,
	}
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

func (b *Box) Background() color.RGBA {
	b.Lock()
	defer b.Unlock()

	return b.background
}

func (b *Box) SetBackground(background color.RGBA) {
	b.Lock()
	defer b.Unlock()

	b.background = background
}

func (b *Box) SetFocus(focus bool) bool {
	return false
}

func (b *Box) SetVisible(visible bool) {
	b.visible = visible
}
func (b *Box) Visible() bool {
	return b.visible
}

func (b *Box) Focus() bool {
	return false
}

func (b *Box) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return true, nil
}

func (b *Box) HandleKeyboard(key ebiten.Key, r rune) (handled bool, err error) {
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

func (b *Box) Draw(screen *ebiten.Image) error {
	return nil
}

var _ Widget = &Box{}
