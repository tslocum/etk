package kibodo

import (
	"bytes"
	"runtime"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TODO test presses registered

func TestKeyboard_Draw(t *testing.T) {
	k := newTestKeyboard()

	// Warm caches
	k.drawBackground()
}

func BenchmarkKeyboard_Draw(b *testing.B) {
	k := newTestKeyboard()

	// Warm caches
	k.drawBackground()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k.drawBackground()
	}
}

func BenchmarkKeyboard_Press(b *testing.B) {
	go func() {
		time.Sleep(2 * time.Second)

		k := newTestKeyboard()

		// Warm caches
		k.drawBackground()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k.drawBackground()
			k.keys[0][0].pressed = true
			k.drawBackground()
			k.keys[0][0].pressed = false
		}
	}()

	runtime.LockOSThread()

	err := ebiten.RunGame(NewDummyGame())
	if err != nil {
		b.Error(err)
	}
}

func newTestKeyboard() *Keyboard {
	k := NewKeyboard(defaultFont(), 64)
	k.SetRect(0, 0, 300, 100)

	return k
}

type DummyGame struct {
	ready bool
}

func (d *DummyGame) Update() error {
	return nil
}

func (d *DummyGame) Draw(screen *ebiten.Image) {
	d.ready = true
}

func (d *DummyGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func NewDummyGame() *DummyGame {
	return &DummyGame{}
}

func defaultFont() *text.GoTextFaceSource {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
	return source
}
