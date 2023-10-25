package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Window displays and passes input to only one child widget at a time.
type Window struct {
	*Box

	allChildren []Widget

	active   int
	labels   []string
	hasLabel bool
}

func NewWindow() *Window {
	return &Window{
		Box: NewBox(),
	}
}

func (w *Window) childrenUpdated() {
	if len(w.allChildren) == 0 {
		w.children = nil
		return
	}
	w.children = []Widget{w.allChildren[w.active]}
}

func (w *Window) SetRect(r image.Rectangle) {
	w.Lock()
	defer w.Unlock()

	w.rect = r

	for _, wgt := range w.children {
		wgt.SetRect(r)
	}
}

func (w *Window) AddChild(wgt ...Widget) {
	w.allChildren = append(w.allChildren, wgt...)

	for _, widget := range wgt {
		widget.SetRect(w.rect)
	}

	blankLabels := make([]string, len(wgt))
	w.labels = append(w.labels, blankLabels...)

	w.childrenUpdated()
}

func (w *Window) AddChildWithLabel(wgt Widget, label string) {
	w.Lock()
	defer w.Unlock()

	wgt.SetRect(w.rect)

	w.allChildren = append(w.allChildren, wgt)
	w.labels = append(w.labels, label)

	if label != "" {
		w.hasLabel = true
	}

	w.childrenUpdated()
}

func (w *Window) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return true, nil
}

func (w *Window) HandleKeyboard() (handled bool, err error) {
	return true, nil
}

func (w *Window) Draw(screen *ebiten.Image) error {
	// TODO draw labels
	return nil
}
