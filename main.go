package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/rivo/tview"
)

var (
	// undopath stores the path to the undofile
	undopath string

	// continueFlag set to true will restore the puzzle from the
	// previous session.
	continueFlag bool
)

func init() {
	flag.BoolVar(&continueFlag, "continue", false, "restore previous sesssions puzzle")
	flag.BoolVar(&continueFlag, "c", false, "restore previous sesssions puzzle")

	// set undopath to the path of the undo file

	localshare, exists := os.LookupEnv("XDG_DATA_HOME")
	if !exists {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalln(err)
		}
		localshare = path.Join(home, `.local/share/sudoku`)
	} else {
		localshare = path.Join(localshare, `sudoku`)
	}

	undopath = path.Join(localshare, `undo`)

	if err := os.MkdirAll(localshare, 0750); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()

	SetTheme("dark", "purple")

	app := tview.NewApplication().EnableMouse(true)

	frame := NewSudokuFrame()
	frame.timer.SetChangedFunc(func() {
		app.Draw()
	})
	if continueFlag {
		func() {
			file, err := os.OpenFile(undopath, os.O_CREATE|os.O_RDONLY, 0750)
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
			frame.grid.ReadUndoHistoryFromFile(file)
		}()
	}

	sidepane := NewSidepane()

	// Restart this game
	solveModal := NewModal()
	InitModalStyle(solveModal)
	solveModal.SetText("Do you want to solve this game?")
	solveModal.AddButtons([]string{"Cancel", "Yes"})

	// Restart this game
	restartModal := NewModal()
	InitModalStyle(restartModal)
	restartModal.SetText("Do you want to restart this game?")
	restartModal.AddButtons([]string{"Cancel", "Yes"})

	// Validate game
	validateModal := NewModal()
	InitModalStyle(validateModal)
	validateModal.SetText("Do you want to validate the puzzle?")
	validateModal.AddButtons([]string{"Cancel", "Yes"})
	validateModal.SetFocus(1)

	accentModal := NewModal()
	InitModalStyle(accentModal)
	accentModal.SetText("Choose color")
	accentModal.AddButtons([]string{
		"Cyan", "Purple", "Pink", "Red", "Orange", "Yellow", "Green",
	})

	switchAppTheme := func(t, accent string) {
		go func() {
			SetTheme(t, accent)
			InitSidepaneStyle(sidepane)
			InitModalStyle(solveModal)
			InitModalStyle(restartModal)
			InitModalStyle(validateModal)
			InitModalStyle(accentModal)
			app.Draw()
		}()
	}

	// Theme changer
	sidepane.GetButton(4).SetSelectedFunc(func() {
		t := "dark"
		if Theme == "dark" {
			t = "light"
		}
		switchAppTheme(t, Accent)
	})

	grid := tview.NewGrid()
	grid.SetRows(0).SetColumns(-1, -3).
		AddItem(sidepane, 0, 0, 1, 1, 0, 0, false).
		AddItem(frame, 0, 1, 1, 1, 0, 0, true)

	pages := tview.NewPages()
	pages.AddPage("grid", grid, true, true)
	pages.AddPage("restart", restartModal, true, false)
	pages.AddPage("solve", solveModal, true, false)
	pages.AddPage("validate", validateModal, true, false)
	pages.AddPage("accent", accentModal, true, false)
	sidepane.GetButton(0).SetSelectedFunc(func() {
		frame.grid.Undo()
	})
	sidepane.GetButton(5).SetSelectedFunc(func() {
		pages.ShowPage("accent")
	})
	accentModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Cyan", "Purple", "Pink", "Red", "Orange", "Yellow", "Green":
			switchAppTheme(Theme, strings.ToLower(buttonLabel))
		}
		pages.SwitchToPage("grid")
	})
	sidepane.GetButton(1).SetSelectedFunc(func() {
		pages.ShowPage("validate")
	})
	validateModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			// TODO: validating logic
		}
		pages.SwitchToPage("grid")
		validateModal.SetFocus(1)
	})
	sidepane.GetButton(2).SetSelectedFunc(func() {
		pages.ShowPage("solve")
	})
	solveModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			// TODO: auto-solve logic
		}
		pages.SwitchToPage("grid")
		solveModal.SetFocus(0)
	})
	sidepane.GetButton(3).SetSelectedFunc(func() {
		pages.ShowPage("restart")
	})
	restartModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Yes" {
			// TODO: restarting logic
		}
		pages.SwitchToPage("grid")
		restartModal.SetFocus(0)
	})

	frame.timer.Start()
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		log.Println(err)
	}

	file, err := os.OpenFile(
		undopath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0750,
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	frame.grid.FlushUndoHistoryToFile(file)
}
