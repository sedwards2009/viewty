package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type Widget interface {
	SetParent(parent Widget)
	Children() []Widget
	Parent() Widget
	Reposition(x, y, width, height int)
	Position() (int, int, int, int)
	Render(screen tcell.Screen)
}
