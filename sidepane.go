package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type button struct {
	*tview.Box

	// icon will appear on the label, just before the main text.
	// NOTE: icon must be a 1 width character.
	icon        rune
	iconDisable bool

	// text is the text that will appear on the label alongside the
	// icon.
	text string

	defaultStyle  tcell.Style
	selectedStyle tcell.Style

	// Optional func that will be triggered when button is pressed.
	selected func()
}

func newButton(icon rune, text string) *button {
	return &button{
		Box: tview.NewBox().
			SetBorder(true).
			SetBorderColor(Accent),
		icon:  icon,
		text: text,
		defaultStyle: tcell.StyleDefault.
			Background(colorBlend(Accent, tview.Styles.PrimitiveBackgroundColor, 20)).
			Foreground(tview.Styles.PrimaryTextColor),
		selectedStyle: tcell.StyleDefault.
			Background(Accent).
			Foreground(tview.Styles.PrimaryTextColor).
			Attributes(tcell.AttrUnderline),
	}
}

// SetIconDisable disables icons if v is true.
func (b *button) SetIconDisable(v bool) *button {
	b.iconDisable = v
	return b
}

type Sidepane struct {
	*tview.Box
}

func NewSidepane() *Sidepane {
	s := &Sidepane{
		Box: tview.NewBox(),
	}
	return s
}
