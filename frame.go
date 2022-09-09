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
