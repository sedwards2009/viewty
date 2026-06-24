package main

import (
	"log"

	viewty "github.com/sedwards2009/viewty"
)

func main() {
	log.Println("Starting Scrollbar Demo...")

	app := viewty.NewApplication()

	verticalScrollbar := viewty.NewScrollbar()
	verticalScrollbar.SetChangedFunc(func(pos int) {
		log.Printf("Vertical scrollbar position: %d", pos)
	})

	innerLayout := viewty.NewHFlex()
	innerLayout.AddWidget(nil, 0, 1)
	innerLayout.AddWidget(verticalScrollbar, 1, 0)
	innerLayout.AddWidget(nil, 0, 1)

	layout := viewty.NewVFlex()
	layout.AddWidget(innerLayout, 0, 10)

	horizontalScrollbar := viewty.NewScrollbar()
	horizontalScrollbar.SetHorizontal(true)
	horizontalScrollbar.SetChangedFunc(func(pos int) {
		log.Printf("Horizontal scrollbar position: %d", pos)
	})

	layout.AddWidget(horizontalScrollbar, 1, 0)
	layout.AddWidget(nil, 0, 1)

	app.AddLayerWidget(layout)

	app.Run()
}
