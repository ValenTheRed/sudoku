package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

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

// NewSudokuFrameFromFile returns an initialised SudokuFrame from
// savefile, with undofile used to restore the undo history.
func NewSudokuFrameFromFile(savefile, undofile *os.File) *SudokuFrame {
	f := &SudokuFrame{
		Grid: tview.NewGrid(),
		grid: NewSudokuGrid(),
	}
	f.difficulty = NewSudokuHeader(f)
	f.timer = NewTimer(f)
	f.numberPad = NewSudokuFooter(f)

	scan := bufio.NewScanner(savefile)

	// puzzle
	if !scan.Scan() {
		log.Fatalln("NewSudokuFrameFromFile: parsing savefile: no puzzle text")
	}
	bytes := scan.Bytes()
	if l := len(bytes); l < 81 || l > 2*81 {
		log.Fatalf("NewSudokuFrameFromFile: parsing puzzle \"%s\": invalid length: have %d, want 81 <= length <= 162", bytes, l)
	}
	for i, j := 0, 0; i < len(bytes); i++ {
		if b := bytes[i]; b != '_' && b != '.' && !(b >= '0' && b <= '9') {
			log.Fatalln("NewSudokuFrameFromFile: parsing puzzle: character must be in the set [_.1-9]")
		}
		r, c := j/9, j%9
		if bytes[i] == '_' {
			f.grid.GetCell(r, c).SetReadonly(true)
			continue
		}
		v := 0
		if bytes[i] != '.' {
			v = int(bytes[i] - '0')
		}
		f.grid.SetCellWithoutUndo(r, c, v)
		j++
	}

	// time
	if !scan.Scan() {
		log.Fatalln("NewSudokuFrameFromFile: parsing savefile: no elapsed time")
	}

	if t, err := strconv.Atoi(scan.Text()); err != nil {
		log.Fatalln("NewSudokuFrameFromFile: parsing elapsed time:", err)
	} else {
		f.timer.elapsed = second(t)
	}

	// difficulty
	if !scan.Scan() {
		log.Fatalln("NewSudokuFrameFromFile: parsing savefile: no difficulty text")
	}
	switch t := scan.Text(); t {
	case "Easy", "Medium", "Hard":
		f.difficulty.SetText(scan.Text())
	default:
		log.Fatalln("NewSudokuFrameFromFile: parsing difficulty: difficulty must be either one of: Easy, Medium, Hard")
	}

	f.grid.ReadUndoHistoryFromFile(undofile)

	f.SetRows(0, 9*SudokuGridRowHeight-1, 0).SetColumns(0, 0)
	f.
		AddItem(f.timer, 0, 1, 1, 1, 0, 0, false).
		AddItem(f.difficulty, 0, 0, 1, 1, 0, 0, false)
	f.AddItem(f.grid, 1, 0, 1, 2, 0, 0, true)
	f.AddItem(f.numberPad, 2, 0, 1, 2, 0, 0, false)
	return f
}

// SavePuzzleToFile saves puzzle, puzzle time, and puzzle difficulty to
// file, in that order. It also saves the undo history.
// NOTE: It uses '.' to denote empty cell.
// NOTE: It appends '_' in front of readonly cells
func (f *SudokuFrame) SavePuzzleToFile(savefile, undofile *os.File) {
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
		savefile.Write(s)
	}
	savefile.Write([]byte{'\n'})
	fmt.Fprintln(savefile, int(f.timer.elapsed))
	fmt.Fprintln(savefile, f.difficulty.GetText(true))

	g.FlushUndoHistoryToFile(undofile)
}
