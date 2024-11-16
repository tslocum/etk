package etk

import (
	"image"
	"image/color"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Button is a clickable button.
type Button struct {
	*Box
	field         *messeji.TextField
	btnBackground color.RGBA
	textFont      *text.GoTextFaceSource
	textSize      int
	borderSize    int
	borderTop     color.RGBA
	borderRight   color.RGBA
	borderBottom  color.RGBA
	borderLeft    color.RGBA
	onSelected    func() error
	pressed       bool
}

// NewButton returns a new Button widget.
func NewButton(label string, onSelected func() error) *Button {
	textColor := Style.ButtonTextColor
	if textColor.A == 0 {
		textColor = Style.TextColorDark
	}
	f := newText()
	f.SetText(label)
	f.SetForegroundColor(textColor)
	f.SetHorizontal(messeji.AlignCenter)
	f.SetVertical(messeji.AlignCenter)
	f.SetAutoResize(true)

	b := &Button{
		Box:           NewBox(),
		field:         f,
		btnBackground: Style.ButtonBgColor,
		textFont:      Style.TextFont,
		textSize:      Scale(Style.TextSize),
		onSelected:    onSelected,
		borderSize:    Scale(Style.ButtonBorderSize),
		borderTop:     Style.ButtonBorderTop,
		borderRight:   Style.ButtonBorderRight,
		borderBottom:  Style.ButtonBorderBottom,
		borderLeft:    Style.ButtonBorderLeft,
	}
	b.SetBackground(Style.ButtonBgColor)
	return b
}

// SetRect sets the position and size of the Button.
func (b *Button) SetRect(r image.Rectangle) {
	b.Box.rect = r

	b.field.SetRect(r)

	for _, w := range b.children {
		w.SetRect(r)
	}
}

// SetBorderSize sets the size of the border around the button.
func (b *Button) SetBorderSize(size int) {
	b.Lock()
	defer b.Unlock()

	b.borderSize = size
}

// SetBorderColors sets the color of the top, right, bottom and left border.
func (b *Button) SetBorderColors(top color.RGBA, right color.RGBA, bottom color.RGBA, left color.RGBA) {
	b.Lock()
	defer b.Unlock()

	b.borderTop = top
	b.borderRight = right
	b.borderBottom = bottom
	b.borderLeft = left
}

// SetForeground sets the color of the button label.
func (b *Button) SetForeground(c color.RGBA) {
	b.Lock()
	defer b.Unlock()

	b.field.SetForegroundColor(c)
}

// SetBackground sets the background color of the button label.
func (b *Button) SetBackground(background color.RGBA) {
	b.Lock()
	defer b.Unlock()

	b.btnBackground = background
	b.background = background
}

// Text returns the content of the text buffer.
func (b *Button) Text() string {
	b.Lock()
	defer b.Unlock()

	return b.field.Text()
}

// SetText sets the text in the field.
func (b *Button) SetText(text string) {
	b.Lock()
	defer b.Unlock()

	b.field.SetText(text)
}

// SetFont sets the font and text size of button label. Scaling is not applied.
func (b *Button) SetFont(fnt *text.GoTextFaceSource, size int) {
	b.Lock()
	defer b.Unlock()

	b.textFont, b.textSize = fnt, size
	b.field.SetFont(b.textFont, b.textSize, fontMutex)
}

// SetHorizontal sets the horizontal alignment of the button label.
func (b *Button) SetHorizontal(h Alignment) {
	b.Lock()
	defer b.Unlock()

	b.field.SetHorizontal(messeji.Alignment(h))
}

// SetVertical sets the vertical alignment of the button label.
func (b *Button) SetVertical(v Alignment) {
	b.Lock()
	defer b.Unlock()

	b.field.SetVertical(messeji.Alignment(v))
}

// Cursor returns the cursor shape shown when a mouse cursor hovers over the
// widget, or -1 to let widgets beneath determine the cursor shape.
func (b *Button) Cursor() ebiten.CursorShapeType {
	return ebiten.CursorShapePointer
}

// HandleKeyboard is called when a keyboard event occurs.
func (b *Button) HandleKeyboard(ebiten.Key, rune) (handled bool, err error) {
	return false, nil
}

// HandleMouse is called when a mouse event occurs.
func (b *Button) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	if !clicked {
		if b.pressed && !pressed {
			b.Lock()
			b.pressed = false
			b.background = b.btnBackground
			b.Unlock()
		}
		return true, nil
	}

	const dim = 0.9
	b.Lock()
	b.pressed = true
	b.background = color.RGBA{uint8(float64(b.btnBackground.R) * dim), uint8(float64(b.btnBackground.G) * dim), uint8(float64(b.btnBackground.B) * dim), 255}
	onSelected := b.onSelected
	if onSelected == nil {
		b.Unlock()
		return true, nil
	}
	b.Unlock()

	return true, onSelected()
}

// Draw draws the button on the screen.
func (b *Button) Draw(screen *ebiten.Image) error {
	r := b.rect

	// Draw label.
	b.field.Draw(screen)

	// Draw border.
	if b.borderSize != 0 {
		if !b.pressed {
			screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+b.borderSize, r.Max.Y)).(*ebiten.Image).Fill(b.borderLeft)
			screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+b.borderSize)).(*ebiten.Image).Fill(b.borderTop)
			screen.SubImage(image.Rect(r.Max.X-b.borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(b.borderRight)
			screen.SubImage(image.Rect(r.Min.X, r.Max.Y-b.borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(b.borderBottom)
		} else {
			screen.SubImage(image.Rect(r.Max.X-b.borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(b.borderLeft)
			screen.SubImage(image.Rect(r.Min.X, r.Max.Y-b.borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(b.borderTop)
			screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+b.borderSize, r.Max.Y)).(*ebiten.Image).Fill(b.borderRight)
			screen.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+b.borderSize)).(*ebiten.Image).Fill(b.borderBottom)
		}
	}

	return nil
}
