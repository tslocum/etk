package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is a resizable image.
type Sprite struct {
	*Box

	img       *ebiten.Image
	imgBounds image.Rectangle

	thumb       *ebiten.Image
	thumbBounds image.Rectangle

	horizontal Alignment
	vertical   Alignment
}

// NewSprite returns a new Sprite widget.
func NewSprite(img *ebiten.Image) *Sprite {
	return &Sprite{
		Box:        NewBox(),
		img:        img,
		imgBounds:  img.Bounds(),
		horizontal: AlignCenter,
		vertical:   AlignCenter,
	}
}

// SetImage sets the image of the Sprite.
func (s *Sprite) SetImage(img *ebiten.Image) {
	s.Lock()
	defer s.Unlock()

	s.img = img
	s.thumbBounds = image.Rectangle{}
}

// SetHorizontal sets the horizontal alignment of the Sprite.
func (s *Sprite) SetHorizontal(h Alignment) {
	s.Lock()
	defer s.Unlock()

	s.horizontal = h
	s.thumbBounds = image.Rectangle{}
}

// SetVertical sets the vertical alignment of the Sprite.
func (s *Sprite) SetVertical(v Alignment) {
	s.Lock()
	defer s.Unlock()

	s.vertical = v
	s.thumbBounds = image.Rectangle{}
}

// Draw draws the Sprite on the screen.
func (s *Sprite) Draw(screen *ebiten.Image) error {
	s.Lock()
	defer s.Unlock()

	if s.rect.Dx() == 0 || s.rect.Dy() == 0 {
		return nil // The Sprite has no size. Don't draw anything.
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.rect.Min.X), float64(s.rect.Min.Y))
	if s.imgBounds.Dx() == s.rect.Dx() && s.imgBounds.Dy() == s.rect.Dy() {
		screen.DrawImage(s.img, op)
		return nil
	} else if s.thumb == nil || s.thumbBounds.Dx() != s.rect.Dx() || s.thumbBounds.Dy() != s.rect.Dy() {
		scale, yScale := float64(s.rect.Dx())/float64(s.imgBounds.Dx()), float64(s.rect.Dy())/float64(s.imgBounds.Dy())
		if yScale < scale {
			scale = yScale
		}
		thumbOp := &ebiten.DrawImageOptions{}
		thumbOp.GeoM.Scale(scale, scale)
		if s.horizontal != AlignStart {
			delta := float64(s.rect.Dx()) - float64(s.imgBounds.Dx())*scale
			if s.horizontal == AlignCenter {
				thumbOp.GeoM.Translate(delta/2, 0)
			} else { // AlignEnd
				thumbOp.GeoM.Translate(delta, 0)
			}
		}
		if s.vertical != AlignStart {
			delta := float64(s.rect.Dy()) - float64(s.imgBounds.Dy())*scale
			if s.vertical == AlignCenter {
				thumbOp.GeoM.Translate(0, delta/2)
			} else { // AlignEnd
				thumbOp.GeoM.Translate(0, delta)
			}
		}
		createThumb := s.thumb == nil
		if !createThumb {
			bounds := s.thumb.Bounds()
			createThumb = bounds.Dx() != s.rect.Dx() || bounds.Dy() != s.rect.Dy()
		}
		if createThumb {
			s.thumb = ebiten.NewImage(s.rect.Dx(), s.rect.Dy())
		}
		s.thumb.DrawImage(s.img, thumbOp)
	}
	screen.DrawImage(s.thumb, op)
	return nil
}
