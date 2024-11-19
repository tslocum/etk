package messeji

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"strings"
	"sync"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Alignment specifies how text is aligned within the field.
type Alignment int

const (
	// AlignStart aligns text at the start of the field.
	AlignStart Alignment = 0

	// AlignCenter aligns text at the center of the field.
	AlignCenter Alignment = 1

	// AlignEnd aligns text at the end of the field.
	AlignEnd Alignment = 2
)

const (
	initialPadding     = 5
	initialScrollWidth = 32
	maxScroll          = 3
)

var (
	initialForeground   = color.RGBA{0, 0, 0, 255}
	initialBackground   = color.RGBA{255, 255, 255, 255}
	initialScrollArea   = color.RGBA{200, 200, 200, 255}
	initialScrollHandle = color.RGBA{108, 108, 108, 255}
)

// TextField is a text display field. Call Update and Draw when your Game's
// Update and Draw methods are called.
//
// Note: A position and size must be set via SetRect before the field will appear.
// Keyboard events are not handled by default, and may be enabled via SetHandleKeyboard.
type TextField struct {
	// r specifies the position and size of the field.
	r image.Rectangle

	// buffer is the text buffer split by newline characters.
	buffer [][]byte

	// incoming is text to be written to the buffer that has not yet been wrapped.
	incoming []byte

	// prefix is the text shown before the content of the field.
	prefix string

	// suffix is the text shown after the content of the field.
	suffix string

	// wordWrap determines whether content is wrapped at word boundaries.
	wordWrap bool

	// bufferWrapped is the content of the field after applying wrapping.
	bufferWrapped []string

	// wrapStart is the first line number in bufferWrapped which corresponds
	// to the last line number in the actual text buffer.
	wrapStart int

	// needWrap is the first line number in the actual text buffer that needs to be wrapped.
	needWrap int

	// wrapScrollBar is whether the scroll bar was visible the last time the field was redrawn.
	wrapScrollBar bool

	// bufferSize is the size (in pixels) of the entire text buffer. When single
	// line mode is enabled,
	bufferSize int

	// lineWidths is the size (in pixels) of each line as it appears on the screen.
	lineWidths []int

	// singleLine is whether the field displays all text on a single line.
	singleLine bool

	// horizontal is the horizontal alignment of the text within field.
	horizontal Alignment

	// vertical is the vertical alignment of the text within field.
	vertical Alignment

	autoResize bool

	// fontSource is the font face source of the text within the field.
	fontSource *text.GoTextFaceSource

	// fontFace is the font face of the text within the field.
	fontFace *text.GoTextFace

	// fontSize is the maximum font size of the text within the field.
	fontSize int

	// overrideFontSize is the actual font size of the text within the field.
	overrideFontSize int

	// fontMutex is the lock which is held whenever utilizing the font.
	fontMutex *sync.Mutex

	// lineHeight is the height of a single line of text.
	lineHeight int

	// overrideLineHeight is the custom height for a line of text, or 0 to disable.
	overrideLineHeight int

	// lineOffset is the offset of the baseline current font.
	lineOffset int

	// textColor is the color of the text within the field.
	textColor color.RGBA

	// backgroundColor is the color of the background of the field.
	backgroundColor color.RGBA

	// padding is the amount of padding around the text within the field.
	padding int

	// follow determines whether the field should automatically scroll to the
	// end when content is added to the buffer.
	follow bool

	// overflow is whether the content of the field is currently larger than the field.
	overflow bool

	// offset is the current view offset of the text within the field, relative to the top.
	offset int

	// handleKeyboard is a flag which, when enabled, causes keyboard input to be handled.
	handleKeyboard bool

	// modified is a flag which, when enabled, causes bufferModified to be called
	// during the next Draw call.
	modified bool

	// scrollRect specifies the position and size of the scrolling area.
	scrollRect image.Rectangle

	// scrollWidth is the width of the scroll bar.
	scrollWidth int

	// scrollAreaColor is the color of the scroll area.
	scrollAreaColor color.RGBA

	// scrollHandleColor is the color of the scroll handle.
	scrollHandleColor color.RGBA

	// scrollBorderSize is the size of the border around the scroll bar handle.
	scrollBorderSize int

	// Scroll bar handle border colors.
	scrollBorderTop    color.RGBA
	scrollBorderRight  color.RGBA
	scrollBorderBottom color.RGBA
	scrollBorderLeft   color.RGBA

	// scrollVisible is whether the scroll bar is visible on the screen.
	scrollVisible bool

	// scrollAutoHide is whether the scroll bar should be automatically hidden
	// when the entire text buffer fits within the screen.
	scrollAutoHide bool

	// scrollDrag is whether the scroll bar is currently being dragged.
	scrollDrag bool

	// scrollDragPoint is the point where the field is being dragged directly.
	scrollDragPoint image.Point

	// scrollDragOffset is the original offset when the field is being dragged directly.
	scrollDragOffset int

	// maskRune is the rune shown instead of the actual buffer contents.
	maskRune rune

	// img is the image of the field.
	img *ebiten.Image

	// visible is whether the field is visible on the screen.
	visible bool

	// redraw is whether the field needs to be redrawn.
	redraw bool

	// keyBuffer is a buffer of key press events.
	keyBuffer []ebiten.Key

	// keyBuffer is a buffer of runes from key presses.
	runeBuffer []rune

	sync.Mutex
}

