package viewty

import (
	tcell "github.com/gdamore/tcell/v2"
)

type ListItem struct {
	Text string
	ID   string
}

type List struct {
	*Flex
	listArea *listArea
	scrollbar *Scrollbar
}

func NewList() *List {
	vFlex := NewHFlex()
	listArea := newListArea()
	scrollbar := NewScrollbar()
	vFlex.AddWidget(listArea, 0, 1)
	vFlex.AddWidget(scrollbar, 1, 0)

	listArea.SetOnScroll(func(position int) {
		scrollbar.SetThumbPosition(position)
	})

	scrollbar.SetOnChange(func(position int) {
		listArea.SetScrollY(position)
	})

	return &List{
		Flex: vFlex,
		listArea: listArea,
		scrollbar: scrollbar,
	}
}

func (l *List) Reposition(x, y, width, height int) {
	l.scrollbar.SetVisible(len(l.listArea.items) > height)
	l.Flex.Reposition(x, y, width, height)
	l.scrollbar.SetThumbSize(height)
}

func (l *List) SetListItems(items []ListItem) {
	l.listArea.SetListItems(items)
	l.scrollbar.SetMax(len(items))
	_, _, _, height := l.Position()
	l.scrollbar.SetThumbSize(height)
	l.scrollbar.SetVisible(len(items) > height)
	l.Relayout()
}

func (l *List) SetOnSelected(onSelected func(id string)) {
	l.listArea.SetOnSelected(onSelected)
}

func (l *List) GetScrollY() int {
	return l.listArea.GetScrollY()
}

func (l *List) SetScrollY(scrollY int) {
	l.listArea.SetScrollY(scrollY)
}

func (l *List) SelectedID() string {
	return l.listArea.SelectedID()
}

func (l *List) ScrollIntoView(id string) {
	l.listArea.ScrollIntoView(id)
}

//-------------------------------------------------------------------------

type listArea struct {
	*WidgetBase
	items      []ListItem
	selectedID string
	scrollY    int
	onSelected func(id string)
	onScroll   func(position int)
}

func newListArea() *listArea {
	return &listArea{
		WidgetBase: NewWidgetBase(),
	}
}

func (l *listArea) Render(painter Painter) {
	_, _, widgetWidth, widgetHeight := l.Position()

	styles := l.GetStyle("List", []string{})
	var listStyle tcell.Style
	if app.HasFocus(l) {
		listStyle = GetTCellStyle(styles, "itemFocusForegroundColor", "itemFocusBackgroundColor")
	} else {
		listStyle = GetTCellStyle(styles, "itemForegroundColor", "itemBackgroundColor")
	}

	ClearRect(painter, 0, 0, widgetWidth, widgetHeight, listStyle)

	for y := range len(l.items) - l.scrollY {
		if y >= widgetHeight {
			break // Don't overflow the widget
		}

		var itemStyle tcell.Style
		item := l.items[l.scrollY + y]
		if l.selectedID == item.ID {
			itemStyle = GetTCellStyle(styles, "selectedItemForegroundColor", "selectedItemBackgroundColor")
		} else {
			itemStyle = listStyle
		}

		PrintString(painter, 0, y, itemStyle, item.Text)

		// Fill the rest of the line with spaces
		if len(item.Text) < widgetWidth-1 {
			for x := len(item.Text); x < widgetWidth; x++ {
				painter.SetContent(x, y, ' ', nil, itemStyle)
			}
		}
	}
}

func (l *listArea) SetListItems(items []ListItem) {
	l.items = items
}

func (l *listArea) SetOnSelected(onSelected func(id string)) {
	l.onSelected = onSelected
}

func (l *listArea) SetOnScroll(onScroll func(position int)) {
    l.onScroll = onScroll
}

func (l *listArea) GetScrollY() int {
	return l.scrollY
}

func (l *listArea) SetScrollY(scrollY int) {
	l.scrollY = scrollY
}

func (l *listArea) SelectedID() string {
	return l.selectedID
}

func (l *listArea) getSelectedIndex() int {
	return l.getIndexByID(l.selectedID)
}

func (l *listArea) getIndexByID(id string) int {
	for i, item := range l.items {
		if item.ID == id {
			return i
		}
	}
	return -1
}

