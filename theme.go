package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Accepted color keys -
// background, foreground
// black, white
// cyan, purple, pink, red, orange, yellow, green
type ColorScheme map[string]tcell.Color

// Default colorschemes.
var (
	ColorSchemes = map[string]ColorScheme{
		"dark": {
			"background":      tcell.GetColor("#121014"),
			"uiSurface":       tcell.GetColor("#363636"),
			"darkerUISurface": tcell.GetColor("#4f4f4f"),
			"foreground":      tcell.GetColor("#dbdbdb"),

			"black": tcell.GetColor("#000000"),
			"white": tcell.GetColor("#ffffff"),
			// Accents
			"cyan":   tcell.GetColor("#03a2d9"),
			"purple": tcell.GetColor("#a229f9"),
			"pink":   tcell.GetColor("#eb17ff"),
			"red":    tcell.GetColor("#ff222f"),
			"orange": tcell.GetColor("#f75500"),
			"yellow": tcell.GetColor("#c0a412"),
			"green":  tcell.GetColor("#319253"),
		},
		"light": {
			"background":      tcell.GetColor("#eeeeee"),
			"foreground":      tcell.GetColor("#151515"),
			"uiSurface":       tcell.GetColor("#cccccc"),
			"darkerUISurface": tcell.GetColor("#999999"),

			"black": tcell.GetColor("#000000"),
			"white": tcell.GetColor("#ffffff"),
			// Accents
			"cyan":   tcell.GetColor("#03bcfc"),
			"purple": tcell.GetColor("#b451fa"),
			"pink":   tcell.GetColor("#f63bf7"),
			"red":    tcell.GetColor("#ff4d57"),
			"orange": tcell.GetColor("#ff742b"),
			"yellow": tcell.GetColor("#f0d322"),
			"green":  tcell.GetColor("#43c571"),
		},
	}
)

// The theme and accent color to be used within the application.
var (
	Theme, Accent string
	BlendAccent   tcell.Color
)

func SetTheme(t, accent string) {
	Theme = t
	// Theme  = LightColorScheme
	Accent = accent
	c := ColorSchemes[Theme]
	BlendAccent = colorBlend(c[Accent], c["background"], 20)
}

// InitModalStyle initialsies Modal m with a custom, default style.
func InitModalStyle(m *Modal) *Modal {
	m.SetBorderColor(ColorSchemes[Theme][Accent])
	m.SetBackgroundColor(ColorSchemes[Theme]["background"])
	m.SetButtonBackgroundColor(ColorSchemes[Theme]["uiSurface"])
	m.SetButtonTextColor(ColorSchemes[Theme]["foreground"])
	m.SetTextColor(ColorSchemes[Theme]["foreground"])
	return m
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
				theme["white"].Hex(),
				theme["foreground"].Hex(),
			))
		tv.SetBackgroundColor(theme["background"])
		return tv
	}

	helpPrim := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).SetText(
		fmt.Sprintf(
			"[#%06x]q[#%06x] quit ??? [#%06x]? [#%06x]help",
			theme["darkerUISurface"].Hex(),
			theme["uiSurface"].Hex(),
			theme["darkerUISurface"].Hex(),
			theme["uiSurface"].Hex(),
		),
	)
	helpPrim.SetBackgroundColor(theme["background"])

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

// colorBlend blends color src with sink. alpha must be in the closed
// interval [0, 100]. A value of 0 for alpha results in sink and a value
// of 100 results in src. Formula:
//	(b*alpha + a*(100 - alpha)) / 2
func colorBlend(src, sink tcell.Color, alpha int32) tcell.Color {
	r, g, b := src.RGB()
	srcRGB := []int32{r, g, b}
	r, g, b = sink.RGB()
	sinkRGB := []int32{r, g, b}
	blendChannel := func(i int) int32 {
		return int32((alpha*srcRGB[i] + (100-alpha)*sinkRGB[i]) / 100)
	}
	return tcell.NewRGBColor(
		blendChannel(0),
		blendChannel(1),
		blendChannel(2),
	)
}