// NewTextField returns a new TextField. See type documentation for more info.
func NewTextField(fontSource *text.GoTextFaceSource, fontSize int, fontMutex *sync.Mutex) *TextField {
	if fontMutex == nil {
		fontMutex = &sync.Mutex{}
	}

	f := &TextField{
		fontSource:        fontSource,
		fontSize:          fontSize,
		fontMutex:         fontMutex,
		textColor:         initialForeground,
		backgroundColor:   initialBackground,
		padding:           initialPadding,
		scrollWidth:       initialScrollWidth,
		scrollAreaColor:   initialScrollArea,
		scrollHandleColor: initialScrollHandle,
		follow:            true,
		wordWrap:          true,
		scrollVisible:     true,
		scrollAutoHide:    true,
		scrollDragPoint:   image.Point{-1, -1},
		visible:           true,
		redraw:            true,
	}

	f.fontMutex.Lock()
	defer f.fontMutex.Unlock()

	f.resizeFont()
	return f
}

// Rect returns the position and size of the field.
func (f *TextField) Rect() image.Rectangle {
	f.Lock()
	defer f.Unlock()

	return f.r
}

// SetRect sets the position and size of the field.
func (f *TextField) SetRect(r image.Rectangle) {
	f.Lock()
	defer f.Unlock()

	if f.r.Eq(r) {
		return
	}

	if f.r.Dx() != r.Dx() || f.r.Dy() != r.Dy() {
		f.bufferWrapped = f.bufferWrapped[:0]
		f.lineWidths = f.lineWidths[:0]
		f.needWrap = 0
		f.wrapStart = 0
		f.modified = true
	}

	f.r = r
	f.resizeFont()
}

func (f *TextField) text() string {
	f.processIncoming()
	f.resizeFont()
	return string(bytes.Join(f.buffer, []byte("\n")))
}

// Text returns the text in the field.
func (f *TextField) Text() string {
	f.Lock()
	defer f.Unlock()

	return f.text()
}

// SetText sets the text in the field.
func (f *TextField) SetText(text string) {
	f.Lock()
	defer f.Unlock()

	f.buffer = f.buffer[:0]
	f.bufferWrapped = f.bufferWrapped[:0]
	f.lineWidths = f.lineWidths[:0]
	f.needWrap = 0
	f.wrapStart = 0
	f.incoming = append(f.incoming[:0], []byte(text)...)
	f.modified = true
	f.redraw = true
	f.resizeFont()
}

// SetLast sets the text of the last line of the field. Newline characters are
// replaced with spaces.
func (f *TextField) SetLast(text string) {
	f.Lock()
	defer f.Unlock()

	t := bytes.ReplaceAll([]byte(text), []byte("\n"), []byte(" "))

	f.processIncoming()

	bufferLen := len(f.buffer)
	if bufferLen == 0 {
		f.incoming = append(f.incoming[:0], t...)
	} else {
		f.buffer[bufferLen-1] = t
		if f.needWrap == -1 || f.needWrap > bufferLen-1 {
			f.needWrap = bufferLen - 1
		}
	}

	f.modified = true
	f.redraw = true
	f.resizeFont()
}

