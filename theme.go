package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var DefaultColorScheme = struct {
	Dark, Light map[string]string
}{
	Dark: map[string]string{
		"background": "#121014",
		"foreground": "#dbdbdb",
		"helpKey":    "#626262",
		"helpDesc":   "#4a4a4a",

		"black": "#000000",
		"white": "#ffffff",
		// Accents
		"cyan":   "#03bcfc",
		"purple": "#a229f9",
		"pink":   "#eb17ff",
		"red":    "#ff222f",
		"orange": "#ff742b",
		"yellow": "#fbd900",
		"green":  "#b6e78d",
	},
	Light: map[string]string{
		"background": "#eeeeee",
		"foreground": "#000033",
		"helpKey":    "#626262",
		"helpDesc":   "#4a4a4a",

		"black": "#000000",
		"white": "#ffffff",
		// Accents
		"cyan":   "#03bcfc",
		"purple": "#a229f9",
		"pink":   "#f520f6",
		"red":    "#ff222f",
		"orange": "#ff742b",
		"yellow": "#fbd900",
		"green":  "#0df50b",
	},
}

// The theme and accent color to be used within the application.
var (
	Theme = DefaultColorScheme.Dark
	Accent = Theme["purple"]
)

// viewDefaultColorScheme is used to display the colorscheme as it would
// be used in the application for testing purposes. It returns a
// Primitive to be set as the root of the application.
func viewDefaultColorScheme(theme map[string]string) tview.Primitive {
	newPrimitive := func(bg string) tview.Primitive {
		tv := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetDynamicColors(true).
			SetText(fmt.Sprintf(
				"[:%s:]  \n[%s::]1 2 3 4 5 6 7 8 9\n[%s:-:]1 2 3 4 5 6 7 8 9",
				bg,
				theme["white"],
				theme["foreground"],
			))
		tv.SetBackgroundColor(tcell.GetColor(theme["background"]))
		return tv
	}

	helpPrim := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).SetText(
		fmt.Sprintf("[%s]q[%s] quit • [%s]? [%s]help",
			theme["helpKey"],
			theme["helpDesc"],
			theme["helpKey"],
			theme["helpDesc"]),
	)
	helpPrim.SetBackgroundColor(tcell.GetColor(theme["background"]))

	grid := tview.NewGrid().
		SetRows(4, 4, 4, 4, 4, 4, 4, 4, 4).
		AddItem(newPrimitive(theme["cyan"]), 0, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme["purple"]), 1, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme["pink"]), 2, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme["red"]), 3, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme["orange"]), 4, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme["yellow"]), 5, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme["green"]), 6, 0, 1, 1, 0, 0, false).
		AddItem(helpPrim, 7, 0, 1, 1, 0, 0, false)

	return grid
}
