package etk

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/sfnt"
)

// Window displays a single child widget at a time, and includes a list to
// view other child widgets. Window.Show must be called after adding a widget.
type Window struct {
	*Box
	font         *sfnt.Font
	fontSize     int
	frameSize    int
	list         *List
	listH        Alignment
	listV        Alignment
	listSize     int
	listWidget   *WithoutFocus
	defaultFocus []Widget
	labels       []string
	active       int
	modified     bool
	listModified bool
	firstDraw    bool
}

// NewWindow returns a new Window widget.
func NewWindow() *Window {
	w := &Window{
		Box:       NewBox(),
		font:      Style.TextFont,
		fontSize:  Scale(Style.TextSize),
		frameSize: Scale(4),
		listSize:  Scale(128),
		listH:     AlignEnd,
		listV:     AlignCenter,
		active:    -1,
		firstDraw: true,
	}
	w.list = NewList(int(float64(Scale(Style.TextSize))*1.5), w.selectItem)
	w.listWidget = &WithoutFocus{w.list}
	return w
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
	w.list.SetItemHeight(size)
}

// SetFrameSize sets the size of the frame around each window.
func (w *Window) SetFrameSize(size int) {
	w.Lock()
	defer w.Unlock()

	w.frameSize = size
	w.modified = true
}

// SetListSize sets the width or height of the window tab list.
func (w *Window) SetListSize(size int) {
	w.Lock()
	defer w.Unlock()

	w.listSize = size
	w.modified = true
}

// SetListHorizontal sets the horizontal alignment of the window tab list.
func (w *Window) SetListHorizontal(h Alignment) {
	w.Lock()
	defer w.Unlock()

	w.listH = h
	w.modified = true
}

// SetListVertical sets the vertical alignment of the window tab list.
func (w *Window) SetListVertical(v Alignment) {
	w.Lock()
	defer w.Unlock()

	w.listV = v
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
		if w.listSize > 0 {
			return []Widget{w.children[w.active], w.listWidget}
		}
		return []Widget{w.children[w.active]}
	} else if w.listSize > 0 {
		return []Widget{w.listWidget}
	}
	return nil
}

// Clear removes all children from the widget.
func (w *Window) Clear() {
	w.Lock()
	defer w.Unlock()

	w.children = w.children[:0]
	w.labels = w.labels[:0]
	w.active = -1
	w.listModified = true
	w.firstDraw = true
}

func (w *Window) selectItem(index int) (accept bool) {
	if index >= 0 && index < len(w.children) {
		w.active = index
		w.modified = true
		SetFocus(w.defaultFocus[index])
	}
	return true
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
	if w.listModified {
		w.list.SetItemHeight(int(float64(Scale(w.fontSize)) * 1.5))
		w.list.Clear()
		for i := range w.labels {
			label := w.labels[i]
			if label == "" {
				label = fmt.Sprintf("#%d", i)
			}
			t := NewText(label)
			t.SetFont(w.font, Scale(w.fontSize))
			t.SetAutoResize(true)
			w.list.AddChildAt(t, 0, i)
			w.list.SetSelectedItem(0, w.active)
		}
		w.listModified = false
	}
	if w.modified {
		if w.active >= 0 {
			wr := w.rect
			if w.listSize > 0 {
				var lr image.Rectangle
				switch w.listH {
				case AlignStart:
					lr = image.Rect(wr.Min.X, wr.Min.Y, wr.Min.X+w.listSize, wr.Max.Y)
				case AlignEnd:
					lr = image.Rect(wr.Max.X-w.listSize, wr.Min.Y, wr.Max.X, wr.Max.Y)
				}
				switch w.listV {
				case AlignStart:
					lr = image.Rect(wr.Min.X, wr.Min.Y, wr.Max.X, wr.Min.Y+w.listSize)
				case AlignEnd:
					lr = image.Rect(wr.Min.X, wr.Max.Y-w.listSize, wr.Max.X, wr.Max.Y)
				}
				dx, dy := lr.Dx(), lr.Dy()
				if dx > 0 && dy > 0 && dx <= wr.Dx() && dy <= wr.Dy() {
					switch w.listH {
					case AlignStart:
						wr.Min.X += w.listSize
					case AlignEnd:
						wr.Max.X -= w.listSize
					}
					switch w.listV {
					case AlignStart:
						wr.Min.Y += w.listSize
					case AlignEnd:
						wr.Max.Y -= w.listSize
					}
					w.list.SetRect(lr)
				}
			}
			w.children[w.active].SetRect(wr)
			if w.firstDraw {
				SetFocus(w.children[w.active])
				w.firstDraw = false
			}
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
		w.defaultFocus = append(w.defaultFocus, widget)
		w.labels = append(w.labels, "")
	}
	w.modified = true
	w.listModified = true
}

// AddChildWithLabel adds a child to the window with the specified default focus and list entry label.
func (w *Window) AddChildWithLabel(wgt Widget, defaultFocus Widget, label string) int {
	w.Lock()
	defer w.Unlock()

	w.children = append(w.children, wgt)
	w.defaultFocus = append(w.defaultFocus, defaultFocus)
	w.labels = append(w.labels, label)

	w.modified = true
	w.listModified = true
	return len(w.children) - 1
}
