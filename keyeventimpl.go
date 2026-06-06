package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type keyEventImpl struct {
	sourceEvent *tcell.EventKey
	phase       EventPhase
	targetWidget Widget
}

func (k *keyEventImpl) SourceEvent() *tcell.EventKey {
	return k.sourceEvent
}

func (k *keyEventImpl) Phase() EventPhase {
	return k.phase
}

func (k *keyEventImpl) TargetWidget() Widget {
	return k.targetWidget
}