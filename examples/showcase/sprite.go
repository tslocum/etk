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

func newSpriteExample() (string, etk.Widget, etk.Widget) {
	sourceImg, _, err := image.Decode(bytes.NewReader(assetLenna))
	if err != nil {
		log.Fatal(err)
	}

	sprite := etk.NewSprite(ebiten.NewImageFromImage(sourceImg))
	sprite.SetHorizontal(etk.AlignStart)
	sprite.SetVertical(etk.AlignStart)

	return "sprite", sprite, nil
}
