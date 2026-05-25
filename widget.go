package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type EventPhase int

const (
	EVENT_PHASE_CAPTURE EventPhase = 0
	EVENT_PHASE_BUBBLE EventPhase = 1
	EVENT_PHASE_TARGET EventPhase = 2
)

type Widget interface {
	Name() string
	SetName(name string)
	SetParent(parent Widget)
	Parent() Widget
	Reposition(x, y, width, height int)
	Position() (int, int, int, int)
	Render(screen tcell.Screen)
    ContainsPoint(x int, y int) bool
    ChildWidgetAt(x int, y int) Widget

	// Returns true if event handling should be stopped.
	HandleMouseEvent(event *tcell.EventMouse, target Widget, phase EventPhase) bool
}
