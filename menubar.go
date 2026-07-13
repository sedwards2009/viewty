package viewty

import (
	"os"
	"strings"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/unilibs/uniwidth"
)

type MenuBar struct {
	*WidgetBase
	menus        []*Menu
	selectedPath []int
	onClose      func()

	menuOverlay *menuOverlay
}

type Menu struct {
	ID        string
	Title     string
	Items     []*MenuItem
	charWidth int
	Shortcut  rune
}

type MenuItem struct {
	ID       string
	Title    string
	Shortcut string

	// Called when the menu item is activated. ID is the ID of the menu item.
	// The optional return value is the primative which should receive the
	// focus once the menu closes.
	Callback func(ID string)
}

const (
	MENU_BAR_SPACING = 1
	MENU_BAR_PADDING = 1
)

func NewMenuBar() *MenuBar {
	result := &MenuBar{
		WidgetBase:   NewWidgetBase(),
		selectedPath: []int{-1},
	}
	result.menuOverlay = newMenuOverlay(result)
	return result
}

func (menuBar *MenuBar) SetMenus(menus []*Menu) {
	menuBar.menus = menus
}

func (menuBar *MenuBar) Open(menuItem int) {
	menuBar.selectMenuBarItem(menuItem)
}

func (menuBar *MenuBar) Render(painter Painter) {
	_, _, width, _ := menuBar.Position()

	styles := menuBar.GetStyle("MenuBar", []string{})
	menuBarStyle := GetTCellStyle(styles, "menubarForegroundColor", "menubarBackgroundColor")

	for i := 0; i < width; i += 1 {
		painter.SetContent(i, 0, ' ', nil, menuBarStyle)
	}

	padding := ""
	for range MENU_BAR_PADDING {
		padding += " "
	}

	fg, bg, _ := menuBarStyle.Decompose()
	reverse := menuBarStyle.Foreground(bg).Background(fg)

	dx := 0
	isLinuxTerm := os.Getenv("TERM") == "linux"
	for i, menu := range menuBar.menus {
		title := menu.Title
		style := menuBarStyle
		isSelected := i == menuBar.selectedPath[0]
		if isSelected {
			style = reverse
		} else if isLinuxTerm {
			// In bare Linux TTYs, underline is often emulated by color 4 (Blue).
			// If the menu bar background is Blue, underlined mnemonics become invisible.
			// We fix this by using a different color (Yellow) and removing the underline attribute.
			title = strings.ReplaceAll(title, "[::u]", "[yellow]") // TODO
			title = strings.ReplaceAll(title, "[::U]", "[-]")
		}
		PrintString(painter, dx, 0, style, padding)
		dx += MENU_BAR_PADDING

		titleWidth := uniwidth.StringWidth(title)
		PrintString(painter, dx, 0, style, title)
		dx += titleWidth
		menu.charWidth = titleWidth

		PrintString(painter, dx, 0, style, padding)
		dx += MENU_BAR_PADDING
		dx += MENU_BAR_SPACING
	}
}

func menuWidthInCells(items []*MenuItem) int {
	titleWidth, shortcutWidth := measureWidths(items)
	return 1 + 1 + titleWidth + 2 + shortcutWidth + 1 + 1
}

func measureWidths(items []*MenuItem) (int, int) {
	maxTitleWidth := 0
	maxShortcutWidth := 0
	for _, item := range items {
		width := uniwidth.StringWidth(item.Title)
		if width > maxTitleWidth {
			maxTitleWidth = width
		}
		shortCutWidth := uniwidth.StringWidth(item.Shortcut)
		if shortCutWidth > maxShortcutWidth {
			maxShortcutWidth = shortCutWidth
		}
	}
	return maxTitleWidth, maxShortcutWidth
}

func (menuBar *MenuBar) HandleKeyEvent(keyEvent KeyEvent) bool {
	event := keyEvent.SourceEvent()

	if menuBar.selectedPath[0] == -1 {
		return false
	}

	selectedMenuIndex := menuBar.selectedPath[0]
	menu := menuBar.menus[selectedMenuIndex]
	selectedItemIndex := menuBar.selectedPath[1]
	item := menu.Items[selectedItemIndex]

	isAlt := (event.Modifiers() & tcell.ModAlt) != 0
	if isAlt {
		for i, menu := range menuBar.menus {
			if menu.Shortcut != 0 && menu.Shortcut == event.Rune() {
				menuBar.selectMenuBarItem(i)
				return true
			}
		}
	}

	switch event.Key() {
	case tcell.KeyEscape:
		menuBar.Close()

	case tcell.KeyLeft:
		selectedMenuIndex--
		if selectedMenuIndex < 0 {
			selectedMenuIndex = 0
		}
		menuBar.selectMenuBarItem(selectedMenuIndex)

	case tcell.KeyRight:
		selectedMenuIndex++
		if selectedMenuIndex == len(menuBar.menus) {
			selectedMenuIndex = len(menuBar.menus) - 1
		}
		menuBar.selectMenuBarItem(selectedMenuIndex)

	case tcell.KeyUp:
		menuBar.selectedPath[1] = nextMenuItem(menu.Items, selectedItemIndex, -1)

	case tcell.KeyDown:
		menuBar.selectedPath[1] = nextMenuItem(menu.Items, selectedItemIndex, 1)

	case tcell.KeyEnter:
		menuBar.executeItem(item)
	}
	return false
}

