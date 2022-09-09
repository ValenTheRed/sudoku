package main

import (
	"github.com/rivo/tview"
)

type SudokuFrame struct {
	*tview.Grid
	difficulty *SudokuHeader
	timer      *Timer
	grid       *SudokuGrid
	numberPad  *SudokuFooter
}

func NewSudokuFrame() *SudokuFrame {
	f := &SudokuFrame{
		Grid: tview.NewGrid(),
		grid: GeneratePuzzle(),
	}
	f.difficulty = NewSudokuHeader(f)
	f.timer = NewTimer(f)
	f.numberPad = NewSudokuFooter(f)
	f.difficulty.SetText("Difficulty")

	f.SetRows(0, 9*SudokuGridRowHeight-1, 0).SetColumns(0, 0)
	f.
		AddItem(f.timer, 0, 1, 1, 1, 0, 0, false).
		AddItem(f.difficulty, 0, 0, 1, 1, 0, 0, false)
	f.AddItem(f.grid, 1, 0, 1, 2, 0, 0, true)
	f.AddItem(f.numberPad, 2, 0, 1, 2, 0, 0, false)
	return f
}
