package main

import (
	"bufio"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SudokuCell struct {
	value    byte
	readonly bool
}

// NewSudokuCell returns a new, modifiable SudokuCell.
func NewSudokuCell(digit int) *SudokuCell {
	return (&SudokuCell{}).SetValue(digit)
}

// NewSudokuCell returns a new, readonly SudokuCell.
func NewReadonlySudokuCell(digit int) *SudokuCell {
	return NewSudokuCell(digit).SetReadonly(true)
}

// Value returns the digit at c.
func (c *SudokuCell) Value() int {
	if c.value == ' ' {
		return 0
	}
	return int(c.value - '0')
}

func (c *SudokuCell) SetValue(digit int) *SudokuCell {
	if digit == 0 {
		c.value = ' '
	} else {
		c.value = byte(digit) + '0'
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

func (c *SudokuCell) Rune() rune {
	return rune(c.value)
}

func (c *SudokuCell) IsEmpty() bool {
	return c.value == ' '
}

type undoItem struct {
	row, col, digit byte
}

type SudokuGrid struct {
	*tview.Box
	selectedRow, selectedColumn int
	contents                    [81]*SudokuCell
	undoHistory                 []undoItem
}

// NewSudokuGrid returns a new SudokuGrid.
func NewSudokuGrid() *SudokuGrid {
	var contents [81]*SudokuCell
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			contents[9*r+c] = NewSudokuCell(0)
		}
	}
	return &SudokuGrid{
		Box:      tview.NewBox(),
		contents: contents,
	}
}

// ClearCells clears all non-readonly cells and the undo history.
func (g *SudokuGrid) ClearCells() *SudokuGrid {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if !g.GetCell(r, c).Readonly() {
				g.SetCellWithUndo(r, c, 0)
			}
		}
	}
	g.undoHistory = nil
	return g
}

// SelectCell focuses the cell at row r and column c.
func (g *SudokuGrid) SelectCell(r, c int) *SudokuGrid {
	g.selectedRow, g.selectedColumn = r, c
	return g
}

// SelectedCell row r and column c of the selected cell.
func (g *SudokuGrid) SelectedCell() (int, int) {
	return g.selectedRow, g.selectedColumn
}

// GetCell returns the cell at row r and column c.
func (g *SudokuGrid) GetCell(r, c int) *SudokuCell {
	return g.contents[9*r+c]
}

// SetCellWithUndo sets the value of cell at row r and column c with the
// value digit. It doesn't store the previous value of the cell in it's
// undo history.
func (g *SudokuGrid) SetCellWithoutUndo(r, c, digit int) *SudokuGrid {
	g.GetCell(r, c).SetValue(digit)
	return g
}

// SetCellWithUndo sets the value of cell at row r and column c with the
// value digit. It also stores the previous value of the cell in it's
// undo history.
func (g *SudokuGrid) SetCellWithUndo(r, c, digit int) *SudokuGrid {
	cell := g.GetCell(r, c)
	if digit == cell.Value() {
		return g
	}
	g.undoHistory = append(g.undoHistory, undoItem{
		byte(r), byte(c), byte(cell.Value()),
	})
	cell.SetValue(digit)
	return g
}

// Undo undos the last move.
func (g *SudokuGrid) Undo() *SudokuGrid {
	if len(g.undoHistory) > 0 {
		item := g.undoHistory[len(g.undoHistory)-1]
		g.undoHistory = g.undoHistory[:len(g.undoHistory)-1]
		g.SetCellWithoutUndo(
			int(item.row),
			int(item.col),
			int(item.digit),
		)
	}
	return g
}

// FlushUndoHistoryToFile writes the entire undo history to file and
// resets the history.
// NOTE: empty cell is denoted by '.'.
func (g *SudokuGrid) FlushUndoHistoryToFile(file *os.File) *SudokuGrid {
	for _, item := range g.undoHistory {
		a := item.row + '0'
		b := item.col + '0'
		c := byte('.')
		if d := item.digit; d != 0 {
			c = d + '0'
		}
		file.Write([]byte{a, ' ', b, ' ', c, '\n'})
	}
	g.undoHistory = nil
	return g
}