func (menuBar *MenuBar) executeItem(item *MenuItem) {
	if item.Title != "" {
		if item.Callback != nil {
			item.Callback(item.ID)
		}
		menuBar.selectMenuBarItem(-1)
		if menuBar.onClose != nil {
			menuBar.onClose()
		}
	}
}

func nextMenuItem(items []*MenuItem, selectedIndex int, direction int) int {
	next := func() {
		selectedIndex += direction

		if selectedIndex < 0 {
			selectedIndex = len(items) - 1
		}
		if selectedIndex >= len(items) {
			selectedIndex = 0
		}
	}
	next()

	for items[selectedIndex].Title == "" {
		next()
	}
	return selectedIndex
}

func (menuBar *MenuBar) HandleMouseEvent(mouseEvent MouseEvent) bool {
	x, y := mouseEvent.Position()
	return menuBar.internalHandleMouseEvent(mouseEvent, x, y)
}

func (menuBar *MenuBar) internalHandleMouseEvent(mouseEvent MouseEvent, x int, y int) bool {
	if y == 0 {
		if mouseEvent.IsLeftMousePress() {
			index, _ := menuBar.menuItemIndexAtX(x)
			if index != -1 {
				selectedIndex := menuBar.selectedPath[0]
				if selectedIndex != index {
					menuBar.selectMenuBarItem(index)
					// setFocus(menuBar)
					return false
				} else {
					menuBar.Close()
					return false
				}
			}
		}
		return false
	}

	selectedIndex := menuBar.selectedPath[0]
	if selectedIndex == -1 { // Is a menu open?
		return false
	}

	var selectedMenuItem *MenuItem = nil
	items := menuBar.menus[selectedIndex].Items
	selectedMenuItemIndex := y - 2
	menuLeft := menuBar.menuIndexLeft(selectedIndex)
	width := menuWidthInCells(items)
	if x >= menuLeft && x < (menuLeft+width) && selectedMenuItemIndex >= 0 && selectedMenuItemIndex < len(items) {
		selectedMenuItem = items[selectedMenuItemIndex]
	}

	if mouseEvent.IsLeftMousePress() {
		if selectedMenuItem != nil {
			menuBar.executeItem(selectedMenuItem)
		} else {
			menuBar.Close()
		}
		return false
	}

	if (mouseEvent.SourceEvent().Buttons()&tcell.Button1 != 0) && selectedMenuItem == nil {
		menuBar.Close()
		return false
	}

	if selectedMenuItem != nil {
		menuBar.selectedPath[1] = selectedMenuItemIndex
	}

	return false
}

func (menuBar *MenuBar) selectMenuBarItem(index int) {
	if index == -1 {
		menuBar.selectedPath = []int{-1}
		app.RemoveLayerWidget(menuBar.menuOverlay)
	} else {
		menuBar.selectedPath = []int{index, 0}
		app.AddLayerWidget(menuBar.menuOverlay)
		app.Focus(menuBar.menuOverlay)
	}
}

func (m *MenuBar) menuItemIndexAtX(posX int) (index int, leftX int) {
	x := 0
	for i, menu := range m.menus {
		if posX < x {
			return -1, -1
		}

		left := x
		x += MENU_BAR_PADDING
		x += menu.charWidth
		x += MENU_BAR_PADDING
		if posX < x {
			return i, left
		}

		x += MENU_BAR_SPACING
	}
	return -1, -1
}

func (m *MenuBar) menuIndexLeft(index int) int {
	x := 0
	for i, menu := range m.menus {
		if i == index {
			return x
		}

		x += MENU_BAR_PADDING
		x += menu.charWidth
		x += MENU_BAR_PADDING
		x += MENU_BAR_SPACING
	}
	return -1
}

func (menuBar *MenuBar) Close() {
	menuBar.selectMenuBarItem(-1)
	if menuBar.onClose != nil {
		menuBar.onClose()
	}
}

func (menuBar *MenuBar) SetOnClose(callback func()) {
	menuBar.onClose = callback
}

type menuOverlay struct {
	*WidgetBase
	menuBar *MenuBar
}

