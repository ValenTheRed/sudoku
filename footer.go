package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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
		b.SetSelectedFunc(func() {
			char := []rune(b.GetLabel())[1]
			if char == cross {
				char = '0'
			}
			g := f.frame.grid
			// NOTE:
			// I expected that after setting a different value to the
			// cell, I would need to manually invoke some some Draw() to
			// get the grid to update, but nope, grid updates
			// automatically and instaneously, and I've no idea how.
			// It updates before the timer updates, so timer's
			// SetChangedFunc() handler can't be the reason.
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

// Draw draws f horizontally centered with one cell width gap at the top.
func (f *SudokuFooter) Draw(screen tcell.Screen) {
	f.DrawForSubclass(screen, f)

	// assumption: no borders around numbers.
	const (
		columns    = 5
		rows       = 2
		cellWidth  = 4
		cellHeight = 2
	)

	X, _ := f.frame.grid.centerCoordinates()
	x, y, _, _ := f.GetInnerRect()

	width := columns*cellWidth - 1
	sudokuWidth := 9*SudokuGridColumnWidth - 1
	x = X + (sudokuWidth-width)/2

	for i, button := range f.buttons {
		if i == 5 {
			y += cellHeight
		}
		// I refrenced tview.Grid.Draw() and tview.Flex.Draw() for
		// writing this Draw() function.

		// This informs the Primitive of it's position so that it
		// can draw from there.
		button.SetRect(x+(cellWidth*(i%columns)), y, cellWidth-1, cellHeight-1)
		// NOTE: I don't really know why we check for focus and defer
		// draw, but I assume tview has a good reason for doing this.
		if button.HasFocus() {
			defer button.Draw(screen)
		} else {
			button.Draw(screen)
		}
	}
}

func (f *SudokuFooter) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return f.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !f.InRect(event.Position()) {
			return
		}
		// Pass mouse events along to the first child item that takes it.
		for _, button := range f.buttons {
			consumed, capture = button.MouseHandler()(action, event, setFocus)
			if consumed {
				return consumed, capture
			}
		}
		return
	})
}
