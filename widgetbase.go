package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)


type WidgetBase struct {
	parent   Widget
	children []Widget
	x int
	y int
	width int
	height int
}

func NewWidgetBase() *WidgetBase {
	return &WidgetBase{}
}

func (w *WidgetBase) SetParent(parent Widget) {
	w.parent = parent
}

func (w *WidgetBase) Children() []Widget {
	return w.children
}

func (w *WidgetBase) Parent() Widget {
	return w.parent
}

func (w *WidgetBase) Reposition(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

func (w *WidgetBase) Position() (int, int, int, int) {
	return w.x, w.y, w.width, w.height
}

func (w *WidgetBase) Render(screen tcell.Screen) {
}
