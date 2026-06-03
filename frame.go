package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
)

type Frame struct {
	*WidgetBase
	title string
	padding int
	drawFrame bool
	contentWidget Widget
}

func NewFrame() *Frame {
	return &Frame{
		WidgetBase: NewWidgetBase(),
		drawFrame: true,
	}
}

func (f *Frame) SetTitle(title string) {
	f.title = title
}

func (f *Frame) Title() string {
	return f.title
}

func (f *Frame) SetDrawFrame(on bool) {
	f.drawFrame = on
}

func (f *Frame) IsDrawFrame() bool {
	return f.drawFrame
}

func (f *Frame) SetPadding(padding int) {
	f.padding = padding
}

func (f *Frame) Padding() int {
	return f.padding
}

func (f *Frame) SetContentWidget(contentWidget Widget) {
	f.contentWidget = contentWidget
}

func (f *Frame) Reposition(x int, y int, width int, height int) {
	f.WidgetBase.Reposition(x, y, width, height)
	if f.contentWidget != nil {
		contentWidth := width - 2 * f.padding
		contentHeight := height - 2 * f.padding
		f.contentWidget.Reposition(0, 0, contentWidth, contentHeight)
	}
}

func (f *Frame) Render(painter Painter) {
	if f.contentWidget != nil {
		_, _, width, height := f.contentWidget.Position()
		f.contentWidget.Render(painter.Translate(f.padding, f.padding).ApplyClipArea(0, 0, width, height))
	}

	var White = tcell.NewHexColor(0xf3f3f3).TrueColor()
    var Black = tcell.NewHexColor(0x000000).TrueColor()

	style := tcell.StyleDefault.Foreground(White).Background(Black)

	painter.SetContent(0, 0, 'X', nil, style)

	FillRect(painter, 0, 0, f.width, 1, '\u2500', style)
	FillRect(painter, 0, f.height-1, f.width, 1, '\u2500', style)
	FillRect(painter, 0, 0, 1, f.height, '\u2502', style)
	FillRect(painter, f.width-1, 0, 1, f.height, '\u2502', style)

	// 0x250C BOX DRAWINGS LIGHT DOWN AND RIGHT
	painter.SetContent(0, 0, '\u250C', nil, style)

	// 0x2510 BOX DRAWINGS LIGHT DOWN AND LEFT
	painter.SetContent(f.width -1, 0, '\u2510', nil, style)

	// 0x2514 BOX DRAWINGS LIGHT UP AND RIGHT
	painter.SetContent(0, f.height-1, '\u2514', nil, style)

	// 0x2518 BOX DRAWINGS LIGHT UP AND LEFT
	painter.SetContent(f.width-1, f.height-1, '\u2518', nil, style)

	if f.title != "" {
		PrintString(painter, 1, 0, style, f.title)
	}
}
