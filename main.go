package main

import (
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication().EnableMouse(true)

	grid := GeneratePuzzle()
	root := tview.NewGrid().
		SetRows(0).SetColumns(0).AddItem(grid, 0, 0, 1, 1, 0, 0, true)

	if err := app.SetRoot(root, true).SetFocus(root).Run(); err != nil {
		panic(err)
	}
}
