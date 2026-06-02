package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/unilibs/uniwidth"
)


func PrintString(painter Painter, x int, y int, style tcell.Style, str string) {
	i := 0
	for _, r := range str {
		painter.SetContent(x+i, y, r, nil, style)
        w := uniwidth.RuneWidth(r)
        for j := range w-1 {
          painter.SetContent(x+j, y, ' ', nil, style)
        }
        i += w
	}
}

func ClearRect(painter Painter, x int, y int, width int, height int, style tcell.Style) {
	for i := range height {
		for j := range width {
			painter.SetContent(x+j, y+i, ' ', nil, style)
		}
	}
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