func newMenuOverlay(menuBar *MenuBar) *menuOverlay {
	return &menuOverlay{
		WidgetBase: NewWidgetBase(),
		menuBar:    menuBar,
	}
}

func (menuOverlay *menuOverlay) HandleKeyEvent(keyEvent KeyEvent) bool {
	return menuOverlay.menuBar.HandleKeyEvent(keyEvent)
}

func (menuOverlay *menuOverlay) HandleMouseEvent(mouseEvent MouseEvent) bool {
	x, y := mouseEvent.Position()
	rx, ry := menuOverlay.menuBar.PointToAbs(0, 0)
	menuOverlay.menuBar.internalHandleMouseEvent(mouseEvent, x-rx, y-ry)
	return false
}

func (menuOverlay *menuOverlay) Render(painter Painter) {
	menuBar := menuOverlay.menuBar
	if menuBar.selectedPath[0] == -1 {
		return
	}
	selectedIndex := menuBar.selectedPath[0]
	menu := menuBar.menus[selectedIndex]

	mx := menuBar.menuIndexLeft(selectedIndex)

	rx, ry := menuBar.PointToAbs(0, 0)
	menuOverlay.drawMenuItems(painter, rx+mx, ry, menu.Items, menuBar.selectedPath[1])
}

func (menuOverlay *menuOverlay) drawMenuItems(painter Painter, menuX int, menuY int, items []*MenuItem, selectedIndex int) {
	titleWidth, _ := measureWidths(items)
	topY := menuY + 1
	y := topY
	menuWidth := menuWidthInCells(items)

	styles := menuOverlay.menuBar.GetStyle("MenuBar", []string{})
	menuStyle := GetTCellStyle(styles, "menuForegroundColor", "menuBackgroundColor")
	menuSelectedStyle := GetTCellStyle(styles, "menuSelectedForegroundColor", "menuSelectedBackgroundColor")

	DrawCappedHorizontalLine(painter, menuX, y, menuWidth, menuStyle, menuStyle, '┌', '─', '┐')
	y++

	for i, item := range items {
		if item.Title == "" {
			DrawCappedHorizontalLine(painter, menuX, y, menuWidth, menuStyle, menuStyle, '├', '─', '┤')
		} else {
			textStyle := menuStyle
			if i == selectedIndex {
				textStyle = menuSelectedStyle
			}
			DrawCappedHorizontalLine(painter, menuX, y, menuWidth, menuStyle, textStyle, '│', ' ', '│')
			PrintString(painter, menuX+2, y, textStyle, item.Title)
			PrintString(painter, menuX+2+titleWidth+2, y, textStyle, item.Shortcut)
		}
		y++
	}

	DrawCappedHorizontalLine(painter, menuX, y, menuWidth, menuStyle, menuStyle, '└', '─', '┘')

	// Draw the drop shadow
	DrawDimVerticalLine(painter, menuX+menuWidth, topY+1, len(items)+1)
	DrawDimVerticalLine(painter, menuX+menuWidth+1, topY+1, len(items)+1)
	DrawDimHorizontalLine(painter, menuX+2, y+1, menuWidth)
	painter.Screen().HideCursor()
}

func DrawCappedHorizontalLine(painter Painter, x int, y int, width int, borderStyle tcell.Style, middleStyle tcell.Style, left rune,
	middle rune, right rune) {

	painter.SetContent(x, y, left, nil, borderStyle)
	for i := 1; i < width-1; i++ {
		painter.SetContent(x+i, y, middle, nil, middleStyle)
	}
	painter.SetContent(x+width-1, y, right, nil, borderStyle)
}

func DrawHorizontalLine(painter Painter, x int, y int, width int, style tcell.Style, char rune) {
	for i := 0; i < width; i++ {
		painter.SetContent(x+i, y, char, nil, style)
	}
}

func DimCell(painter Painter, x int, y int) {
	cellRune, cellRunes, style, _ := painter.GetContent(x, y)
	fg, bg, _ := style.Decompose()

	fgR, fgG, fgB := fg.TrueColor().RGB()
	fg = tcell.NewRGBColor(fgR/2, fgG/2, fgB/2)

	bgR, bgG, bgB := bg.TrueColor().RGB()
	bg = tcell.NewRGBColor(bgR/2, bgG/2, bgB/2)

	style = style.Foreground(fg).Background(bg)
	painter.SetContent(x, y, cellRune, cellRunes, style)
}

func DrawDimHorizontalLine(painter Painter, x int, y int, length int) {
	for i := 0; i < length; i++ {
		DimCell(painter, x+i, y)
	}
}

func DrawDimVerticalLine(painter Painter, x int, y int, length int) {
	for i := 0; i < length; i++ {
		DimCell(painter, x, y+i)
	}
}
