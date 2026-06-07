package main

import (
	"fmt"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/viewty"
)

func main() {
	app := viewty.NewApplication()
	app.EnableLogging(true)

	rootFlexH := viewty.NewHFlex()
	rootFlexH.SetGapSize(1)
	rootFlexH.SetName("Flex")

	White := tcell.NewHexColor(0xf3f3f3).TrueColor()
    Blue := tcell.NewHexColor(0x007ace).TrueColor()
    Red := tcell.NewHexColor(0xce7a00).TrueColor()
    Green := tcell.NewHexColor(0x00ce7a).TrueColor()
    Yellow := tcell.NewHexColor(0xcece00).TrueColor()

    leftVFlex := viewty.NewVFlex()
    leftVFlex.SetName("leftVFlex")

    blueBox := viewty.NewBox()
    blueBox.SetName("BlueBox")
	blueBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Blue))
	leftVFlex.AddWidget(blueBox, 0, 1)

	textInput := viewty.NewTextInput()
	textInput.SetText("TextInput")
	leftVFlex.AddWidget(textInput, 1, 0)

	scrollArea := viewty.NewScrollArea()

	scrollContent := viewty.NewHFlex()
	scrollContent.SetName("scrollContent")
    redBox := viewty.NewBox()
    redBox.SetName("RedBox")
    redBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Red))
    scrollContent.AddWidget(redBox, 0, 1)

    greenBox := viewty.NewDotBox()
    greenBox.SetName("GreenBox")
    greenBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Green))
    scrollContent.AddWidget(greenBox, 0, 1)

    yellowBox := viewty.NewBox()
    yellowBox.SetName("yellowBox")
    yellowBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Yellow))
    scrollContent.AddWidget(yellowBox, 0, 1)

    scrollArea.SetContentWidget(scrollContent)
    scrollArea.SetMinimumSize(80, 60)

	leftVFlex.AddWidget(scrollArea, 0, 1)

	scrollLeftButton := viewty.NewButton()
	scrollLeftButton.SetText("Left")
	scrollLeftButton.SetOnClick(func (id string) {
		scrollArea.SetOffsetX(scrollArea.OffsetX()-1)
	})
	scrollRightButton := viewty.NewButton()
	scrollRightButton.SetText("Right")
	scrollRightButton.SetOnClick(func (id string) {
		scrollArea.SetOffsetX(scrollArea.OffsetX()+1)
	})

	buttonFlex := viewty.NewHFlex()
	buttonFlex.SetGapSize(1)
	buttonFlex.SetName("ButtonFlex")
	buttonFlex.AddWidget(scrollLeftButton, 0, 1)
    buttonFlex.AddWidget(scrollRightButton, 0, 1)
    leftVFlex.AddWidget(buttonFlex, 1, 0)

	vFlex := viewty.NewVFlex()
	vFlex.SetName("vFlex")
	vFlex.SetGapSize(1)

	whiteBox := viewty.NewDotBox()
	whiteBox.SetName("WhiteDotBox")
	whiteBox.SetBackgroundStyle(tcell.StyleDefault.Foreground(Blue).Background(White))
	vFlex.AddWidget(whiteBox, 10, 1)

	frame := viewty.NewFrame()
	frame.SetName("frame")
	frame.SetTitle("A Frame")
	vFlex.AddWidget(frame, 10, 1)

	button := viewty.NewButton()
	button.SetName("Button")
	button.SetText("Button clicks: 0")
	clickCount := 0
	button.SetOnClick(func(id string) {
		clickCount++
		button.SetText(fmt.Sprintf("Button clicks: %d", clickCount))
	})
	frame.SetContentWidget(button)

	rootFlexH.AddWidget(leftVFlex, 10, 2)
	rootFlexH.AddWidget(vFlex, 10, 1)

	app.SetRootWidget(rootFlexH)

    app.Run()
}
