package main

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
