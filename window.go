package etk

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/sfnt"
)

// Window displays child widgets in floating or maximized windows.
type Window struct {
	*Box
	font       *sfnt.Font
	fontSize   int
	frameSize  int
	titleSize  int
	titles     []*Text
	floating   []bool
	fullscreen []Widget
	modified   bool
}

// NewWindow returns a new Window widget.
func NewWindow() *Window {
	return &Window{
		Box:       NewBox(),
		font:      Style.TextFont,
		fontSize:  Scale(Style.TextSize),
		frameSize: Scale(4),
		titleSize: Scale(40),
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

	for _, title := range w.titles {
		title.SetFont(w.font, w.fontSize)
	}
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

// SetFullscreen expands the specified widget to fill the netire screen, hiding
// the title bar. When -1 is provided, the currently fullscreen widget is
// restored to its a normal size.
func (w *Window) SetFullscreen(index int) {
	w.Lock()
	defer w.Unlock()

	if index == -1 {
		w.fullscreen = w.fullscreen[:0]
	} else if index >= 0 && index < len(w.children) {
		w.fullscreen = append(w.fullscreen[:0], w.children[index])
	}
	w.modified = true
}

// Children returns the children of the widget.
func (w *Window) Children() []Widget {
	w.Lock()
	defer w.Unlock()

	if len(w.fullscreen) != 0 {
		return w.fullscreen
	}
	return w.children
}

// Clear removes all children from the widget.
func (w *Window) Clear() {
	w.Lock()
	defer w.Unlock()

	w.children = w.children[:0]
	w.titles = w.titles[:0]
	w.floating = w.floating[:0]
	w.fullscreen = w.fullscreen[:0]
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
		if len(w.fullscreen) != 0 {
			w.fullscreen[0].SetRect(w.rect)
		} else {
			for i, wgt := range w.children {
				r := wgt.Rect()
				if r.Empty() || (!w.floating[i] && !r.Eq(w.rect)) {
					r = w.rect
				}
				if r.Max.X >= w.rect.Max.X {
					r = r.Sub(image.Point{r.Max.X - w.rect.Max.X, 0})
				}
				if r.Max.Y >= w.rect.Max.Y {
					r = r.Sub(image.Point{0, r.Max.Y - w.rect.Max.Y})
				}
				if r.Min.X < w.rect.Min.X {
					r.Min.X = w.rect.Min.X
				}
				if r.Min.Y < w.rect.Min.Y {
					r.Min.Y = w.rect.Min.Y
				}
				wgt.SetRect(r)
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
		t := NewText("")
		t.SetFont(w.font, w.fontSize)

		w.children = append(w.children, &windowWidget{NewBox(), t, widget, w})
		w.titles = append(w.titles, t)
		w.floating = append(w.floating, false)
	}
	w.modified = true
}

// AddChildWithTitle adds a child to the window with the specified window title.
func (w *Window) AddChildWithTitle(wgt Widget, title string) int {
	w.Lock()
	defer w.Unlock()

	t := NewText(title)
	t.SetFont(w.font, w.fontSize)

	w.children = append(w.children, &windowWidget{NewBox(), t, wgt, w})
	w.titles = append(w.titles, t)
	w.floating = append(w.floating, false)

	w.modified = true
	return len(w.children) - 1
}

type windowWidget struct {
	*Box
	title *Text
	wgt   Widget
	w     *Window
}

func (w *windowWidget) SetRect(r image.Rectangle) {
	w.Lock()
	defer w.Unlock()

	w.rect = r
	w.title.SetRect(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+w.w.titleSize))
	w.wgt.SetRect(image.Rect(r.Min.X, r.Min.Y+w.w.titleSize, r.Max.X, r.Max.Y))
}

func (w *windowWidget) Background() color.RGBA {
	return color.RGBA{0, 0, 0, 255}
}

func (w *windowWidget) Draw(screen *ebiten.Image) error {
	w.title.Draw(screen)

	background := w.wgt.Background()
	if background.A != 0 {
		screen.SubImage(w.wgt.Rect()).(*ebiten.Image).Fill(background)
	}
	return w.wgt.Draw(screen)
}
