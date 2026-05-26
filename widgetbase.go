package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type WidgetBase struct {
	parent Widget
	name string
	x      int
	y      int
	width  int
	height int
}

func NewWidgetBase() *WidgetBase {
	return &WidgetBase{}
}

func (w *WidgetBase) SetParent(parent Widget) {
	w.parent = parent
}

func (w *WidgetBase) Parent() Widget {
	return w.parent
}

func (w *WidgetBase) Name() string {
	return w.name
}

func (w *WidgetBase) SetName(name string) {
	w.name = name
}

func (w *WidgetBase) Reposition(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

func (w *WidgetBase) Position() (x int, y int, width int, height int) {
	return w.x, w.y, w.width, w.height
}

func (w *WidgetBase) ContainsPoint(x int, y int) bool {
	return x >= w.x && x < w.x+w.width && y >= w.y && y < w.y+w.height
}

func (w *WidgetBase) ChildWidgetAt(x int, y int) Widget {
	return nil
}

func (w *WidgetBase) Render(screen tcell.Screen) {
}

func (w *WidgetBase) HandleMouseEvent(event *tcell.EventMouse, target Widget, phase EventPhase) bool {
	return false
}

func (w *WidgetBase) Focus() {}
