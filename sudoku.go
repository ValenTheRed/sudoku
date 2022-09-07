package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SudokuCell struct {
	value    byte
	readonly bool
}

// NewSudokuCell returns a new, modifiable SudokuCell.
func NewSudokuCell(val byte) *SudokuCell {
	return &SudokuCell{
		value: val,
	}
}

// NewSudokuCell returns a new, readonly SudokuCell.
func NewReadonlySudokuCell(val byte) *SudokuCell {
	return NewSudokuCell(val).SetReadonly(true)
}

func (c *SudokuCell) Value() byte {
	return c.value
}

func (c *SudokuCell) SetValue(v byte) *SudokuCell {
	c.value = v
	if v == '0' {
		c.value = ' '
	}
	return c
}

func (c *SudokuCell) Readonly() bool {
	return c.readonly
}

func (c *SudokuCell) SetReadonly(v bool) *SudokuCell {
	c.readonly = v
	return c
}

func (c SudokuCell) Rune() rune {
	return rune(c.value)
}

type SudokuGrid struct {
	*tview.Box
	selectedRow, selectedColumn int
	contents                    [81]*SudokuCell
}

// NewSudokuGrid returns a new SudokuGrid.
func NewSudokuGrid() *SudokuGrid {
	var contents [81]*SudokuCell
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			contents[9*r+c] = NewSudokuCell('0')
		}
	}
	return &SudokuGrid{
		Box:      tview.NewBox(),
		contents: contents,
	}
}

// SelectCell focuses the cell at row r and column c.
func (g *SudokuGrid) SelectCell(r, c int) *SudokuGrid {
	g.selectedRow, g.selectedColumn = r, c
	return g
}

// GetCell returns the cell at row r and column c.
func (g *SudokuGrid) GetCell(r, c int) *SudokuCell {
	return g.contents[9*r+c]
}

// InputHandler returns the handler for this primitive.
func (g *SudokuGrid) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	forward := func(pos *int) func() {
		return func() {
			*pos = (*pos + 1) % 9
		}
	}
	backward := func(pos *int) func() {
		return func() {
			*pos--
			if *pos < 0 {
				*pos += 9
			}
		}
	}
	down := forward(&g.selectedRow)
	up := backward(&g.selectedRow)
	left := backward(&g.selectedColumn)
	right := forward(&g.selectedColumn)
	return g.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch key := event.Key(); key {
		case tcell.KeyRune:
			switch r := event.Rune(); r {
			case 'j':
				down()
			case 'k':
				up()
			case 'h':
				left()
			case 'l':
				right()
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				cell := g.GetCell(g.selectedRow, g.selectedColumn)
				if !cell.Readonly() {
					cell.SetValue(byte(r))
				}
			}
		case tcell.KeyDown:
			down()
		case tcell.KeyUp:
			up()
		case tcell.KeyLeft:
			left()
		case tcell.KeyRight:
			right()
		}
	})
}

// Draw draws the sudoku grid onto the screen.
func (g *SudokuGrid) Draw(screen tcell.Screen) {
	const (
		hBorderRt    = tview.BoxDrawingsLightRight
		hBorder      = tview.BoxDrawingsLightHorizontal
		hBorderLt    = tview.BoxDrawingsLightLeft
		hBorderHeavy = tview.BoxDrawingsHeavyHorizontal

		vBorder      = tview.BoxDrawingsLightVertical
		vBorderHeavy = tview.BoxDrawingsHeavyVertical

		crossBorder = tview.BoxDrawingsHeavyVerticalAndHorizontal
	)

	g.Box.DrawForSubclass(screen, g)
	X, Y, _, _ := g.Box.GetInnerRect()

	heavyBorderStyle := tcell.StyleDefault.Foreground(Accent).Background(Theme.background)
	lightBorderStyle := tcell.StyleDefault.Foreground(Theme.helpKey).Background(Theme.background)
	cellStyle := tcell.StyleDefault.Foreground(Theme.foreground).Background(Theme.background)
	readonlyStyle := tcell.StyleDefault.Foreground(Theme.foreground).Background(Accent)

	// helper function to draw i-th cell at row y and column x.
	drawCell := func(c *SudokuCell, style func(tcell.Style) tcell.Style, x, y int) {
		if c.Readonly() {
			screen.SetContent(X+x, Y+y, c.Rune(), nil, style(readonlyStyle))
		} else {
			screen.SetContent(X+x, Y+y, c.Rune(), nil, style(cellStyle))
		}
	}

	// One row for the numbers, second row for the borders, and we won't
	// draw anything after the last number row.
	for y := 0; y < (9*2)-1; y++ {
		switch {
		// border between subgrid row
		case y == (3*2)-1 || y == (6*2)-1:
			for x := 0; x < (9*4)-1; x++ {
				screen.SetContent(X+x, Y+y, hBorderHeavy, nil, heavyBorderStyle)
			}
			screen.SetContent(X+(4*3)-1, Y+y, crossBorder, nil, heavyBorderStyle)
			screen.SetContent(X+(4*6)-1, Y+y, crossBorder, nil, heavyBorderStyle)
		// border inside subgrid row
		case y%2 != 0:
			runes := []rune{hBorderRt, hBorder, hBorderLt, ' '}
			for x := 0; x < (9*4)-1; x++ {
				screen.SetContent(X+x, Y+y, runes[x%len(runes)], nil, lightBorderStyle)
			}
			screen.SetContent(X+(4*3)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
			screen.SetContent(X+(4*6)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
		// number row
		default:
			for x := 0; x < (9*4)-1; x++ {
				if x%4 == 3 {
					screen.SetContent(X+x, Y+y, vBorder, nil, lightBorderStyle)
					continue
				}
				r, c := y/2, x/4
				cell := g.GetCell(r, c)
				if x%4 == 0 || x%4 == 2 {
					if cell.Readonly() {
						cell = NewReadonlySudokuCell(0)
					} else {
						cell = NewSudokuCell(0)
					}
				}
				if g.selectedRow == r && g.selectedColumn == c {
					drawCell(cell, func(s tcell.Style) tcell.Style {
						return s.Reverse(true)
					}, x, y)
				} else {
					drawCell(cell, func(s tcell.Style) tcell.Style {
						return s
					}, x, y)
				}
			}
			screen.SetContent(X+(4*3)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
			screen.SetContent(X+(4*6)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
		}
	}
}
