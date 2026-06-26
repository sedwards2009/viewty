package viewty

import (
	tcell "github.com/gdamore/tcell/v2"
	"github.com/unilibs/uniwidth"
)

type CheckBox struct {
	*WidgetBase
	checked bool
	id      string
	label   string
	onChange func(id string)
}

func NewCheckBox() *CheckBox {
	return &CheckBox{
		WidgetBase: NewWidgetBase(),
	}
}

func (c *CheckBox) SetChecked(checked bool) {
	c.checked = checked
}

func (c *CheckBox) SetLabel(label string) {
	c.label = label
}

func (c *CheckBox) SetId(id string) {
	c.id = id
}

func (c *CheckBox) Id() string {
	return c.id
}

func (c *CheckBox) SetOnChange(onChange func(id string)) {
	c.onChange = onChange
}

func (c *CheckBox) Render(painter Painter) {
	_, _, w, h := c.Position()

	styles := c.GetStyle("CheckBox", []string{})
	var checkboxStyle tcell.Style
	if app.HasFocus(c) {
		checkboxStyle = GetTCellStyle(styles, "foregroundFocusColor", "backgroundFocusColor")
	} else {
		checkboxStyle = GetTCellStyle(styles, "checkboxForegroundColor", "checkboxBackgroundColor")
	}

	ClearRect(painter, 0, 0, w, h, checkboxStyle)

	// Draw checkbox
	x := 1
	y := (h - 1) / 2
	checkboxText := "[ ]"
	if c.checked {
		checkboxText = "[X]"
	}
	PrintString(painter, x, y, checkboxStyle, checkboxText)

	// Draw label if exists
	if c.label != "" {
		checkboxWidth := uniwidth.StringWidth(checkboxText)
		x += checkboxWidth + 1 // checkbox width + 1 space
		var labelStyle tcell.Style
		if app.HasFocus(c) {
			labelStyle = GetTCellStyle(styles, "labelFocusForegroundColor", "labelFocusBackgroundColor")
		} else {
			labelStyle = GetTCellStyle(styles, "labelForegroundColor", "labelBackgroundColor")
		}
		PrintString(painter, x, y, labelStyle, c.label)
	}
}

func (c *CheckBox) HandleMouseEvent(mouseEvent MouseEvent) bool {
	if mouseEvent.SourceEvent().Buttons() == tcell.Button1 {
		app.Focus(c)
		c.checked = !c.checked
		if c.onChange != nil {
			c.onChange(c.id)
		}
	}
	return false
}

func (c *CheckBox) HandleKeyEvent(keyEvent KeyEvent) bool {
	key := keyEvent.SourceEvent().Key()
	if (key == tcell.KeyEnter || key == tcell.KeyRune && keyEvent.SourceEvent().Rune() == ' ') {
		c.checked = !c.checked
		if c.onChange != nil {
			c.onChange(c.id)
		}
	}
	return false
}
