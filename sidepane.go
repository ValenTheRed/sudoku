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

	// Optional func that will be triggered when button is pressed.
	selected func()
}

func newButton(icon rune, text string) *button {
	return &button{
		Box: tview.NewBox().
			SetBorder(true),
		icon: icon,
		text: text,
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

// SetSelectedFunc sets f as the optional handler to fire when b is
// selected.
func (b *button) SetSelectedFunc(f func()) *button {
	b.selected = f
	return b
}

func (b *button) Draw(screen tcell.Screen) {
	defaultStyle := tcell.StyleDefault.
		Background(colorBlend(Accent, Theme.background, 20)).
		Foreground(Theme.foreground)
	selectedStyle := tcell.StyleDefault.
		Background(Accent).
		Foreground(Theme.foreground).
		Attributes(tcell.AttrUnderline)

	style := defaultStyle
	if b.HasFocus() {
		style = selectedStyle
	}
	_, bg, _ := style.Decompose()
	b.SetBorderColor(Accent)
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
	Buttons [6]*button
}

func NewSidepane() *Sidepane {
	s := &Sidepane{
		Box: tview.NewBox(),
	}

	s.SetBorderPadding(1, 1, 1, 1)

	for i, item := range [6]struct {
		icon  rune
		label string
	}{
		{'', "Undo"},
		{'', "Validate"},
		{'', "Solve"},
		{'', "Restart"},
		{'', "Switch theme"},
		{'', "Change Accent"},
	} {
		s.Buttons[i] = newButton(item.icon, item.label)
	}

	return s
}

func (s *Sidepane) Draw(screen tcell.Screen) {
	s.SetBackgroundColor(
		colorBlend(Accent, Theme.background, 20),
	)
	s.DrawForSubclass(screen, s)
	x, y, width, _ := s.GetInnerRect()

	const (
		buttonHeight  = 3
		buttonPadding = 0
	)

	for _, button := range s.Buttons {
		button.SetRect(x, y, width, buttonHeight)
		button.Draw(screen)
		y += buttonHeight + buttonPadding
	}
}

func (s *Sidepane) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return s.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !s.InRect(event.Position()) {
			return
		}
		// Pass mouse events along to the first child item that takes it.
		for _, button := range s.Buttons {
			consumed, capture = button.MouseHandler()(action, event, setFocus)
			if consumed {
				return consumed, capture
			}
		}
		return
	})
}
