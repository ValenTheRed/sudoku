package main

import (
	"github.com/rivo/tview"
)

type SudokuFooter struct {
	*tview.Box
	frame   *SudokuFrame
	buttons [10]*tview.Button
}
