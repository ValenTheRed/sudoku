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
		Background(BlendAccent).
		Foreground(ColorSchemes[Theme]["foreground"])
	selectedStyle := tcell.StyleDefault.
		Background(ColorSchemes[Theme][Accent]).
		Foreground(ColorSchemes[Theme]["foreground"]).
		Attributes(tcell.AttrUnderline)

	style := defaultStyle
	if b.HasFocus() {
		style = selectedStyle
	}
	_, bg, _ := style.Decompose()
	b.SetBorderColor(ColorSchemes[Theme][Accent])
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
	*tview.Flex
}

func NewSidepane() *Sidepane {
	s := &Sidepane{
		Flex: tview.NewFlex(),
	}
	s.Box = tview.NewBox()
	s.SetDirection(tview.FlexRow)
	InitSidepaneStyle(s)

	s.SetBorderPadding(1, 1, 1, 1)

	for _, item := range [6]struct {
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
		s.AddItem(newButton(item.icon, item.label), 3, 1, false)
	}

	return s
}

func (s *Sidepane) GetButton(index int) *button {
	return s.GetItem(index).(*button)
}

func InitSidepaneStyle(s *Sidepane) *Sidepane {
	s.Box.SetBackgroundColor(BlendAccent)
	return s
}