// ReadUndoHistoryFromFile reads file and appends it's undo history to
// g.
func (g *SudokuGrid) ReadUndoHistoryFromFile(file *os.File) *SudokuGrid {
	for s := bufio.NewScanner(file); s.Scan(); {
		line := s.Bytes()
		if len(line) != 5 {
			log.Fatalf("ReadUndoHistoryFromFile: parsing undo \"%s\": line length must be 5\n", line)
		}
		a, b, c := line[0], line[2], line[4]
		if a < '0' || a > '9' {
			log.Fatalf("ReadUndoHistoryFromFile: parsing undo: first character is %c, must be in the set [1-9]", a)
		}
		a -= '0'
		if b < '0' || b > '9' {
			log.Fatalf("ReadUndoHistoryFromFile: parsing undo: second character is %c, must be in the set [1-9]", b)
		}
		b -= '0'
		if c == '.' {
			c = 0
		} else if c >= '0' && c <= '9' {
			c = c - '0'
		} else {
			log.Fatalf("ReadUndoHistoryFromFile: parsing undo: third character is %c, must be in the set [.1-9]", b)
		}
		g.undoHistory = append(g.undoHistory, undoItem{a, b, c})
	}
	return g
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
					g.SetCellWithUndo(g.selectedRow, g.selectedColumn, int(r-'0'))
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

const (
	// I assume that a cell in the grid is composed of the characters
	// for its contents then one border character. So,
	// cell length = len(cell contents + border character)

	// cell vertical length = len('<number>' + '-') = 2
	SudokuGridRowHeight = 2
	// cell horizontal length = len(' ' + '<number>' ' ' + '|') = 4
	SudokuGridColumnWidth = 4
)

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

	g.Box.SetBackgroundColor(ColorSchemes[Theme]["background"])
	g.Box.DrawForSubclass(screen, g)
	X, Y := g.centerCoordinates()

	heavyBorderStyle := tcell.StyleDefault.Foreground(ColorSchemes[Theme][Accent]).Background(ColorSchemes[Theme]["background"])
	lightBorderStyle := tcell.StyleDefault.Foreground(ColorSchemes[Theme]["uiSurface"]).Background(ColorSchemes[Theme]["background"])
	cellStyle := tcell.StyleDefault.Foreground(ColorSchemes[Theme]["foreground"]).Background(ColorSchemes[Theme]["background"])
	readonlyStyle := tcell.StyleDefault.Foreground(ColorSchemes[Theme]["foreground"]).Background(ColorSchemes[Theme][Accent])

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
	for y := 0; y < (9*SudokuGridRowHeight)-1; y++ {
		switch {
		// border between subgrid row
		case y == (3*SudokuGridRowHeight)-1 || y == (6*SudokuGridRowHeight)-1:
			for x := 0; x < (9*SudokuGridColumnWidth)-1; x++ {
				screen.SetContent(X+x, Y+y, hBorderHeavy, nil, heavyBorderStyle)
			}
			screen.SetContent(X+(SudokuGridColumnWidth*3)-1, Y+y, crossBorder, nil, heavyBorderStyle)
			screen.SetContent(X+(SudokuGridColumnWidth*6)-1, Y+y, crossBorder, nil, heavyBorderStyle)
		// border inside subgrid row
		case y%SudokuGridRowHeight != 0:
			runes := []rune{hBorderRt, hBorder, hBorderLt, ' '}
			for x := 0; x < (9*SudokuGridColumnWidth)-1; x++ {
				screen.SetContent(X+x, Y+y, runes[x%len(runes)], nil, lightBorderStyle)
			}
			screen.SetContent(X+(SudokuGridColumnWidth*3)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
			screen.SetContent(X+(SudokuGridColumnWidth*6)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
		// number row
		default:
			for x := 0; x < (9*SudokuGridColumnWidth)-1; x++ {
				if x%SudokuGridColumnWidth == 3 {
					screen.SetContent(X+x, Y+y, vBorder, nil, lightBorderStyle)
					continue
				}
				r, c := y/SudokuGridRowHeight, x/SudokuGridColumnWidth
				cell := g.GetCell(r, c)
				if x%SudokuGridColumnWidth == 0 || x%SudokuGridColumnWidth == 2 {
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
			screen.SetContent(X+(SudokuGridColumnWidth*3)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
			screen.SetContent(X+(SudokuGridColumnWidth*6)-1, Y+y, vBorderHeavy, nil, heavyBorderStyle)
		}
	}
}

// centerCoordinates calculates and returns the (X, Y) coordinates of
// the point within the SudokuGrid bounding box from where, if drawn,
// the SudokuGrid looks centered.
func (g *SudokuGrid) centerCoordinates() (X, Y int) {
	X, Y, width, height := g.Box.GetInnerRect()
	if width := width - (9 * SudokuGridColumnWidth) - 1; width > 0 {
		X += width / 2
	}
	if height := height - (9 * SudokuGridRowHeight) - 1; height > 0 {
		Y += height / 2
	}
	return X, Y
}

func (g *SudokuGrid) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return g.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		if !g.InRect(x, y) {
			return
		}
		// From my investigation, both, MouseLeftDown and
		// MouseLeftClick, are sent whenever you left click.
		// MouseLeftDown is sent first and then MouseLeftClick.
		//
		// The follwing switch structure has been referenced from
		// tview.TextView.
		switch action {
		case tview.MouseLeftDown:
			setFocus(g)
			consumed = true
		case tview.MouseLeftClick:
			r, c := g.cellAtCoordinate(event.Position())
			if r != -1 {
				g.SelectCell(r, c)
			}
			consumed = true
		}
		return consumed, nil
	})
}

// cellAtCoordinate returns the row and column of the cell enclosing
// the point at (x, y). Returns (-1, -1) if point outside g's bounding
// box, or point resides on a border character.
func (g *SudokuGrid) cellAtCoordinate(x, y int) (r, c int) {
	if !g.InRect(x, y) {
		return -1, -1
	}
	X, Y := g.centerCoordinates()
	x, y = x-X, y-Y
	width, height := 9*SudokuGridColumnWidth-1, 9*SudokuGridRowHeight-1
	if x < 0 || y < 0 || x > width || y > height {
		return -1, -1
	}
	if y%SudokuGridRowHeight == SudokuGridRowHeight-1 ||
		x%SudokuGridColumnWidth == SudokuGridColumnWidth-1 {
		return -1, -1
	}
	return y / SudokuGridRowHeight, x / SudokuGridColumnWidth
}
