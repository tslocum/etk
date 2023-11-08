package etk

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/llgcode/draw2d/draw2dimg"
)

// Checkbox is a toggleable Checkbox. It automatically resizes itself ensure
// a square shape.
type Checkbox struct {
	*Box

	selected    bool
	checkColor  color.RGBA
	borderSize  int
	borderColor color.RGBA
	baseImg     *image.RGBA
	img         *ebiten.Image
	onSelect    func() error
}

// NewCheckbox returns a new Checkbox widget.
func NewCheckbox(onSelect func() error) *Checkbox {
	return &Checkbox{
		Box:         NewBox(),
		checkColor:  Style.TextColorDark,
		borderSize:  2,
		borderColor: Style.BorderColor,
		onSelect:    onSelect,
	}
}

// SetRect sets the position and size of the Checkbox.
func (c *Checkbox) SetRect(r image.Rectangle) {
	if c.Box.rect.Eq(r) {
		return
	}

	bounds := r.Bounds()
	newSize := bounds.Dx()
	if bounds.Dy() < newSize {
		newSize = bounds.Dy()
	}

	if r.Dx() != newSize {
		r.Max.X = r.Min.X + newSize
	}
	if r.Dy() != newSize {
		r.Max.Y = r.Min.Y + newSize
	}
	c.Box.rect = r

	c.updateImage()

	for _, w := range c.children {
		w.SetRect(r)
	}
}

// SetCheckColor sets the check mark color of the Checkbox.
func (c *Checkbox) SetCheckColor(checkColor color.RGBA) {
	c.checkColor = checkColor
	c.updateImage()
}

// SetBorderColor sets the border color of the Checkbox.
func (c *Checkbox) SetBorderColor(borderColor color.RGBA) {
	c.borderColor = borderColor
	c.updateImage()
}

// Selected returns the selection state of the Checkbox.
func (c *Checkbox) Selected() bool {
	return c.selected
}

// SetSelected sets the Checkbox selection state. The onSelect function is not
// called when the value is set manually via SetSelected.
func (c *Checkbox) SetSelected(selected bool) {
	if c.selected == selected {
		return
	}
	c.selected = selected
	c.updateImage()
}

// HandleKeyboard is called when a keyboard event occurs.
func (c *Checkbox) HandleKeyboard(ebiten.Key, rune) (handled bool, err error) {
	return false, nil
}

// HandleMouse is called when a mouse event occurs.
func (c *Checkbox) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	if !clicked {
		return true, nil
	}

	c.selected = !c.selected
	c.updateImage()

	c.Lock()
	onSelect := c.onSelect
	if onSelect == nil {
		c.Unlock()
		return true, nil
	}
	c.Unlock()

	return true, onSelect()
}

func (c *Checkbox) updateImage() {
	r := c.Rect()
	if r.Empty() {
		c.img = nil
		return
	}

	bounds := r.Bounds()
	newSize := bounds.Dx()
	var initializeImg bool
	if c.img == nil {
		initializeImg = true
	} else {
		imgBounds := c.img.Bounds()
		imgSize := imgBounds.Dx()
		if imgSize != newSize {
			initializeImg = true
		}
	}
	if initializeImg {
		c.baseImg = image.NewRGBA(rectAtOrigin(r))
		c.img = ebiten.NewImage(newSize, newSize)
	}

	// Draw border.
	c.img.Fill(c.borderColor)
	c.img.SubImage(rectAtOrigin(r).Inset(c.borderSize)).(*ebiten.Image).Fill(color.RGBA64{0, 0, 0, 0})

	// Draw check mark.
	if !c.selected {
		return
	}
	gc := draw2dimg.NewGraphicContext(c.baseImg)
	gc.SetStrokeColor(c.checkColor)
	gc.SetLineWidth(float64(c.borderSize))
	gc.MoveTo(0, 0)
	gc.LineTo(float64(r.Dx()), float64(r.Dy()))
	gc.MoveTo(0, float64(r.Dy()))
	gc.LineTo(float64(r.Dx()), 0)
	gc.Stroke()
	c.img.DrawImage(ebiten.NewImageFromImage(c.baseImg), nil)

}

// Draw draws the Checkbox on the screen.
func (c *Checkbox) Draw(screen *ebiten.Image) error {
	if c.img == nil {
		return nil
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.rect.Min.X), float64(c.rect.Min.Y))
	screen.DrawImage(c.img, op)
	return nil
}
