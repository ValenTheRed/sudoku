package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
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

// GetLabel returns the label text of b.
func (b *button) GetLabel() string {
	if b.iconDisable {
		return b.text
	}
	return string([]rune{b.icon, ' '}) + b.text
}

func (b *button) Draw(screen tcell.Screen) {
	style := b.defaultStyle
	if b.HasFocus() {
		style = b.selectedStyle
	}
	_, bg, _ := style.Decompose()
	b.SetBackgroundColor(bg)
	b.DrawForSubclass(screen, b)

	label := b.GetLabel()
	x, y, width, height := b.GetInnerRect()
	y += height / 2
	x += (width - runewidth.StringWidth(label)) / 2

	for _, r := range label {
		screen.SetContent(x, y, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
}

func (b *button) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch key := event.Key(); key {
		case tcell.KeyEnter:
			// Selected
			if b.selected != nil {
				b.selected()
			}
		}
	})
}

func (b *button) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return b.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !b.InRect(event.Position()) {
			return
		}
		switch action {
		case tview.MouseLeftDown:
			setFocus(b)
			consumed = true
		case tview.MouseLeftClick:
			// Selected Func
			if b.selected != nil {
				b.selected()
			}
			consumed = true
		}
		return
	})
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
