package main

import (
	"fmt"

	tcell "github.com/gdamore/tcell/v2"
	tt "github.com/sedwards2009/termtronic"
)

func main() {
	app := tt.NewApplication()
	app.EnableLogging(true)

	rootFlexH := tt.NewHFlex()
	rootFlexH.SetGapSize(1)
	rootFlexH.SetName("Flex")

	var White = tcell.NewHexColor(0xf3f3f3).TrueColor()
    var Blue = tcell.NewHexColor(0x007ace).TrueColor()

    blueBox := tt.NewBox()
    blueBox.SetName("BlueBox")
	blueBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Blue))

	vFlex := tt.NewVFlex()
	vFlex.SetGapSize(1)

	whiteBox := tt.NewDotBox()
	whiteBox.SetName("WhiteDotBox")
	whiteBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(Blue).Background(White))
	vFlex.AddWidget(whiteBox, 10, 1)

	button := tt.NewButton()
	button.SetName("Button")
	button.SetText("Button clicks: 0")
	clickCount := 0
	button.SetOnClick(func(id string) {
		clickCount++
		button.SetText(fmt.Sprintf("Button clicks: %d", clickCount))
	})
	vFlex.AddWidget(button, 1, 0)

	rootFlexH.AddWidget(blueBox, 10, 2)
	rootFlexH.AddWidget(vFlex, 10, 1)

	app.SetRootWidget(rootFlexH)

    app.Run()
}
