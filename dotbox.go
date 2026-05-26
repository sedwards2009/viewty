package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type DotBox struct {
	*WidgetBase
	backgroundStyle tcell.Style

	dotX int
	dotY int
}

func NewDotBox() *DotBox {
	return &DotBox{
		WidgetBase: NewWidgetBase(),
	}
}

func (d *DotBox) SetBackgroundStyle(style tcell.Style) {
	d.backgroundStyle = style
}

func (d *DotBox) Render(screen tcell.Screen) {
	x, y, w, h := d.Position()

	ClearRect(screen, x, y, w, h, d.backgroundStyle)

	screen.SetContent(x + d.dotX, y + d.dotY, 'X', nil, d.backgroundStyle)
}

func (d *DotBox) HandleMouseEvent(event *tcell.EventMouse, target Widget, phase EventPhase) bool {
	if event.Buttons() != tcell.Button1 {
		return false
	}
	x, y := event.Position()
	d.dotX = x - d.x
	d.dotY = y - d.y

	return false
}
