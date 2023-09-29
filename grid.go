package etk

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Grid struct {
	*Box

	columnSizes []int
	rowSizes    []int

	columnPadding int
	rowPadding    int

	cellPositions [][2]int
	cellSpans     [][2]int

	childrenUpdated bool
}

func NewGrid() *Grid {
	return &Grid{
		Box: NewBox(),
	}
}

func (g *Grid) SetRect(r image.Rectangle) {
	g.Lock()
	defer g.Unlock()

	g.Box.rect = r
	g.reposition()
}

func (g *Grid) SetColumnSizes(size ...int) {
	g.Lock()
	defer g.Unlock()

	g.columnSizes = size
	g.reposition()
}

func (g *Grid) SetColumnPadding(padding int) {
	g.Lock()
	defer g.Unlock()

	g.columnPadding = padding
	g.reposition()
}

func (g *Grid) SetRowSizes(size ...int) {
	g.Lock()
	defer g.Unlock()

	g.rowSizes = size
	g.reposition()
}

func (g *Grid) SetRowPadding(padding int) {
	g.Lock()
	defer g.Unlock()

	g.rowPadding = padding
	g.reposition()
}

func (g *Grid) AddChild(wgt ...Widget) {
	g.Box.AddChild(wgt...)

	for i := 0; i < len(wgt); i++ {
		g.cellPositions = append(g.cellPositions, [2]int{0, 0})
		g.cellSpans = append(g.cellSpans, [2]int{1, 1})
	}

	g.childrenUpdated = true
}

func (g *Grid) AddChildAt(wgt Widget, x int, y int, columns int, rows int) {
	g.Box.AddChild(wgt)

	g.cellPositions = append(g.cellPositions, [2]int{x, y})
	g.cellSpans = append(g.cellSpans, [2]int{columns, rows})

	g.childrenUpdated = true
}

func (g *Grid) HandleMouse(cursor image.Point, pressed bool, clicked bool) (handled bool, err error) {
	return false, nil
}

func (g *Grid) HandleKeyboard() (handled bool, err error) {
	return false, nil
}

func (g *Grid) Draw(screen *ebiten.Image) error {
	g.Lock()
	defer g.Unlock()

	for _, child := range g.children {
		err := child.Draw(screen)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Grid) reposition() {
	gridW, gridH := g.rect.Dx(), g.rect.Dy()

	// Determine max column and row sizes and proportions.
	var (
		numColumns int
		numRows    int

		maxColumnProportion = 1
		maxRowProportion    = 1

		numColumnProportions = make(map[int]int)
		numRowProportions    = make(map[int]int)
	)
	for i := range g.children {
		position := g.cellPositions[i]
		x, y := position[0], position[1]

		span := g.cellSpans[i]
		w, h := span[0], span[1]

		if x+w > numColumns {
			numColumns = x + w
		}
		if y+h > numRows {
			numRows = y + h
		}

		if -w > maxColumnProportion {
			maxColumnProportion = -w
		}
		if -h > maxRowProportion {
			maxRowProportion = -h
		}
	}

	// Determine actual column and row sizes and proportions.
	numColumnSizes := len(g.columnSizes)
	numRowSizes := len(g.rowSizes)

	columnWidths := make([]int, numColumns)
	var usedWidth int
	for i := 0; i < numColumns; i++ {
		if i >= numColumnSizes {
			columnWidths[i] = -1
		} else {
			columnWidths[i] = g.columnSizes[i]

			if g.columnSizes[i] > 0 {
				usedWidth += g.columnSizes[i]
			}
		}

		if columnWidths[i] < 0 {
			numColumnProportions[-columnWidths[i]]++
		}
	}
	remainingWidth := gridW - usedWidth - (g.columnPadding * (numColumns + 1))
	columnProportions := make([]int, maxColumnProportion)
	for i := 0; i < maxColumnProportion; i++ {
		columnProportions[i] = remainingWidth / (i + 1)
	}
	for i := 0; i < numColumns; i++ {
		if columnWidths[i] < 0 {
			columnWidths[i] = columnProportions[-columnWidths[i]-1] / numColumnProportions[-columnWidths[i]]
		}
	}

	rowHeights := make([]int, numRows)
	var usedHeight int
	for i := 0; i < numRows; i++ {
		if i >= numRowSizes {
			rowHeights[i] = -1
		} else {
			rowHeights[i] = g.rowSizes[i]

			if g.rowSizes[i] > 0 {
				usedHeight += g.rowSizes[i]
			}
		}

		if rowHeights[i] < 0 {
			numRowProportions[-rowHeights[i]]++
		}
	}
	remainingHeight := gridH - usedHeight - (g.rowPadding * (numRows + 1))
	rowProportions := make([]int, maxRowProportion)
	for i := 0; i < maxRowProportion; i++ {
		rowProportions[i] = remainingHeight / (i + 1)
	}
	for i := 0; i < numRows; i++ {
		if rowHeights[i] < 0 {
			rowHeights[i] = rowProportions[-rowHeights[i]-1] / numRowProportions[-rowHeights[i]]
		}
	}

	columnPositions := make([]int, numColumns)
	{
		x := g.columnPadding
		for i := 0; i < numColumns; i++ {
			columnPositions[i] = x
			x += columnWidths[i] + g.columnPadding
		}
	}

	rowPositions := make([]int, numRows)
	{
		y := g.rowPadding
		for i := 0; i < numRows; i++ {
			rowPositions[i] = y
			y += rowHeights[i] + g.rowPadding
		}
	}

	// Reposition and resize all children.
	for i, child := range g.children {
		position := g.cellPositions[i]
		span := g.cellSpans[i]

		x := columnPositions[position[0]]
		y := rowPositions[position[1]]

		var w, h int
		for j := 0; j < span[0]; j++ {
			if j > 0 {
				w += g.columnPadding
			}
			w += columnWidths[position[0]+j]
		}
		for j := 0; j < span[1]; j++ {
			if j > 0 {
				h += g.rowPadding
			}
			h += rowHeights[position[1]+j]
		}

		child.SetRect(image.Rect(x, y, x+w, y+h))
	}
}
