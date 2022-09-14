package main

import (
	"fmt"
	"os"

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

// SavePuzzleToFile saves puzzle, puzzle time, and puzzle difficulty to
// file, in that order.
// NOTE: It uses '.' to denote empty cell.
// NOTE: It appends '_' in front of readonly cells
func (f *SudokuFrame) SavePuzzleToFile(file *os.File) {
	g := f.grid
	for r := 0; r < 9; r++ {
		var s []byte
		for c := 0; c < 9; c++ {
			cell := g.GetCell(r, c)
			var v byte
			if digit := cell.Value(); digit == 0 {
				v = '.'
			} else {
				v = byte(digit) + '0'
			}
			if cell.Readonly() {
				s = append(s, '_', v)
			} else {
				s = append(s, v)
			}
		}
		file.Write(s)
	}
	file.Write([]byte{'\n'})
	fmt.Fprintln(file, int(f.timer.elapsed))
	fmt.Fprintln(file, f.difficulty.GetText(true))
}
