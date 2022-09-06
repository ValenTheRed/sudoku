package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ColorScheme struct {
	background, foreground tcell.Color
	helpKey, helpDesc      tcell.Color

	black, white tcell.Color
	// Accents
	cyan   tcell.Color
	purple tcell.Color
	pink   tcell.Color
	red    tcell.Color
	orange tcell.Color
	yellow tcell.Color
	green  tcell.Color
}

// Default colorschemes.
var (
	DarkColorScheme = ColorScheme{
		background: tcell.GetColor("#121014"),
		foreground: tcell.GetColor("#dbdbdb"),
		helpKey:    tcell.GetColor("#626262"),
		helpDesc:   tcell.GetColor("#4a4a4a"),

		black: tcell.GetColor("#000000"),
		white: tcell.GetColor("#ffffff"),
		// Accents
		cyan:   tcell.GetColor("#03bcfc"),
		purple: tcell.GetColor("#a229f9"),
		pink:   tcell.GetColor("#eb17ff"),
		red:    tcell.GetColor("#ff222f"),
		orange: tcell.GetColor("#ff742b"),
		yellow: tcell.GetColor("#fbd900"),
		green:  tcell.GetColor("#b6e78d"),
	}
	LightColorScheme = ColorScheme{
		background: tcell.GetColor("#eeeeee"),
		foreground: tcell.GetColor("#000033"),
		helpKey:    tcell.GetColor("#626262"),
		helpDesc:   tcell.GetColor("#4a4a4a"),

		black: tcell.GetColor("#000000"),
		white: tcell.GetColor("#ffffff"),
		// Accents
		cyan:   tcell.GetColor("#03bcfc"),
		purple: tcell.GetColor("#a229f9"),
		pink:   tcell.GetColor("#f520f6"),
		red:    tcell.GetColor("#ff222f"),
		orange: tcell.GetColor("#ff742b"),
		yellow: tcell.GetColor("#fbd900"),
		green:  tcell.GetColor("#0df50b"),
	}
)

// The theme and accent color to be used within the application.
var (
	Theme  = DarkColorScheme
	Accent = Theme.purple
)

func init() {
	tview.Styles.PrimaryTextColor = Theme.foreground
	tview.Styles.PrimitiveBackgroundColor = Theme.background
}

// viewDefaultColorScheme is used to display the colorscheme as it would
// be used in the application for testing purposes. It returns a
// Primitive to be set as the root of the application.
func viewDefaultColorScheme(theme ColorScheme) tview.Primitive {
	newPrimitive := func(bg tcell.Color) tview.Primitive {
		tv := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetDynamicColors(true).
			SetText(fmt.Sprintf(
				"[:#%06x:]  \n[#%06x::]1 2 3 4 5 6 7 8 9\n[#%06x:-:]1 2 3 4 5 6 7 8 9",
				bg.Hex(),
				theme.white.Hex(),
				theme.foreground.Hex(),
			))
		tv.SetBackgroundColor(theme.background)
		return tv
	}

	helpPrim := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).SetText(
		fmt.Sprintf("[#%06x]q[#%06x] quit • [#%06x]? [#%06x]help",
			theme.helpKey.Hex(),
			theme.helpDesc.Hex(),
			theme.helpKey.Hex(),
			theme.helpDesc.Hex()),
	)
	helpPrim.SetBackgroundColor(theme.background)

	grid := tview.NewGrid().
		SetRows(4, 4, 4, 4, 4, 4, 4, 4, 4).
		AddItem(newPrimitive(theme.cyan), 0, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme.purple), 1, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme.pink), 2, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme.red), 3, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme.orange), 4, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme.yellow), 5, 0, 1, 1, 0, 0, false).
		AddItem(newPrimitive(theme.green), 6, 0, 1, 1, 0, 0, false).
		AddItem(helpPrim, 7, 0, 1, 1, 0, 0, false)

	return grid
}
