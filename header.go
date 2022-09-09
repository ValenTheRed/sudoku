package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SudokuHeader struct {
	*tview.TextView
	align int
}

func NewSudokuHeader() *SudokuHeader {
	return &SudokuHeader{
		TextView: tview.NewTextView(),
		align:    tview.AlignLeft,
	}
}

func (h *SudokuHeader) SetTextAlign(align int) *SudokuHeader {
	h.align = align
	return h
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

// Timer counts the number of seconds it took to complete the puzzle.
//
// I've not included any way to check whether the Timer is running or
// not since, Timer will start at the beginning of the application and
// stop only when the puzzle is correctly completed or the application
// exits.
type Timer struct {
	*SudokuHeader
	elapsed second
	stopCh  chan struct{}
}

// NewTimer returns a new initialised Timer.
//
// NOTE: NewTimer doesn't set handler for queueing redraws via
// SetChangedFunc(), i.e.
//	Timer.SetChangedFunc(func() {
//		app.Draw()
//	})
// users of NewTimer have to do that themselves. Without installing the
// handler, Timer text would not be updated.
func NewTimer() *Timer {
	t := &Timer{
		SudokuHeader: NewSudokuHeader(),
		stopCh:       make(chan struct{}),
	}
	t.SetText(t.elapsed.String())
	return t
}

func (t *Timer) Start() {
	go worker(func() {
		t.elapsed++
		t.SetText(t.elapsed.String())
	}, t.stopCh)
}

func (t *Timer) Stop() {
	t.stopCh <- struct{}{}
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

// Worker executes work after every second. If a message is sent to
// quit, Worker returns.
func worker(work func(), quit <-chan struct{}) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			work()
		case <-quit:
			return
		}
	}
}
