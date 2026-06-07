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
}

func (m *mouseEventImpl) Position() (x int, y int) {
    return m.x, m.y
}

func (m *mouseEventImpl) SourceEvent() *tcell.EventMouse {
    return m.sourceEvent
}

func (m *mouseEventImpl) Phase() EventPhase {
    return m.phase
}

func (m *mouseEventImpl) TargetWidget() Widget {
    return m.targetWidget
}
