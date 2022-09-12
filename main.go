package main

import (
	"github.com/rivo/tview"
)

func main() {
	SetTheme(DarkColorScheme)

	app := tview.NewApplication().EnableMouse(true)

	frame := NewSudokuFrame()
	frame.timer.SetChangedFunc(func() {
		app.Draw()
	})

	sidepane := NewSidepane()

	// Restart this game
	restartModal := NewModal()
	InitModalStyle(restartModal)
	restartModal.SetText("Do you want to restart this game?")
	restartModal.AddButtons([]string{"Cancel", "Yes"})

	// Theme changer
	sidepane.Buttons[4].SetSelectedFunc(func() {
		go func() {
			t := DarkColorScheme
			if Theme == DarkColorScheme {
				t = LightColorScheme
			}
			SetTheme(t)
			InitModalStyle(restartModal)
			app.Draw()
		}()
	})

	grid := tview.NewGrid()
	grid.SetRows(0).SetColumns(-1, -3).
		AddItem(sidepane, 0, 0, 1, 1, 0, 0, false).
		AddItem(frame, 0, 1, 1, 1, 0, 0, true)

	pages := tview.NewPages()
	pages.AddPage("grid", grid, true, true)
	pages.AddPage("restart", restartModal, true, false)
	sidepane.Buttons[3].SetSelectedFunc(func() {
		pages.ShowPage("restart")
	})
	restartModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			// TODO: restarting logic
		}
		pages.SwitchToPage("grid")
	})

	frame.timer.Start()
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
