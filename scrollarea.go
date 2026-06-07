package viewty

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
	s.contentWidget.SetParent(s)
}

func (s *ScrollArea) Reposition(x int, y int, width int, height int) {
	s.WidgetBase.Reposition(x, y, width, height)
	contentWidth := max(width, s.minWidth)
    contentHeight := max(height, s.minHeight)
    s.contentWidget.Reposition(0, 0, contentWidth, contentHeight)
}

func (s *ScrollArea) Render(painter Painter) {
	if s.contentWidget == nil {
		return
	}
	s.contentWidget.Render(painter.Translate(s.offsetX, s.offsetY))	// TODO: clip
}

func (s *ScrollArea) ChildWidgetAt(x int, y int) Widget {
	if s.contentWidget == nil {
		return nil
	}

	myX, myY, _, _ := s.Position()
	childX := x - myX - s.offsetX
	childY := y - myY - s.offsetY
	return s.contentWidget.ChildWidgetAt(childX, childY)
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

func (s *ScrollArea) PointToAbs(x int, y int) (ax int, ay int) {
	if s.Parent() == nil {
		return x, y
	}
	return s.Parent().PointToAbs(x+s.x+s.offsetX, y+s.y+s.offsetY)
}
