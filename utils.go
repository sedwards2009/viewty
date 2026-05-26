package termtronic

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/unilibs/uniwidth"
)


func PrintString(screen tcell.Screen, x int, y int, style tcell.Style, str string) {
	i := 0
	for _, r := range str {
		screen.SetContent(x+i, y, r, nil, style)
        w := uniwidth.RuneWidth(r)
        for j := range w-1 {
          screen.SetContent(x+j, y, ' ', nil, style)
        }
        i += w
	}
}

func ClearRect(screen tcell.Screen, x int, y int, width int, height int, style tcell.Style) {
	for i := range height {
		for j := range width {
			screen.SetContent(x+j, y+i, ' ', nil, style)
		}
	}
}

func PrintCenteredString(screen tcell.Screen, x int, y int, width int, style tcell.Style, str string) {
	ClearRect(screen, x, y, width, 1, style)
	strWidth := uniwidth.StringWidth(str)
	offset := (width - strWidth) / 2
	if offset < 0 {
		offset = 0
	}
	PrintString(screen, x+offset, y, style, str)
}
