package viewty

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/sedwards2009/viewty/micro/action"
	"github.com/sedwards2009/viewty/micro/buffer"
	"github.com/sedwards2009/viewty/micro/config"
	"github.com/sedwards2009/viewty/micro/display"
	"github.com/sedwards2009/viewty/runtime"
)

type TextEditor struct {
    *WidgetBase
    outerThis Widget
	buffer    *buffer.Buffer
	bufWindow *display.BufWindow
	bufPane   *action.BufPane
}

type ActionController struct {
	*action.BufPane
}

func NewTextEditor(buffer *buffer.Buffer) *TextEditor {
	v := &TextEditor{
	    WidgetBase: NewWidgetBase(),
		buffer:    buffer,
		bufWindow: display.NewBufWindow(0, 0, 10, 10, buffer),
	}
	v.outerThis = v

	buffer.RegisterRedrawCallback(func() {
	  app.ForceRender()
	})

	v.bufPane = action.NewBufPane(v.buffer, v.bufWindow)
	v.buffer.UpdateRules()
	return v
}

func (t *TextEditor) SetOuterThis(w Widget) {
    t.outerThis = w
}

func (t *TextEditor) Render(painter Painter) {
	_, _, width, height := t.Position()
	foreground := tcell.NewHexColor(0xf3f3f3).TrueColor()
	background := tcell.NewHexColor(0x0b835c).TrueColor()

	buttonStyle := tcell.StyleDefault.Foreground(foreground).Background(background)
painter.SetContent(0,0,'?', nil,buttonStyle)
	absX, absY :=t.PointToAbs(0, 0)
	t.bufWindow.X = absX
	t.bufWindow.Y = absY
	if t.bufWindow.Width != width || t.bufWindow.Height != height {
		t.bufWindow.Resize(width, height)
	}
	t.bufWindow.Display(painter.Screen(), app.HasFocus(t.outerThis))
}

func (t *TextEditor) HandleKeyEvent(keyEvent KeyEvent) bool {
	takeFocus := func() {
	    app.Focus(t.outerThis)
	}
	t.bufPane.HandleEvent(keyEvent.SourceEvent(), takeFocus)
	return false
}

// MouseHandler returns the mouse handler for this primitive.
func (t *TextEditor) HandleMouseEvent(mouseEvent MouseEvent) bool {
	takeFocus := func() {
		app.Focus(t.outerThis)
	}
	t.bufPane.HandleEvent(mouseEvent.SourceEvent(), takeFocus)
	return false
}

// func (t *TextEditor) PasteHandler() func(pastedText string, setFocus func(p tview.Primitive)) {
// 	return t.WrapPasteHandler(func(pastedText string, setFocus func(p tview.Primitive)) {
// 	t.bufPane.PasteString(pastedText)
// 	})
// }

func (t *TextEditor) SetColorscheme(cs Colorscheme) {
	t.bufWindow.Colorscheme = config.Colorscheme(cs)
	t.buffer.UpdateRules()
}

func (t *TextEditor) Buffer() *buffer.Buffer {
	return t.buffer
}

func (t *TextEditor) Cursor() *buffer.Cursor {
	return t.bufPane.Cursor
}

func (t *TextEditor) Relocate() {
	t.bufWindow.Relocate()
}

func (t *TextEditor) ActionController() *ActionController {
	return &ActionController{t.bufPane}
}

type Keybindings struct {
	*action.KeyTree
}

func ParseKeybindings(config map[string]string) Keybindings {
	return Keybindings{action.BindingMappingToKeyTree(config)}
}

func (t *TextEditor) SetKeybindings(keybindings Keybindings) {
	t.bufPane.SetBindings(keybindings.KeyTree)
}

func (t *TextEditor) GoToLoc(loc buffer.Loc) {
	t.bufPane.GotoLoc(loc)
}

func NewBufferFromString(content string, path string) *buffer.Buffer {
	return buffer.NewBufferFromString(content, path)
}

type Colorscheme config.Colorscheme

func (colorscheme Colorscheme) GetColor(color string) tcell.Style {
	return config.Colorscheme(colorscheme).GetColor(color)
}

func LoadInternalColorscheme(name string) (Colorscheme, bool) {
	data, err := runtime.Asset("runtime/colorschemes/" + name + ".micro")
	if err != nil {
		return nil, false
	}
	return ParseColorscheme(string(data)), true
}

func ParseColorscheme(data string) Colorscheme {
	return Colorscheme(config.ParseColorscheme(data))
}

func ListColorschemes() []string {
	files, err := runtime.AssetDir("runtime/colorschemes")
	if err != nil {
		return nil
	}
	var schemes []string
	for _, f := range files {
		schemes = append(schemes, f[:len(f)-6]) // Remove .micro extension
	}
	return schemes
}

func ListSyntaxes() []string {
	files, err := runtime.AssetDir("runtime/syntax")
	if err != nil {
		return nil
	}
	var syntaxes []string
	for _, f := range files {
		if strings.HasSuffix(f, ".yaml") {
			syntaxes = append(syntaxes, f[:len(f)-5]) // Remove .yaml extension
		}
	}
	return syntaxes
}

