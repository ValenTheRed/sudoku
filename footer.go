package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type SudokuFooter struct {
	*tview.Box
	frame   *SudokuFrame
	buttons [10]*tview.Button
}

func NewSudokuFooter(frame *SudokuFrame) *SudokuFooter {
	f := &SudokuFooter{
		Box:   tview.NewBox(),
		frame: frame,
	}

	const cross = '✗'

	f.Box.SetBorderPadding(1, 0, 0, 0)

	newBtn := func(char rune) *tview.Button {
		b := tview.NewButton(fmt.Sprintf(" %c ", char))
		b.SetLabelColorActivated(Accent)
		b.SetBackgroundColor(Theme.helpDesc)
		b.SetSelectedFunc(func() {
			char := []rune(b.GetLabel())[1]
			if char == cross {
				char = '0'
			}
			g := f.frame.grid
			if cell := g.GetCell(g.SelectedCell()); !cell.Readonly() {
				cell.SetValue(int(char - '0'))
			}
		})
		return b
	}

	for i := 0; i < 9; i++ {
		f.buttons[i] = newBtn(rune(i) + '1')
	}
	f.buttons[9] = newBtn(cross)

	return f
}
