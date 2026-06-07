package viewty

import (
	"log"
	"github.com/gdamore/tcell/v2"
)


type TextInput struct {
	*TextEditor
	userText       string
	history        []string
	historyPointer int

	done    func(tcell.Key)
	changed func(text string)
}

func NewTextInput() *TextInput {
	buffer := NewBufferFromString("", "")
	editor := NewTextEditor(buffer)
	buffer.Settings["ruler"] = false
	buffer.Settings["hidecursoronblur"] = true

	t := &TextInput{
		TextEditor: editor,
		historyPointer: -1,
	}
	t.SetOuterThis(t)
	return t
}

func (t *TextInput) SetHistory(historyText []string) {
	t.history = make([]string, len(historyText))
	copy(t.history, historyText)
	t.historyPointer = len(historyText)
}

func (t *TextInput) SetKeybindings(keybindings Keybindings) {
	t.TextEditor.SetKeybindings(keybindings)
}

func (t *TextInput) SetTextColor(foreground tcell.Color, background tcell.Color) {
	scheme := make(Colorscheme)
	scheme["default"] = tcell.StyleDefault.Foreground(foreground).Background(background)
	t.TextEditor.SetColorscheme(scheme)
}

func (t *TextInput) GetText() string {
	text := t.TextEditor.Buffer().Line(0)
	return text
}

func (t *TextInput) SetText(text string) {
	t.userText = text
	t.internalSetText(text)
}

func (t *TextInput) internalSetText(text string) {
	t.TextEditor.ActionController().DeleteLine()
	t.TextEditor.Buffer().Insert(t.TextEditor.Buffer().Start(), text)
	t.TextEditor.ActionController().StartOfLine()
	t.TextEditor.Relocate()
}

func (t *TextInput) HandleMouseEvent(mouseEvent MouseEvent) bool {
	if mouseEvent.SourceEvent().Buttons() == tcell.Button1 {
		app.Focus(t)
	}
	return false
}

func (t *TextInput) HandleKeyEvent(keyEvent KeyEvent) bool {
    event := keyEvent.SourceEvent()
	switch event.Key() {
	case tcell.KeyEnter, tcell.KeyEscape, tcell.KeyTab, tcell.KeyBacktab:
		if t.done != nil {
			t.done(event.Key())
		}
		return false
	case tcell.KeyUp:
		if len(t.history) != 0 && t.historyPointer > 0 {
			t.historyPointer--
			t.internalSetText(t.history[t.historyPointer])
			if t.changed != nil {
				t.changed(t.GetText())
			}
			return false
		}
	case tcell.KeyDown:
		if len(t.history) != 0 {
			if t.historyPointer < len(t.history)-1 {
				t.historyPointer++
				t.internalSetText(t.history[t.historyPointer])
			} else {
				t.internalSetText(t.userText)
				t.historyPointer = len(t.history)
				// Set to one past the end of history, so pressing up goes to the most recent entry
			}
			if t.changed != nil {
				t.changed(t.GetText())
			}
			return false
		}
	}

	t.TextEditor.Buffer().ClearModified()
	t.TextEditor.HandleKeyEvent(keyEvent)
	if t.TextEditor.Buffer().Modified() {
		t.TextEditor.Buffer().ClearModified()
		t.userText = t.GetText()
		if len(t.history) != 0 {
			t.historyPointer = len(t.history)
			// Reset history pointer to end of history, so pressing up goes to the most recent entry
		}
		if t.changed != nil {
			t.changed(t.GetText())
		}
	}
	return false
}

func (t *TextInput) SetDoneFunc(handler func(key tcell.Key)) *TextInput {
	t.done = handler
	return t
}

func (t *TextInput) SetChangedFunc(handler func(text string)) *TextInput {
	t.changed = handler
	return t
}
