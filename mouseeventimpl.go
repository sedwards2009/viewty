package viewty

import (
	tcell "github.com/gdamore/tcell/v2"
)


type mouseEventImpl struct {
  	x int
    y int
    sourceEvent *tcell.EventMouse
    phase EventPhase
    targetWidget Widget
    previousSourceEvent *tcell.EventMouse
}

func (m *mouseEventImpl) Position() (x int, y int) {
    return m.x, m.y
}

func (m *mouseEventImpl) SourceEvent() *tcell.EventMouse {
    return m.sourceEvent
}

func (m *mouseEventImpl) PreviousSourceEvent() *tcell.EventMouse {
    return m.previousSourceEvent
}

func (m *mouseEventImpl) Phase() EventPhase {
    return m.phase
}

func (m *mouseEventImpl) TargetWidget() Widget {
    return m.targetWidget
}

func (m *mouseEventImpl) IsLeftMousePress() bool {
	if m.previousSourceEvent == nil {
		return false
	}
	return (m.previousSourceEvent.Buttons() & tcell.Button1 == 0) && (m.sourceEvent.Buttons() & tcell.Button1 != 0)
}

func (m *mouseEventImpl) IsLeftMouseRelease() bool {
	if m.previousSourceEvent == nil {
		return false
	}
	return (m.previousSourceEvent.Buttons() & tcell.Button1 != 0) && (m.sourceEvent.Buttons() & tcell.Button1 == 0)
}
