package viewty

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/unilibs/uniwidth"
)

func PrintString(painter Painter, x int, y int, style tcell.Style, str string) {
	i := 0
	for _, r := range str {
		painter.SetContent(x+i, y, r, nil, style)
		w := uniwidth.RuneWidth(r)
		for j := range w - 1 {
			painter.SetContent(x+j, y, ' ', nil, style)
		}
		i += w
	}
}

func FillRect(painter Painter, x int, y int, width int, height int, char rune, style tcell.Style) {
	for i := range height {
		for j := range width {
			painter.SetContent(x+j, y+i, char, nil, style)
		}
	}
}

func ClearRect(painter Painter, x int, y int, width int, height int, style tcell.Style) {
	FillRect(painter, x, y, width, height, ' ', style)
}

func PrintCenteredString(painter Painter, x int, y int, width int, style tcell.Style, str string) {
	ClearRect(painter, x, y, width, 1, style)
	strWidth := uniwidth.StringWidth(str)
	offset := (width - strWidth) / 2
	if offset < 0 {
		offset = 0
	}
	PrintString(painter, x+offset, y, style, str)
}

/**
 * A label is a string where th e ampersand character is used to preceed and
 * introduce a shortcut character which should be rendered differently.
 */
func LabelWidth(label string) int {
	w := 0
	for _, r := range label {
		if r != '&' {
			w += uniwidth.RuneWidth(r)
		}
	}
	return w
}

func PrintLabel(painter Painter, x int, y int, str string, style tcell.Style, shortcutStyle tcell.Style) {
	i := 0
	currentStyle := style
	for _, r := range str {
		if r == '&' {
			currentStyle = shortcutStyle
			continue
		}
		painter.SetContent(x+i, y, r, nil, currentStyle)
		w := uniwidth.RuneWidth(r)
		for j := range w - 1 {
			painter.SetContent(x+j, y, ' ', nil, currentStyle)
		}
		i += w
		currentStyle = style
	}
}