func init() {
	config.InitRuntimeFiles()
}

type Action func() bool

func (t *TextEditor) MapActionNameToAction(name string) Action {
	if f, ok := action.BufKeyActions[name]; ok {
		return func() bool {
			return f(t.bufPane)
		}
	}
	return nil
}

// Actions
const (
	ActionCursorUp                = "CursorUp"
	ActionCursorDown              = "CursorDown"
	ActionCursorPageUp            = "CursorPageUp"
	ActionCursorPageDown          = "CursorPageDown"
	ActionCursorLeft              = "CursorLeft"
	ActionCursorRight             = "CursorRight"
	ActionCursorStart             = "CursorStart"
	ActionCursorEnd               = "CursorEnd"
	ActionSelectToStart           = "SelectToStart"
	ActionSelectToEnd             = "SelectToEnd"
	ActionSelectUp                = "SelectUp"
	ActionSelectDown              = "SelectDown"
	ActionSelectLeft              = "SelectLeft"
	ActionSelectRight             = "SelectRight"
	ActionWordRight               = "WordRight"
	ActionWordLeft                = "WordLeft"
	ActionSelectWordRight         = "SelectWordRight"
	ActionSelectWordLeft          = "SelectWordLeft"
	ActionDeleteWordRight         = "DeleteWordRight"
	ActionDeleteWordLeft          = "DeleteWordLeft"
	ActionSelectLine              = "SelectLine"
	ActionSelectToStartOfLine     = "SelectToStartOfLine"
	ActionSelectToEndOfLine       = "SelectToEndOfLine"
	ActionParagraphPrevious       = "ParagraphPrevious"
	ActionParagraphNext           = "ParagraphNext"
	ActionInsertNewline           = "InsertNewline"
	ActionInsertSpace             = "InsertSpace"
	ActionBackspace               = "Backspace"
	ActionDelete                  = "Delete"
	ActionInsertTab               = "InsertTab"
	ActionCenter                  = "Center"
	ActionUndo                    = "Undo"
	ActionRedo                    = "Redo"
	ActionCopy                    = "Copy"
	ActionCut                     = "Cut"
	ActionCutLine                 = "CutLine"
	ActionDuplicateLine           = "DuplicateLine"
	ActionDeleteLine              = "DeleteLine"
	ActionMoveLinesUp             = "MoveLinesUp"
	ActionMoveLinesDown           = "MoveLinesDown"
	ActionIndentSelection         = "IndentSelection"
	ActionOutdentSelection        = "OutdentSelection"
	ActionOutdentLine             = "OutdentLine"
	ActionPaste                   = "Paste"
	ActionSelectAll               = "SelectAll"
	ActionStart                   = "Start"
	ActionEnd                     = "End"
	ActionPageUp                  = "PageUp"
	ActionPageDown                = "PageDown"
	ActionSelectPageUp            = "SelectPageUp"
	ActionSelectPageDown          = "SelectPageDown"
	ActionHalfPageUp              = "HalfPageUp"
	ActionHalfPageDown            = "HalfPageDown"
	ActionStartOfLine             = "StartOfLine"
	ActionEndOfLine               = "EndOfLine"
	ActionToggleRuler             = "ToggleRuler"
	ActionToggleOverwriteMode     = "ToggleOverwriteMode"
	ActionEscape                  = "Escape"
	ActionScrollUp                = "ScrollUp"
	ActionScrollDown              = "ScrollDown"
	ActionSpawnMultiCursor        = "SpawnMultiCursor"
	ActionSpawnMultiCursorSelect  = "SpawnMultiCursorSelect"
	ActionRemoveMultiCursor       = "RemoveMultiCursor"
	ActionRemoveAllMultiCursors   = "RemoveAllMultiCursors"
	ActionSkipMultiCursor         = "SkipMultiCursor"
	ActionJumpToMatchingBrace     = "JumpToMatchingBrace"
	ActionInsertEnter             = "InsertEnter"
	ActionUnbindKey               = "UnbindKey"
	ActionStartOfTextToggle       = "StartOfTextToggle"
	ActionMousePress              = "MousePress"
	ActionMouseDrag               = "MouseDrag"
	ActionMouseRelease            = "MouseRelease"
	ActionSetManualSelectionStart = "SetManualSelectionStart"
	ActionSetManualSelectionEnd   = "SetManualSelectionEnd"
	ActionToggleBookmark          = "ToggleBookmark"
)

type KeyDesc struct {
	KeyCode   tcell.Key
	Modifiers tcell.ModMask
	R         rune
}

// Utility function for parsing key sequences which have the same format as Smidgen keybindings
func ParseKeySequence(k string) (KeyDesc, bool) {
	kd, ok := action.ParseKeyboardSequence(k)
	if !ok {
		return KeyDesc{}, false
	}
	return KeyDesc{
		KeyCode:   kd.KeyCode,
		Modifiers: kd.Modifiers,
		R:         kd.R,
	}, true
}