// SetPrefix sets the text shown before the content of the field.
func (f *TextField) SetPrefix(text string) {
	f.Lock()
	defer f.Unlock()

	f.prefix = text
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// SetSuffix sets the text shown after the content of the field.
func (f *TextField) SetSuffix(text string) {
	f.Lock()
	defer f.Unlock()

	f.suffix = text
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// SetFollow sets whether the field should automatically scroll to the end when
// content is added to the buffer.
func (f *TextField) SetFollow(follow bool) {
	f.Lock()
	defer f.Unlock()

	f.follow = follow
}

// SetSingleLine sets whether the field displays all text on a single line.
// When enabled, the field scrolls horizontally. Otherwise, it scrolls vertically.
func (f *TextField) SetSingleLine(single bool) {
	f.Lock()
	defer f.Unlock()

	if f.singleLine == single {
		return
	}

	f.singleLine = single
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// SetHorizontal sets the horizontal alignment of the text within the field.
func (f *TextField) SetHorizontal(h Alignment) {
	f.Lock()
	defer f.Unlock()

	if f.horizontal == h {
		return
	}

	f.horizontal = h
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
}

// SetVertical sets the veritcal alignment of the text within the field.
func (f *TextField) SetVertical(v Alignment) {
	f.Lock()
	defer f.Unlock()

	if f.vertical == v {
		return
	}

	f.vertical = v
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
}

// LineHeight returns the line height for the field.
func (f *TextField) LineHeight() int {
	f.Lock()
	defer f.Unlock()

	if f.overrideLineHeight != 0 {
		return f.overrideLineHeight
	}
	return f.lineHeight
}

// SetLineHeight sets a custom line height for the field. Setting a line
// height of 0 restores the automatic line height detection based on the font.
func (f *TextField) SetLineHeight(height int) {
	f.Lock()
	defer f.Unlock()

	f.overrideLineHeight = height
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// ForegroundColor returns the color of the text within the field.
func (f *TextField) ForegroundColor() color.RGBA {
	f.Lock()
	defer f.Unlock()

	return f.textColor
}

// SetForegroundColor sets the color of the text within the field.
func (f *TextField) SetForegroundColor(c color.RGBA) {
	f.Lock()
	defer f.Unlock()

	f.textColor = c
	f.modified = true
}

// SetBackgroundColor sets the color of the background of the field.
func (f *TextField) SetBackgroundColor(c color.RGBA) {
	f.Lock()
	defer f.Unlock()

	f.backgroundColor = c
	f.modified = true
}

// SetFont sets the font face of the text within the field.
func (f *TextField) SetFont(faceSource *text.GoTextFaceSource, size int, mutex *sync.Mutex) {
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	f.Lock()
	defer f.Unlock()

	mutex.Lock()
	defer mutex.Unlock()

	f.fontSource = faceSource
	f.fontSize = size
	f.fontMutex = mutex
	f.overrideFontSize = 0

	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// SetAutoResize sets whether the font is automatically scaled down when it is
// too large to fit the entire text buffer on one line.
func (f *TextField) SetAutoResize(resize bool) {
	f.Lock()
	defer f.Unlock()

	f.autoResize = resize
	f.resizeFont()
}

// Padding returns the amount of padding around the text within the field.
func (f *TextField) Padding() int {
	f.Lock()
	defer f.Unlock()

	return f.padding
}

// SetPadding sets the amount of padding around the text within the field.
func (f *TextField) SetPadding(padding int) {
	f.Lock()
	defer f.Unlock()

	f.padding = padding
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// Visible returns whether the field is currently visible on the screen.
func (f *TextField) Visible() bool {
	return f.visible
}

// SetVisible sets whether the field is visible on the screen.
func (f *TextField) SetVisible(visible bool) {
	f.Lock()
	defer f.Unlock()

	if f.visible == visible {
		return
	}

	f.visible = visible
	if visible {
		f.redraw = true
	}
}

// SetScrollBarWidth sets the width of the scroll bar.
func (f *TextField) SetScrollBarWidth(width int) {
	f.Lock()
	defer f.Unlock()

	if f.scrollWidth == width {
		return
	}

	f.scrollWidth = width
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// SetScrollBarColors sets the color of the scroll bar area and handle.
func (f *TextField) SetScrollBarColors(area color.RGBA, handle color.RGBA) {
	f.Lock()
	defer f.Unlock()

	f.scrollAreaColor, f.scrollHandleColor = area, handle
	f.redraw = true
}

// SetScrollBorderSize sets the size of the border around the scroll bar handle.
func (f *TextField) SetScrollBorderSize(size int) {
	f.Lock()
	defer f.Unlock()

	f.scrollBorderSize = size
	f.redraw = true
	f.resizeFont()
}

// SetScrollBorderColor sets the color of the top, right, bottom and left border
// of the scroll bar handle.
func (f *TextField) SetScrollBorderColors(top color.RGBA, right color.RGBA, bottom color.RGBA, left color.RGBA) {
	f.Lock()
	defer f.Unlock()

	f.scrollBorderTop = top
	f.scrollBorderRight = right
	f.scrollBorderBottom = bottom
	f.scrollBorderLeft = left
	f.redraw = true
}

// SetScrollBarVisible sets whether the scroll bar is visible on the screen.
func (f *TextField) SetScrollBarVisible(scrollVisible bool) {
	f.Lock()
	defer f.Unlock()

	if f.scrollVisible == scrollVisible {
		return
	}

	f.scrollVisible = scrollVisible
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// SetAutoHideScrollBar sets whether the scroll bar is automatically hidden
// when the entire text buffer is visible.
func (f *TextField) SetAutoHideScrollBar(autoHide bool) {
	f.Lock()
	defer f.Unlock()

	if f.scrollAutoHide == autoHide {
		return
	}

	f.scrollAutoHide = autoHide
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
	f.resizeFont()
}

// WordWrap returns the current text wrap mode.
func (f *TextField) WordWrap() bool {
	f.Lock()
	defer f.Unlock()

	return f.wordWrap
}

// SetWordWrap sets a flag which, when enabled, causes text to wrap without breaking words.
func (f *TextField) SetWordWrap(wrap bool) {
	f.Lock()
	defer f.Unlock()

	if f.wordWrap == wrap {
		return
	}

	f.wordWrap = wrap
	f.needWrap = 0
	f.wrapStart = 0
	f.modified = true
}

func (f *TextField) resizeFont() {
	if !f.autoResize {
		if f.overrideFontSize == f.fontSize {
			return
		}
		f.overrideFontSize = f.fontSize
		f.fontFace = fontFace(f.fontSource, f.overrideFontSize)
		f.fontUpdated()
		f.bufferModified()
		return
	}

	w, h := f.r.Dx()-f.padding*2, f.r.Dy()-f.padding*2
	if w <= 0 || h <= 0 {
		if f.overrideFontSize == f.fontSize {
			return
		}
		f.overrideFontSize = f.fontSize
		f.fontFace = fontFace(f.fontSource, f.overrideFontSize)
		f.fontUpdated()
		f.bufferModified()
		return
	}

	f.processIncoming()

	for size := f.fontSize; size > 0; size-- {
		f.fontFace = fontFace(f.fontSource, size)
		f.fontUpdated()
		if f.lineHeight > h {
			continue
		}
		f.needWrap = 0
		f.wrapStart = 0
		f.wrap()
		if len(f.bufferWrapped) <= 1 {
			break
		}
	}
}

// SetHandleKeyboard sets a flag controlling whether keyboard input should be handled
// by the field. This can be used to facilitate focus changes between multiple inputs.
func (f *TextField) SetHandleKeyboard(handle bool) {
	f.Lock()
	defer f.Unlock()

	f.handleKeyboard = handle
}

// SetMask sets the rune used to mask the text buffer contents. Set to 0 to disable.
func (f *TextField) SetMask(r rune) {
	f.Lock()
	defer f.Unlock()

	if f.maskRune == r {
		return
	}

	f.maskRune = r
	f.modified = true
	f.resizeFont()
}

// Write writes to the field's buffer.
func (f *TextField) Write(p []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	return f._write(p)
}

func (f *TextField) _write(p []byte) (n int, err error) {
	f.incoming = append(f.incoming, p...)
	f.modified = true
	f.redraw = true
	return len(p), nil
}

// HandleKeyboardEvent passes the provided key or rune to the TextField.
func (f *TextField) HandleKeyboardEvent(key ebiten.Key, r rune) (handled bool, err error) {
	f.Lock()
	defer f.Unlock()

	if !f.visible || rectIsZero(f.r) || !f.handleKeyboard {
		return false, nil
	}

	return f._handleKeyboardEvent(key, r)
}

func (f *TextField) _handleKeyboardEvent(key ebiten.Key, r rune) (handled bool, err error) {
	if key != -1 {
		// Handle keyboard PageUp/PageDown.
		offsetAmount := 0
		switch key {
		case ebiten.KeyPageUp:
			offsetAmount = 100
		case ebiten.KeyPageDown:
			offsetAmount = -100
		}
		if offsetAmount != 0 {
			f.offset += offsetAmount
			f.clampOffset()
			f.redraw = true
			return true, nil
		}
		return true, err
	}
	return true, nil
}

func (f *TextField) HandleMouseEvent(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	f.Lock()
	defer f.Unlock()

	if !f.visible || rectIsZero(f.r) {
		return false, nil
	}

	return f._handleMouseEvent(cursor, pressed, clicked)
}

func (f *TextField) _handleMouseEvent(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	if !cursor.In(f.r) {
		return false, nil
	}

	// Handle mouse wheel.
	_, scroll := ebiten.Wheel()
	if scroll != 0 {
		if scroll < -maxScroll {
			scroll = -maxScroll
		} else if scroll > maxScroll {
			scroll = maxScroll
		}
		const offsetAmount = 25
		f.offset += int(scroll * offsetAmount)
		f.clampOffset()
		f.redraw = true
	}

	// Handle scroll bar click (and drag).
	if !f.showScrollBar() {
		return true, nil
	} else if pressed || f.scrollDrag {
		p := image.Point{cursor.X - f.r.Min.X, cursor.Y - f.r.Min.Y}
		if pressed {
			// Handle dragging the text field directly.
			if !f.scrollDrag && !p.In(f.scrollRect) && f.scrollDragPoint.X == -1 && f.scrollDragPoint.Y == -1 {
				f.scrollDragPoint = p
				f.scrollDragOffset = f.offset
			}
			if f.scrollDragPoint.X != -1 {
				delta := f.scrollDragPoint.Y - p.Y
				f.offset = f.scrollDragOffset - delta
			} else { // Handle dragging the scroll bar handle.
				dragY := cursor.Y - f.r.Min.Y - f.scrollWidth/4
				if dragY < 0 {
					dragY = 0
				} else if dragY > f.scrollRect.Dy() {
					dragY = f.scrollRect.Dy()
				}

				pct := float64(dragY) / float64(f.scrollRect.Dy()-f.scrollWidth/2)
				if pct < 0 {
					pct = 0
				} else if pct > 1 {
					pct = 1
				}

				h := f.r.Dy()
				f.offset = -int(float64(f.bufferSize-h-f.lineOffset+f.padding*2) * pct)
			}
			f.clampOffset()

			f.redraw = true
			f.scrollDrag = true
		} else if !pressed {
			f.scrollDrag = false
			f.scrollDragPoint = image.Point{-1, -1}
			f.scrollDragOffset = 0
		}
	}
	return true, nil
}

// Update updates the field. This function should be called when
// Game.Update is called.
func (f *TextField) Update() error {
	f.Lock()
	defer f.Unlock()

	if !f.visible || rectIsZero(f.r) {
		return nil
	}

	f.keyBuffer = inpututil.AppendJustPressedKeys(f.keyBuffer[:0])
	for _, key := range f.keyBuffer {
		handled, err := f._handleKeyboardEvent(key, 0)
		if err != nil {
			return err
		} else if handled {
			f.redraw = true
		}
	}

	f.runeBuffer = ebiten.AppendInputChars(f.runeBuffer[:0])
	for _, r := range f.runeBuffer {
		handled, err := f._handleKeyboardEvent(-1, r)
		if err != nil {
			return err
		} else if handled {
			f.redraw = true
		}
	}

	cx, cy := ebiten.CursorPosition()
	if cx != 0 || cy != 0 {
		handled, err := f._handleMouseEvent(image.Point{X: cx, Y: cy}, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft), inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))
		if err != nil {
			return err
		} else if handled {
			f.redraw = true
		}
	}

	return nil
}

// Draw draws the field on the screen. This function should be called
// when Game.Draw is called.
func (f *TextField) Draw(screen *ebiten.Image) {
	f.Lock()
	defer f.Unlock()

	if f.modified {
		f.fontMutex.Lock()

		f.bufferModified()
		f.modified = false

		f.fontMutex.Unlock()
	}

	if !f.visible || rectIsZero(f.r) {
		return
	}

	if f.redraw {
		f.fontMutex.Lock()

		f.drawImage()
		f.redraw = false

		f.fontMutex.Unlock()
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.r.Min.X), float64(f.r.Min.Y))
	screen.DrawImage(f.img, op)
}

func (f *TextField) fontUpdated() {
	m := f.fontFace.Metrics()
	f.lineHeight = int(m.HAscent + m.HDescent)
	f.lineOffset = int(m.HLineGap)
	if f.lineOffset < 0 {
		f.lineOffset *= -1
	}
}

func (f *TextField) wrapContent(withScrollBar bool) {
	if withScrollBar != f.wrapScrollBar {
		f.needWrap = 0
		f.wrapStart = 0
	} else if f.needWrap == -1 {
		return
	}
	f.wrapScrollBar = withScrollBar

	if len(f.buffer) == 0 || (f.singleLine && !f.autoResize) {
		buffer := f.prefix + string(bytes.Join(f.buffer, nil)) + f.suffix
		w, _ := text.Measure(buffer, f.fontFace, float64(f.lineHeight))

		f.bufferWrapped = []string{buffer}
		f.wrapStart = 0
		f.lineWidths = append(f.lineWidths[:0], int(w))

		f.needWrap = -1
		return
	}

	w := f.r.Dx()
	if withScrollBar {
		w -= f.scrollWidth
	}
	bufferLen := len(f.buffer)
	j := f.wrapStart
	for i := f.needWrap; i < bufferLen; i++ {
		var line string
		if i == 0 {
			line = f.prefix + string(f.buffer[i])
		} else {
			line = string(f.buffer[i])
		}
		if i == bufferLen-1 {
			line += f.suffix
		}
		l := len(line)
		availableWidth := w - (f.padding * 2)

		f.wrapStart = j

		// BoundString returns 0 for strings containing only whitespace.
		if strings.TrimSpace(line) == "" {
			if len(f.bufferWrapped) <= j {
				f.bufferWrapped = append(f.bufferWrapped, "")
			} else {
				f.bufferWrapped[j] = ""
			}
			if len(f.lineWidths) <= j {
				f.lineWidths = append(f.lineWidths, 0)
			} else {
				f.lineWidths[j] = 0
			}
			j++
			continue
		}

		// Add characters one at a time until the line doesn't fit. When word
		// wrapping is enabled, break the line at the last whitespace character.
		var start int
		var lastSpace int
		var lastSpaceSize int
		var boundsWidth int
		var lastBoundsWidth int
	WRAPTEXT:
		for start < l {
			lastSpace = -1
			lastBoundsWidth = -1
			var e int
			for _, r := range line[start:] {
				runeSize := len(string(r))
				if e > l-start-runeSize {
					e = l - start - runeSize
				}
				if unicode.IsSpace(r) {
					lastSpace = e
					lastSpaceSize = runeSize
				}
				w, _ := text.Measure(line[start:start+e+runeSize], f.fontFace, float64(f.lineHeight))
				boundsWidth = int(w)
				if boundsWidth > availableWidth {
					var addSpace bool
					if e > 0 {
						if e > 0 {
							e -= runeSize
						}
						if f.wordWrap && lastSpace != -1 {
							e = lastSpace
							addSpace = true
						}
					}
					if lastBoundsWidth == -1 {
						w, _ := text.Measure(line[start:start+e], f.fontFace, float64(f.lineHeight))
						boundsWidth = int(w)
					} else {
						boundsWidth = lastBoundsWidth
					}

					if len(f.bufferWrapped) <= j {
						f.bufferWrapped = append(f.bufferWrapped, line[start:start+e])
					} else {
						f.bufferWrapped[j] = line[start : start+e]
					}
					if len(f.lineWidths) <= j {
						f.lineWidths = append(f.lineWidths, boundsWidth)
					} else {
						f.lineWidths[j] = boundsWidth
					}
					j++

					if addSpace {
						e += lastSpaceSize
					}

					start += e
					if e == 0 {
						start += runeSize
					}
					continue WRAPTEXT
				}
				lastBoundsWidth = boundsWidth
				e += runeSize
			}

			if len(f.bufferWrapped) <= j {
				f.bufferWrapped = append(f.bufferWrapped, line[start:])
			} else {
				f.bufferWrapped[j] = line[start:]
			}
			if len(f.lineWidths) <= j {
				f.lineWidths = append(f.lineWidths, boundsWidth)
			} else {
				f.lineWidths[j] = boundsWidth
			}
			j++
			break
		}
	}

	if len(f.bufferWrapped) >= j {
		f.bufferWrapped = f.bufferWrapped[:j]
	}

	f.needWrap = -1
}

// drawContent draws the text buffer to img.
func (f *TextField) drawContent() (overflow bool) {
	if f.backgroundColor.A != 0 {
		f.img.Fill(f.backgroundColor)
	} else {
		f.img.Clear()
	}
	fieldWidth := f.r.Dx()
	fieldHeight := f.r.Dy()
	if f.showScrollBar() {
		fieldWidth -= f.scrollWidth
	}
	lines := len(f.bufferWrapped)

	h := f.r.Dy()
	lineHeight := f.overrideLineHeight
	if lineHeight == 0 {
		lineHeight = f.lineHeight
	}
	var firstVisible, lastVisible int
	firstVisible = 0
	lastVisible = len(f.bufferWrapped) - 1
	if !f.singleLine {
		firstVisible = (f.offset * -1) / f.lineHeight
		lastVisible = firstVisible + (f.r.Dy() / f.lineHeight) + 1
		if lastVisible > len(f.bufferWrapped)-1 {
			lastVisible = len(f.bufferWrapped) - 1
		}
	}
	numVisible := lastVisible - firstVisible
	// Calculate buffer size (width for single-line fields or height for multi-line fields).
	if f.singleLine {
		w, _ := text.Measure(f.bufferWrapped[firstVisible], f.fontFace, float64(f.lineHeight))
		f.bufferSize = int(w)
		if f.bufferSize > fieldWidth-f.padding*2 {
			overflow = true
		}
	} else {
		f.bufferSize = (len(f.bufferWrapped)) * lineHeight
		if f.bufferSize > fieldHeight-f.padding*2 {
			overflow = true
		}
	}
	for i := firstVisible; i <= lastVisible; i++ {
		line := f.bufferWrapped[i]
		if f.maskRune != 0 {
			line = strings.Repeat(string(f.maskRune), len(line))
			if i == lastVisible && len(line) > 0 && len(line) >= len(f.suffix) {
				line = line[:len(line)-len(f.suffix)] + f.suffix
			}
		}
		lineX := f.padding
		lineY := 1 + f.padding + -f.lineOffset + lineHeight*i

		// Calculate whether the line overflows the visible area.
		lineOverflows := lineY < 0 || lineY >= h-f.padding
		if lineOverflows {
			overflow = true
		}

		// Skip drawing off-screen lines.
		if lineY < 0 {
			continue
		}

		// Apply scrolling transformation.
		if f.singleLine {
			lineX += f.offset
		} else {
			lineY += f.offset
		}

		// Align horizontally.
		if f.horizontal == AlignCenter {
			lineX = (fieldWidth - f.lineWidths[i]) / 2
		} else if f.horizontal == AlignEnd {
			lineX = (fieldWidth - f.lineWidths[i]) - f.padding - 1
		}

		// Align vertically.
		totalHeight := f.lineOffset + lineHeight*(lines)
		if f.vertical == AlignCenter && (f.autoResize || totalHeight <= h) {
			lineY = fieldHeight/2 - totalHeight/2 + f.lineOffset + (lineHeight * (i)) - 2
		} else if f.vertical == AlignEnd && (f.autoResize || totalHeight <= h) {
			lineY = fieldHeight - lineHeight*(numVisible+1-i) - f.padding
		}

		// Draw line.
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(lineX), float64(lineY))
		op.ColorScale.ScaleWithColor(f.textColor)
		text.Draw(f.img, line, f.fontFace, op)
	}

	return overflow
}

func (f *TextField) clampOffset() {
	fieldSize := f.r.Dy()
	if f.singleLine {
		fieldSize = f.r.Dx()
	}
	minSize := -(f.bufferSize - fieldSize + f.padding*2)
	if !f.singleLine {
		minSize += f.lineOffset
	}
	if minSize > 0 {
		minSize = 0
	}
	maxSize := 0
	if f.offset < minSize {
		f.offset = minSize
	} else if f.offset > maxSize {
		f.offset = maxSize
	}
}

func (f *TextField) showScrollBar() bool {
	return !f.autoResize && !f.singleLine && f.scrollVisible && (f.overflow || !f.scrollAutoHide)
}

func (f *TextField) wrap() {
	w, h := f.r.Dx(), f.r.Dy()

	var newImage bool
	if f.img == nil {
		newImage = true
	} else {
		imgRect := f.img.Bounds()
		imgW, imgH := imgRect.Dx(), imgRect.Dy()
		newImage = imgW != w || imgH != h
	}
	if newImage {
		f.img = ebiten.NewImage(w, h)
	}

	showScrollBar := f.showScrollBar()
	f.wrapContent(showScrollBar)
	f.overflow = f.drawContent()
	if f.showScrollBar() != showScrollBar {
		f.wrapContent(!showScrollBar)
		f.drawContent()
	}
}

// drawImage draws the field to img (caching it for future draws).
func (f *TextField) drawImage() {
	if rectIsZero(f.r) {
		f.img = nil
		return
	}

	f.wrap()

	// Draw scrollbar.
	if f.showScrollBar() {
		w, h := f.r.Dx(), f.r.Dy()

		scrollAreaX, scrollAreaY := w-f.scrollWidth, 0
		f.scrollRect = image.Rect(scrollAreaX, scrollAreaY, scrollAreaX+f.scrollWidth, h)

		scrollBarH := f.scrollWidth / 2
		if scrollBarH < 4 {
			scrollBarH = 4
		}

		scrollX, scrollY := w-f.scrollWidth, 0
		pct := float64(-f.offset) / float64(f.bufferSize-h-f.lineOffset+f.padding*2)
		scrollY += int(float64(h-scrollBarH) * pct)
		scrollBarRect := image.Rect(scrollX, scrollY, scrollX+f.scrollWidth, scrollY+scrollBarH)

		// Draw scroll area.
		f.img.SubImage(f.scrollRect).(*ebiten.Image).Fill(f.scrollAreaColor)

		// Draw scroll handle.
		f.img.SubImage(scrollBarRect).(*ebiten.Image).Fill(f.scrollHandleColor)

		// Draw scroll handle border.
		if f.scrollBorderSize != 0 {
			r := scrollBarRect
			f.img.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+f.scrollBorderSize, r.Max.Y)).(*ebiten.Image).Fill(f.scrollBorderLeft)
			f.img.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+f.scrollBorderSize)).(*ebiten.Image).Fill(f.scrollBorderTop)
			f.img.SubImage(image.Rect(r.Max.X-f.scrollBorderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(f.scrollBorderRight)
			f.img.SubImage(image.Rect(r.Min.X, r.Max.Y-f.scrollBorderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(f.scrollBorderBottom)
		}
	}
}

func (f *TextField) processIncoming() {
	if len(f.incoming) == 0 {
		return
	}

	line := len(f.buffer) - 1
	if line < 0 {
		line = 0
		f.buffer = append(f.buffer, nil)
	}
	if f.needWrap == -1 {
		f.needWrap = line
	}
	for _, b := range f.incoming {
		if b == '\n' {
			line++
			f.buffer = append(f.buffer, nil)
			continue
		}
		f.buffer[line] = append(f.buffer[line], b)
	}
	f.incoming = f.incoming[:0]
}

func (f *TextField) bufferModified() {
	f.processIncoming()
	f.resizeFont()

	f.drawImage()

	lastOffset := f.offset
	if f.follow {
		f.offset = -math.MaxInt
	}
	f.clampOffset()
	if f.offset != lastOffset {
		f.drawImage()
	}

	f.redraw = false
}

func rectIsZero(r image.Rectangle) bool {
	return r.Dx() == 0 || r.Dy() == 0
}

func fontFace(source *text.GoTextFaceSource, size int) *text.GoTextFace {
	return &text.GoTextFace{
		Source: source,
		Size:   float64(size),
	}
}
