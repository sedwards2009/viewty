package viewty

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

func (d *DotBox) Render(painter Painter) {
	_, _, w, h := d.Position()
	ClearRect(painter, 0, 0, w, h, d.backgroundStyle)
	painter.SetContent(d.dotX, d.dotY, 'X', nil, d.backgroundStyle)
}

func (d *DotBox) HandleMouseEvent(mouseEvent MouseEvent) bool {
	if mouseEvent.SourceEvent().Buttons() != tcell.Button1 {
		return false
	}
	x, y := mouseEvent.Position()
	d.dotX = x
	d.dotY = y
	return false
}
