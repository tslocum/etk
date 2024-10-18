package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/sfnt"
)

// Window displays a single child widget at a time, and includes a list to
// view other child widgets. Window.Show must be called after adding a widget.
type Window struct {
	*Box
	font      *sfnt.Font
	fontSize  int
	frameSize int
	titleSize int
	titles    []string
	active    int
	modified  bool
}

// NewWindow returns a new Window widget.
func NewWindow() *Window {
	return &Window{
		Box:       NewBox(),
		font:      Style.TextFont,
		fontSize:  Scale(Style.TextSize),
		frameSize: Scale(4),
		titleSize: Scale(40),
		active:    -1,
	}
}

// SetRect sets the position and size of the widget.
func (w *Window) SetRect(r image.Rectangle) {
	w.Lock()
	defer w.Unlock()

	w.rect = r
	w.modified = true
}

// SetFont sets the font and text size of the window titles. Scaling is not applied.
func (w *Window) SetFont(fnt *sfnt.Font, size int) {
	w.Lock()
	defer w.Unlock()

	w.font = fnt
	w.fontSize = size
}

// SetFrameSize sets the size of the frame around each window.
func (w *Window) SetFrameSize(size int) {
	w.Lock()
	defer w.Unlock()

	w.frameSize = size
	w.modified = true
}

// SetTitleSize sets the height of the title bars.
func (w *Window) SetTitleSize(size int) {
	w.Lock()
	defer w.Unlock()

	w.titleSize = size
	w.modified = true
}

// Show displays the specified child widget within the Window.
func (w *Window) Show(index int) {
	w.Lock()
	defer w.Unlock()

	if index >= 0 && index < len(w.children) {
		w.active = index
	} else {
		w.active = -1
	}
	w.modified = true
}

// Hide hides the currently visible child widget.
func (w *Window) Hide() {
	w.Lock()
	defer w.Unlock()

	w.active = -1
	w.modified = true
}

// Children returns the children of the widget.
func (w *Window) Children() []Widget {
	w.Lock()
	defer w.Unlock()

	if w.active >= 0 && w.active < len(w.children) {
		return []Widget{w.children[w.active]}
	}
	return nil
}

// Clear removes all children from the widget.
func (w *Window) Clear() {
	w.Lock()
	defer w.Unlock()

	w.children = w.children[:0]
	w.titles = w.titles[:0]
	w.active = -1
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
	if w.modified {
		if w.active >= 0 {
			w.children[w.active].SetRect(w.rect)
		}
		w.modified = false
	}
	return nil
}

// AddChild adds a child to the window.
func (w *Window) AddChild(wgt ...Widget) {
	w.Lock()
	defer w.Unlock()

	for _, widget := range wgt {
		w.children = append(w.children, widget)
		w.titles = append(w.titles, "")
	}
	w.modified = true
}

// AddChildWithTitle adds a child to the window with the specified window title.
func (w *Window) AddChildWithTitle(wgt Widget, title string) int {
	w.Lock()
	defer w.Unlock()

	w.children = append(w.children, wgt)
	w.titles = append(w.titles, title)

	w.modified = true
	return len(w.children) - 1
}
