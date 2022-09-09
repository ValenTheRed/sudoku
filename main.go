package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	app.EnableMouse(true)

	frame := NewSudokuFrame()

	root := tview.NewGrid().
		SetRows(0).SetColumns(0).
		AddItem(frame, 0, 0, 1, 1, 0, 0, true)

	frame.timer.SetChangedFunc(func() {
		app.Draw()
	})

	frame.timer.Start()
	if err := app.SetRoot(root, true).SetFocus(root).Run(); err != nil {
		panic(err)
	}
}
