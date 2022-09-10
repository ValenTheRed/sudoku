package main

import (
	"github.com/rivo/tview"
)

type Sidepane struct {
	*tview.Box
}

func NewSidepane() *Sidepane {
	s := &Sidepane{
		Box: tview.NewBox(),
	}
	return s
}
