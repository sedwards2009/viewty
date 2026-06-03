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

type MouseEvent interface {
	Position() (x int, y int)
	SourceEvent() *tcell.EventMouse
	Phase() EventPhase
	TargetWidget() Widget
}


type Widget interface {
	// The name of this widget.
	Name() string

	// Set the name of this widget
	SetName(name string)

	// Set the parent of this widget
	//
	// Note: This is mostly for internal use.
	SetParent(parent Widget)

	// The parent of this widget
	//
	// This will be `nil` if the widget has no parent or is the root widget.
	Parent() Widget

	// Size and position the widget relative to its parent.
	Reposition(x, y, width, height int)

	// Size and position of the widget relative to its parent.
	Position() (int, int, int, int)

	// Render the widget to the screen
	//
	// Note: This is for internal use.
	Render(screen Painter)

    // True if point is on the widget.
    //
    // The point is relative to the parent widget.
	ContainsPoint(x int, y int) bool

	PointToAbs(x int, y int) (ax int, ay int)

	// Get any child widget at the given point.
	//
    // The point is relative to the parent widget.
    ChildWidgetAt(x int, y int) Widget

	// Handle a mouse event
	//
	// Returns true if event handling should be stopped.
	HandleMouseEvent(mouseEvent MouseEvent) bool

	Focus()

	SetVisible(visible bool)
	IsVisible() bool
}
