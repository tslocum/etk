package messeji

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"sync"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed testdata
var testDataFS embed.FS

var testTextField *TextField

func defaultFont() *text.GoTextFaceSource {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
	return source
}

func cacheGlyphs(source *text.GoTextFaceSource, size int) {
	testFiles := []string{"loremipsum", "long"}
	for _, fileName := range testFiles {
		content, err := testDataFS.ReadFile(fmt.Sprintf("testdata/%s.txt", fileName))
		if err != nil {
			panic("failed to open testdata")
		}
		face := &text.GoTextFace{
			Source: source,
			Size:   float64(size),
		}
		text.CacheGlyphs(string(content), face)
	}
}

func TestWrapContent(t *testing.T) {
	const fontSize = 24
	fontSource := defaultFont()

	cacheGlyphs(fontSource, fontSize)

	testCases := []struct {
		long     bool // Test data type.
		wordWrap bool // Enable wordwrap.

	}{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	testRect := image.Rect(0, 0, 200, 400)

	for _, c := range testCases {
		var name string
		if !c.long {
			name = "loremipsum"
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
			textField := NewTextField(fontSource, fontSize, &sync.Mutex{})
			testTextField = textField
			textField.SetRect(testRect)
			textField.SetWordWrap(c.wordWrap)
			textField.Write(content)
			textField.bufferModified()
		})
	}
}

func BenchmarkWrapContent(b *testing.B) {
	const fontSize = 24
	fontSource := defaultFont()

	cacheGlyphs(fontSource, fontSize)

	testCases := []struct {
		long     bool // Test data type.
		wordWrap bool // Enable wordwrap.

	}{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	testRect := image.Rect(0, 0, 200, 400)

	for _, c := range testCases {
		var name string
		if !c.long {
			name = "loremipsum"
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

		textField := NewTextField(fontSource, fontSize, &sync.Mutex{})
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
