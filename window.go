package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Window is a widget paging mechanism. Only one widget added to a window is
// displayed at a time.
type Window struct {
	*Box

	allChildren []Widget

	active   int
	labels   []string
	hasLabel bool
}

// NewWindow returns a new Window widget.
func NewWindow() *Window {
	return &Window{
		Box: NewBox(),
	}
}

// SetRect sets the position and size of the widget.
func (w *Window) SetRect(r image.Rectangle) {
	w.Lock()
	defer w.Unlock()

	w.rect = r

	for _, wgt := range w.children {
		wgt.SetRect(r)
	}
}

// HandleKeyboard is called when a keyboard event occurs.
func (w *Window) HandleKeyboard(ebiten.Key, rune) (handled bool, err error) {
	return true, nil
}

// HandleMouse is called when a mouse event occurs.
func (w *Window) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return true, nil
}

// Draw draws the widget on the screen.
func (w *Window) Draw(screen *ebiten.Image) error {
	// TODO draw labels
	return nil
}

func (w *Window) childrenUpdated() {
	if len(w.allChildren) == 0 {
		w.children = nil
		return
	}
	w.children = []Widget{w.allChildren[w.active]}
}

// AddChild adds a child to the window.
func (w *Window) AddChild(wgt ...Widget) {
	w.allChildren = append(w.allChildren, wgt...)

	for _, widget := range wgt {
		widget.SetRect(w.rect)
	}

	blankLabels := make([]string, len(wgt))
	w.labels = append(w.labels, blankLabels...)

	w.childrenUpdated()
}

// AddChildWithLabel adds a child to the window with the specified label.
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
