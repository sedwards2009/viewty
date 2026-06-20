package viewty

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

func (b *Button) Render(painter Painter) {
	_, _, w, h := b.Position()

	styles := b.GetStyle("Button", []string{})
	var buttonStyle tcell.Style
	if app.HasFocus(b) {
		buttonStyle = GetTCellStyle(styles, "foregroundFocusColor", "backgroundFocusColor")
	} else {
		buttonStyle = GetTCellStyle(styles, "foregroundColor", "backgroundColor")
	}
	ClearRect(painter, 0, 0, w, h, buttonStyle)
	PrintCenteredString(painter, 0, 0, b.width, buttonStyle, b.text)
}

func (b *Button) HandleMouseEvent(mouseEvent MouseEvent) bool {
    if mouseEvent.SourceEvent().Buttons() == tcell.Button1 {
        app.Focus(b)
	    if b.onClick != nil {
			b.onClick(b.id)
		}
	}
	return false
}

func (b *Button) HandleKeyEvent(keyEvent KeyEvent) bool {
    key := keyEvent.SourceEvent().Key()
	if (key == tcell.KeyEnter || key == tcell.KeyRune && keyEvent.SourceEvent().Rune() == ' ') && b.onClick != nil {
		b.onClick(b.id)
	}
	return false
}
