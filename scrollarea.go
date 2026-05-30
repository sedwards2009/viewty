package termtronic

type ScrollArea struct {
	*WidgetBase
	offsetX int
	offsetY int
	minWidth int
	minHeight int
	contentWidget Widget
}


func NewScrollArea() *ScrollArea {
	return &ScrollArea{
		WidgetBase: NewWidgetBase(),
	}
}

func (s *ScrollArea) SetMinimumSize(minWidth int, minHeight int) {
	s.minWidth = minWidth
	s.minHeight = minHeight
}

func (s *ScrollArea) SetContentWidget(contentWidget Widget) {
	s.contentWidget = contentWidget
}

func (s *ScrollArea) Reposition(x int, y int, width int, height int) {
	s.WidgetBase.Reposition(x, y, width, height)
	contentWidth := max(width, s.minWidth)
    contentHeight := max(height, s.minHeight)
    s.contentWidget.Reposition(0, 0, contentWidth, contentHeight)
}

func (s *ScrollArea) Render(screen TranslateScreenWriter) {
	s.contentWidget.Render(screen.NewTranslate(s.offsetX, s.offsetY))	// TODO: clip
}

func (s *ScrollArea) ChildWidgetAt(x int, y int) Widget {
	return s.contentWidget.ChildWidgetAt(x+s.offsetX, y+s.offsetY)
}

func (s *ScrollArea) SetOffsetX(offset int) {
	s.offsetX = offset
}

func (s *ScrollArea) SetOffsetY(offset int) {
	s.offsetY = offset
}

func (s *ScrollArea) OffsetX() int {
	return s.offsetX
}

func (s *ScrollArea) OffsetY() int {
	return s.offsetY
}
