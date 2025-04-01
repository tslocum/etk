//go:build example

package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"codeberg.org/tslocum/etk"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed asset/lenna.png
var assetLenna []byte

func newSpriteExample() (etk.Widget, etk.Widget) {
	sourceImg, _, err := image.Decode(bytes.NewReader(assetLenna))
	if err != nil {
		log.Fatal(err)
	}

	s := etk.NewSprite(ebiten.NewImageFromImage(sourceImg))
	s.SetHorizontal(etk.AlignStart)
	s.SetVertical(etk.AlignStart)
	return s, nil
}
