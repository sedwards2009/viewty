package viewty

import (
	"github.com/unilibs/uniwidth"
)

type TableHeader struct {
	label     string
	alignment int
}

type TableContents struct {
}

type TableContentAdapter interface {
	HeightInCells() int
	WidthInCells() int
	GetTextAt(x int, y int) string
	GetClassAt(x int, y int) string
}

type SelectionType int

const (
	SELECTION_NONE SelectionType = iota
	SELECTION_VERTICAL
	SELECTION_HORIZONTAL
	SELECTION_CELL
)

type Table struct {
	*WidgetBase
	selectionType SelectionType

	topHeaders  []TableHeader
	leftHeaders []TableHeader

	dirtyHeaderData bool
	topHeaderData   []headerData
	leftHeaderData  []headerData

	showTopHeaders  bool
	showLeftHeaders bool

	contentAdapter TableContentAdapter
}

type headerData struct {
	width int
}

func NewTable() *Table {
	return &Table{
		WidgetBase: NewWidgetBase(),
	}
}

func (t *Table) SetContent(contentAdapter TableContentAdapter) {
	t.contentAdapter = contentAdapter
}

func (t *Table) SetSelectionType(selectionType SelectionType) {
	t.selectionType = selectionType
}

func (t *Table) ShowTopHeaders() bool {
	return t.showTopHeaders
}

func (t *Table) SetShowTopHeaders(show bool) {
	t.showTopHeaders = show
	t.dirtyHeaderData = true
}

func (t *Table) ShowLeftHeaders() bool {
	return t.showLeftHeaders
}

func (t *Table) SetShowLeftHeaders(show bool) {
	t.showLeftHeaders = show
	t.dirtyHeaderData = true
}

func (t *Table) computerHeaderData() {
	if !t.dirtyHeaderData {
		return
	}

	if t.showTopHeaders && len(t.topHeaders) > 0 {
		t.topHeaderData = make([]headerData, len(t.topHeaders))
		for i := range t.topHeaders {
			maxWidth := 0
			for j := 0; j < t.contentAdapter.HeightInCells(); j++ {
				text := t.contentAdapter.GetTextAt(i, j)
				w := uniwidth.StringWidth(text)
				if w > maxWidth {
					maxWidth = w
				}
			}
			t.topHeaderData[i] = headerData{width: maxWidth}
		}
	}

	if t.showLeftHeaders && len(t.leftHeaders) > 0 {
		t.leftHeaderData = make([]headerData, len(t.leftHeaders))
		for i := range t.leftHeaders {
			maxWidth := 0
			for j := 0; j < t.contentAdapter.WidthInCells(); j++ {
				text := t.contentAdapter.GetTextAt(j, i)
				w := uniwidth.StringWidth(text)
				if w > maxWidth {
					maxWidth = w
				}
			}
			t.leftHeaderData[i] = headerData{width: maxWidth}
		}
	}
}

func (t *Table) Render(painter Painter) {
	t.computerHeaderData()

	_, _, w, h := t.Position()

	style := t.GetStyle("Table", nil)
	if t.showTopHeaders || t.showLeftHeaders {
		headerStyle := GetTCellStyle(style, "headerForegroundColor", "headerBackgroundColor")
		if t.showTopHeaders && len(t.topHeaders) > 0 {
			for j := range h {
				x := 0
				for i, header := range t.topHeaders {
					labelWidth := LabelWidth(header.label)
					if t.topHeaderData != nil {
						labelWidth = max(labelWidth, t.topHeaderData[i].width)
					}
					PrintString(painter, x, j, headerStyle, header.label)
					x += labelWidth
				}
			}
		}

		if t.showLeftHeaders && len(t.leftHeaders) > 0 {
			for j := range h {
				x := 0
				for i, header := range t.leftHeaders {
					labelWidth := LabelWidth(header.label)
					if t.leftHeaderData != nil {
						labelWidth = max(labelWidth, t.leftHeaderData[i].width)
					}
					PrintString(painter, x, j, headerStyle, header.label)
					x += labelWidth
				}
			}
		}
	}

	rowCount := t.contentAdapter.HeightInCells()
	columnCount := t.contentAdapter.WidthInCells()

	cellStyle := GetTCellStyle(style, "foregroundColor", "backgroundColor")
	for y := 0; y < rowCount && y < h; y++ {
		for x := 0; x < columnCount && x < w; x++ {
			text := t.contentAdapter.GetTextAt(x, y)
			// class := t.contentAdapter.GetClassAt(x, y)
			PrintString(painter, x, y, cellStyle, text)
		}
	}
}
