package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type Box struct {
	*WidgetBase
	backgroundStyle tcell.Style
}

func NewBox() *Box {
	return &Box{
		WidgetBase: NewWidgetBase(),
	}
}

func (b *Box) SetBackgroundStyle(style tcell.Style) {
	b.backgroundStyle = style
}

func (b *Box) Render(screen tcell.Screen) {
	x, y, w, h := b.Position()

	for i := range h {
		for j := range w {
			screen.SetContent(x+j, y+i, 'A', nil, b.backgroundStyle)
		}
	}
}