func (l *listArea) HandleKeyEvent(keyEvent KeyEvent) bool {
	key := keyEvent.SourceEvent().Key()

	switch key {
	case tcell.KeyUp:
		if l.selectedID == "" {
			// No item selected, select the last item (above the top)
			if len(l.items) > 0 {
				l.selectedID = l.items[len(l.items)-1].ID
				if l.onSelected != nil {
					l.onSelected(l.selectedID)
				}
			}
		} else {
			// Select the item above
			selectedIndex := l.getSelectedIndex()
			if selectedIndex > 0 {
				newSelectedID := l.items[selectedIndex-1].ID
				l.selectedID = newSelectedID
				if l.onSelected != nil {
					l.onSelected(l.selectedID)
				}
				l.ScrollIntoView(newSelectedID)
			}
		}
	case tcell.KeyPgUp:
		if l.selectedID != "" {
			// Move selected item to the top of the visible area
			selectedIndex := l.getSelectedIndex()
			_, _, _, widgetHeight := l.Position()
			newSelectedIndex := max(0, selectedIndex - widgetHeight + 1)
			newSelectedID := l.items[newSelectedIndex].ID
			l.selectedID = newSelectedID
			if l.onSelected != nil {
				l.onSelected(l.selectedID)
			}
			l.ScrollIntoView(l.selectedID)
		}
	case tcell.KeyPgDn:
		if l.selectedID != "" {
			selectedIndex := l.getSelectedIndex()
			_, _, _, widgetHeight := l.Position()
			newSelectedIndex := min(len(l.items)-1, selectedIndex + widgetHeight - 1)
			newSelectedID := l.items[newSelectedIndex].ID
			l.selectedID = newSelectedID
			if l.onSelected != nil {
				l.onSelected(l.selectedID)
			}
			l.ScrollIntoView(newSelectedID)
		}
	case tcell.KeyDown:
		if l.selectedID == "" {
			// No item selected, select the first item
			if len(l.items) > 0 {
				l.selectedID = l.items[0].ID
				if l.onSelected != nil {
					l.onSelected(l.selectedID)
				}
			}
		} else {
			// Select the item below
			selectedIndex := l.getSelectedIndex()
			if selectedIndex < len(l.items)-1 {
				newSelectedID := l.items[selectedIndex+1].ID
				l.selectedID = newSelectedID
				if l.onSelected != nil {
					l.onSelected(l.selectedID)
				}
				l.ScrollIntoView(newSelectedID)
			}
		}
	}

	return false
}

func (l *listArea) setScrollY(scrollY int) {
	if l.scrollY == scrollY {
		return
	}
	l.scrollY = scrollY
	if l.onScroll != nil {
		l.onScroll(scrollY)
	}
}

func (l *listArea) ScrollIntoView(id string) {
	index := l.getIndexByID(id)
	if index < 0 {
		return	// Coding error somewhere
	}
	_, _, _, widgetHeight := l.Position()

	above := index - l.scrollY
	if above < 0 {
		l.setScrollY(index)
	}

	under := index - (l.scrollY + widgetHeight)
	if under >= 0 {
		l.setScrollY(index - widgetHeight + 1)
	}
}

func (l *listArea) Reposition(x, y, width, height int) {
	l.WidgetBase.Reposition(x, y, width, height)
	if l.scrollY > len(l.items) - height {
		l.setScrollY(max(0, len(l.items) - height-1))
	}
}

func (l *listArea) HandleMouseEvent(mouseEvent MouseEvent) bool {
	switch mouseEvent.SourceEvent().Buttons() {
	case tcell.Button1:
		_, relY := mouseEvent.Position()
		app.Focus(l)

		clickIndex := relY + l.scrollY
		if clickIndex >=0 && clickIndex < len(l.items) {
			clickedItemID := l.items[relY].ID
			if l.selectedID != clickedItemID {
				l.selectedID = clickedItemID
				if l.onSelected != nil {
					l.onSelected(clickedItemID)
				}
			}
		}

	case tcell.WheelUp:
		// Scroll up by one line
		if l.scrollY > 0 {
			l.setScrollY(l.scrollY-1)
		}

	case tcell.WheelDown:
		// Scroll down by one line
		_, _, _, widgetHeight := l.Position()
		if l.scrollY < len(l.items)- widgetHeight {
			l.setScrollY(l.scrollY+1)
		}
	}

	return false
}
