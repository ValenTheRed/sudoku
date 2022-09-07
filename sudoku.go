package main

import (
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
