package viewty

type Scrollbar struct {
	*Flex
	track       *ScrollbarTrack
	upBtn       *Button
	downBtn     *Button
	onChange    func(position int)
}

func NewScrollbar() *Scrollbar {
	track := NewScrollbarTrack()
	upBtn := NewButton()
	upBtn.SetText("▲")	// TODO: Get these chars from the style
	downBtn := NewButton()
	downBtn.SetText("▼")

	flex := NewVFlex()
	flex.AddWidget(track, 0, 1)
	flex.AddWidget(upBtn, 1, 0)
	flex.AddWidget(downBtn, 1, 0)

	scrollbar := &Scrollbar{
		Flex:    flex,
		track:   track,
		upBtn:   upBtn,
		downBtn: downBtn,
	}

	upBtn.SetOnClick(func(id string) {
		pos := scrollbar.track.ThumbPosition() - max(scrollbar.track.ThumbSize()/2, 1)
		newPos := scrollbar.track.SetThumbPosition(pos)
		if scrollbar.onChange != nil {
			scrollbar.onChange(newPos)
		}
	})

	downBtn.SetOnClick(func(id string) {
		pos := scrollbar.track.ThumbPosition() + max(scrollbar.track.ThumbSize()/2, 1)
		newPos := scrollbar.track.SetThumbPosition(pos)
		if scrollbar.onChange != nil {
			scrollbar.onChange(newPos)
		}
	})

	return scrollbar
}

func (s *Scrollbar) SetHorizontal(isHorizontal bool) {
	s.Flex.vertical = !isHorizontal
	s.track.SetHorizontal(isHorizontal)
	if isHorizontal {
		s.upBtn.SetText("\u25c4")
		s.downBtn.SetText("\u25ba")
	} else {
		s.upBtn.SetText("\u25b2")
		s.downBtn.SetText("\u25bc")
	}
}

func (s *Scrollbar) SetOnChange(changedFunc func(position int)) {
	s.onChange = changedFunc
	s.track.SetOnChange(changedFunc)
}

func (s *Scrollbar) HandleMouseEvent(mouseEvent MouseEvent) bool {
	return s.Flex.HandleMouseEvent(mouseEvent)
}

func (s *Scrollbar) HandleKeyEvent(keyEvent KeyEvent) bool {
	return s.Flex.HandleKeyEvent(keyEvent)
}

func (s *Scrollbar) Render(painter Painter) {
	s.Flex.Render(painter)
}

func (s *Scrollbar) SetThumbPosition(position int) int {
	return s.track.SetThumbPosition(position)
}

func (s *Scrollbar) ThumbPosition() int {
	return s.track.ThumbPosition()
}

func (s *Scrollbar) SetThumbSize(size int) {
	s.track.SetThumbSize(size)
}

func (s *Scrollbar) ThumbSize() int {
	return s.track.ThumbSize()
}

func (s *Scrollbar) SetMax(max int) {
	s.track.SetMax(max)
}

func (s *Scrollbar) Max() int {
	return s.track.Max()
}
