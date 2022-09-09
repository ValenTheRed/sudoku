package main

import (
	"github.com/rivo/tview"
)

type SudokuFrame struct {
	*tview.Grid
	difficulty *SudokuHeader
	timer      *Timer
	grid       *SudokuGrid
}

func NewSudokuFrame() *SudokuFrame {
	f := &SudokuFrame{
		Grid: tview.NewGrid(),
		grid: GeneratePuzzle(),
	}
	f.difficulty = NewSudokuHeader(f)
	f.timer = NewTimer(f)
	f.difficulty.SetText("Difficulty")
	f.
		SetRows(2, 0).SetColumns(0, 0).
		AddItem(f.timer, 0, 1, 1, 1, 0, 0, false).
		AddItem(f.difficulty, 0, 0, 1, 1, 0, 0, false).
		AddItem(f.grid, 1, 0, 1, 2, 0, 0, true)
	return f
}
