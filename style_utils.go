package viewty

import (
	"github.com/gdamore/tcell/v2"
)

func GetTCellStyle(style StyleMap, foregroundName string, backgroundName string) tcell.Style {
	var foregroundColor tcell.Color
	foregroundColorAny, ok := style[foregroundName]
	if ok {
		foregroundColor, ok = foregroundColorAny.(tcell.Color)
	}
	if !ok {
		foregroundColor = tcell.NewHexColor(0xff00ff)
	}

	var backgroundColor tcell.Color
	backgroundColorAny, ok := style[backgroundName]
	if ok {
		backgroundColor, ok = backgroundColorAny.(tcell.Color)
	}
	if !ok {
		backgroundColor = tcell.NewHexColor(0xffff00)
	}
	return tcell.StyleDefault.Foreground(foregroundColor).Background(backgroundColor)
}
