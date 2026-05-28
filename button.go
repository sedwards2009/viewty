package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type Button struct {
	*WidgetBase
	text string
	id string
	onClick func(id string)
}

func NewButton() *Button {
	return &Button{
		WidgetBase: NewWidgetBase(),
	}
}

func (b *Button) SetText(text string) {
	b.text = text
}

func (b *Button) SetOnClick(onClick func(id string)) {
	b.onClick = onClick
}

func (b *Button) SetId(id string) {
	b.id = id
}

func (b *Button) Id() string {
	return b.id
}

func (b *Button) Render(screen TranslateScreenWriter) {
	_, _, w, h := b.Position()

	var White = tcell.NewHexColor(0xf3f3f3).TrueColor()
    var Green = tcell.NewHexColor(0x0b835c).TrueColor()

	buttonStyle := tcell.StyleDefault.Foreground(White).Background(Green)

	for i := range h {
		for j := range w {
			screen.SetContent(j, i, ' ', nil, buttonStyle)
		}
	}
	PrintCenteredString(screen, 0, 0, b.width, buttonStyle, b.text)
}

func (b *Button) HandleMouseEvent(mouseEvent MouseEvent) bool {
	if mouseEvent.SourceEvent().Buttons() == tcell.Button1 && b.onClick != nil {
		b.onClick(b.id)
	}
	return false
}
