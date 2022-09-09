package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SudokuHeader struct {
	*tview.TextView
}

func NewSudokuHeader() *SudokuHeader {
	return &SudokuHeader{tview.NewTextView()}
}

// Draw draws SudokuHeader left aligned and at the bottom left of the
// bounding box.
func (h *SudokuHeader) Draw(screen tcell.Screen) {
	h.DrawForSubclass(screen, h)

	x, y, _, height := h.GetRect()
	y = y + height - 2

	// first row
	textStyle := tcell.StyleDefault.Background(Theme.background).Foreground(Theme.foreground)
	for i, r := range h.GetText(true) {
		screen.SetContent(x+i, y, r, nil, textStyle)
	}

	y++
	// second row
	underlineStyle := tcell.StyleDefault.Background(Theme.background).Foreground(Accent)
	for i := range h.GetText(true) {
		screen.SetContent(x+i, y, '▔', nil, underlineStyle)
	}
}

type second int

func (s second) String() string {
	hrs := s / 3600
	min := (s / 60) % 60
	sec := s % 60

	var ret strings.Builder
	if hrs != 0 {
		ret.WriteString(fmt.Sprintf("%dh ", hrs))
	}
	if (hrs != 0 && min == 0) || min != 0 {
		ret.WriteString(fmt.Sprintf("%dm ", min))
	}
	ret.WriteString(fmt.Sprintf("%ds", sec))
	return ret.String()
}
