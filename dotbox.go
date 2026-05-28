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

func (d *DotBox) Render(screen TranslateScreenWriter) {
	_, _, w, h := d.Position()
	ClearRect(screen, 0, 0, w, h, d.backgroundStyle)
	screen.SetContent(d.dotX, d.dotY, 'X', nil, d.backgroundStyle)
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
