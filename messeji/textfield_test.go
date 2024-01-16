package messeji

import (
	"embed"
	"fmt"
	"image"
	"log"
	"sync"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed testdata
var testDataFS embed.FS

var testTextField *TextField

func TestWrapContent(t *testing.T) {
	testCases := []struct {
		long     bool // Short or long text.
		wordWrap bool // Enable wordwrap.

	}{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const size = 24
	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	testRect := image.Rect(0, 0, 200, 400)

	for _, c := range testCases {
		var name string
		if !c.long {
			name = "short"
		} else {
			name = "long"
		}

		content, err := testDataFS.ReadFile(fmt.Sprintf("testdata/%s.txt", name))
		if err != nil {
			t.Errorf("failed to open testdata: %s", err)
		}

		if !c.wordWrap {
			name += "/wrapChar"
		} else {
			name += "/wrapWord"
		}

		t.Run(name, func(t *testing.T) {
			textField := NewTextField(face, &sync.Mutex{})
			testTextField = textField
			textField.SetRect(testRect)
			textField.SetWordWrap(c.wordWrap)
			textField.Write(content)
			textField.bufferModified()
		})
	}
}

func BenchmarkWrapContent(b *testing.B) {
	testCases := []struct {
		long     bool // Short or long text.
		wordWrap bool // Enable wordwrap.

	}{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const size = 24
	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	testRect := image.Rect(0, 0, 200, 400)

	for _, c := range testCases {
		var name string
		if !c.long {
			name = "short"
		} else {
			name = "long"
		}

		content, err := testDataFS.ReadFile(fmt.Sprintf("testdata/%s.txt", name))
		if err != nil {
			b.Errorf("failed to open testdata: %s", err)
		}

		if !c.wordWrap {
			name += "/wrapChar"
		} else {
			name += "/wrapWord"
		}

		textField := NewTextField(face, &sync.Mutex{})
		testTextField = textField
		textField.SetRect(testRect)
		textField.SetWordWrap(c.wordWrap)

		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				textField.SetText("")
				textField.Write(content)
				textField.bufferModified()
			}
		})
	}
}
