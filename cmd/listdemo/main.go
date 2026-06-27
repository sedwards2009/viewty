package main

import (
	"fmt"

	"github.com/sedwards2009/viewty"
)

func main() {
	app := viewty.NewApplication()
	app.EnableLogging(true)

	// Create the list widget
	list := viewty.NewList()
	list.SetName("My List")

	// Set up list items
	items := []viewty.ListItem{
		{Text: "Item 1", ID: "item1"},
		{Text: "Item 2", ID: "item2"},
		{Text: "Item 3", ID: "item3"},
		{Text: "Item 4", ID: "item4"},
		{Text: "Item 5", ID: "item5"},
		{Text: "Item 6", ID: "item6"},
		{Text: "Item 7", ID: "item7"},
		{Text: "Item 8", ID: "item8"},
		{Text: "Item 9", ID: "item9"},
		{Text: "Item 10", ID: "item10"},
		{Text: "Item 11", ID: "item11"},
		{Text: "Item 12", ID: "item12"},
		{Text: "Item 13", ID: "item13"},
		{Text: "Item 14", ID: "item14"},
		{Text: "Item 15", ID: "item15"},
		{Text: "Item 16", ID: "item16"},
		{Text: "Item 17", ID: "item17"},
		{Text: "Item 18", ID: "item18"},
		{Text: "Item 19", ID: "item19"},
		{Text: "Item 20", ID: "item20"},
		{Text: "Item 21", ID: "item21"},
		{Text: "Item 22", ID: "item22"},
		{Text: "Item 23", ID: "item23"},
		{Text: "Item 24", ID: "item24"},
		{Text: "Item 25", ID: "item25"},
		{Text: "Item 26", ID: "item26"},
		{Text: "Item 27", ID: "item27"},
		{Text: "Item 28", ID: "item28"},
		{Text: "Item 29", ID: "item29"},
		{Text: "Item 30", ID: "item30"},
	}
	list.SetListItems(items)

	// Create a button to clear the list
	clearButton := viewty.NewButton()
	clearButton.SetText("Clear List")
	clearButton.SetOnClick(func(id string) {
		list.SetListItems(nil)
	})

	// Add widgets to a flex layout
	vFlex := viewty.NewVFlex()
	vFlex.SetGapSize(1)
	vFlex.AddWidget(list, 10, 1)
	vFlex.AddWidget(clearButton, 0, 1)

	// Add to application
	app.AddLayerWidget(vFlex)

	fmt.Println("List demo - click the list to focus it, click the button to clear")

	app.Run()
}
