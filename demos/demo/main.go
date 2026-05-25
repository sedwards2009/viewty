package main

import (
	tcell "github.com/gdamore/tcell/v2"
	tt "github.com/sedwards2009/termtronic"
)

func main() {
	app := tt.NewApplication()
	app.EnableLogging(true)

	rootFlexWidget := tt.NewHFlex()
	rootFlexWidget.SetGapSize(1)
	rootFlexWidget.SetName("Flex")

	var White = tcell.NewHexColor(0xf3f3f3).TrueColor()
    var Blue = tcell.NewHexColor(0x007ace).TrueColor()

    blueBox := tt.NewBox()
    blueBox.SetName("BlueBox")
	blueBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Blue))

	whiteBox := tt.NewDotBox()
	whiteBox.SetName("WhiteDotBox")
	whiteBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(Blue).Background(White))
	rootFlexWidget.AddWidget(blueBox, 10, 1)
	rootFlexWidget.AddWidget(whiteBox, 10, 2)

	app.SetRootWidget(rootFlexWidget)

    app.Run()
}
