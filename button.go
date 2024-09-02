package etk

import (
	"image"
	"image/color"

	"code.rocket9labs.com/tslocum/etk/messeji"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
)

// Button is a clickable button.
type Button struct {
	*Box
	field        *messeji.TextField
	textFont     *sfnt.Font
	textSize     int
	textAutoSize int
	borderSize   int
	borderTop    color.RGBA
	borderRight  color.RGBA
	borderBottom color.RGBA
	borderLeft   color.RGBA
	onSelected   func() error
	pressed      bool
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
	f.SetScrollBarVisible(false)

	b := &Button{
		Box:          NewBox(),
		field:        f,
		textFont:     Style.TextFont,
		textSize:     Scale(Style.TextSize),
		onSelected:   onSelected,
		borderSize:   Scale(Style.BorderSize),
		borderTop:    Style.BorderColorTop,
		borderRight:  Style.BorderColorRight,
		borderBottom: Style.BorderColorBottom,
		borderLeft:   Style.BorderColorLeft,
	}
	b.SetBackground(Style.ButtonBgColor)
	b.resizeFont()
	return b
}

// SetRect sets the position and size of the Button.
func (b *Button) SetRect(r image.Rectangle) {
	b.Box.rect = r

	b.field.SetRect(r)
	b.resizeFont()

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
	b.resizeFont()
}

// SetFont sets the font and text size of button label. Scaling is not applied.
func (b *Button) SetFont(fnt *sfnt.Font, size int) {
	b.Lock()
	defer b.Unlock()

	b.textFont, b.textSize = fnt, size
	b.resizeFont()
}

func (b *Button) resizeFont() {
	w, h := b.rect.Dx()-b.field.Padding()*2, b.rect.Dy()-b.field.Padding()*2
	if w == 0 || h == 0 {
		if b.textAutoSize == b.textSize {
			return
		}
		b.textAutoSize = b.textSize
		ff := FontFace(b.textFont, b.textSize)
		b.field.SetFont(ff, fontMutex)
		return
	}

	var autoSize int
	var ff font.Face
	for autoSize = b.textSize; autoSize > 0; autoSize-- {
		ff = FontFace(b.textFont, autoSize)
		bounds := BoundString(ff, b.field.Text())
		if bounds.Dx() <= w && bounds.Dy() <= h {
			break
		}
	}
	if b.textAutoSize == autoSize {
		return
	}

	b.field.SetFont(ff, fontMutex)
	b.textAutoSize = autoSize
}

// SetHorizontal sets the horizontal alignment of the button label.
func (b *Button) SetHorizontal(h Alignment) {
	b.Lock()
	defer b.Unlock()

	b.field.SetHorizontal(messeji.Alignment(h))
}

// SetVertical sets the vertical alignment of the button label.
func (b *Button) SetVertical(h Alignment) {
	b.Lock()
	defer b.Unlock()

	b.field.SetVertical(messeji.Alignment(h))
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
			b.background = Style.ButtonBgColor
			b.Unlock()
		}
		return true, nil
	}

	b.Lock()
	b.pressed = true
	b.background = color.RGBA{uint8(float64(Style.ButtonBgColor.R) * 0.95), uint8(float64(Style.ButtonBgColor.G) * 0.95), uint8(float64(Style.ButtonBgColor.B) * 0.95), 255}
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
