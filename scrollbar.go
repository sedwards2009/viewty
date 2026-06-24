package viewty

type Scrollbar struct {
	*Flex
	track       *ScrollbarTrack
	upBtn       *Button
	downBtn     *Button
	changedFunc func(position int)
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
		if scrollbar.changedFunc != nil {
			scrollbar.changedFunc(newPos)
		}
	})

	downBtn.SetOnClick(func(id string) {
		pos := scrollbar.track.ThumbPosition() + max(scrollbar.track.ThumbSize()/2, 1)
		newPos := scrollbar.track.SetThumbPosition(pos)
		if scrollbar.changedFunc != nil {
			scrollbar.changedFunc(newPos)
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

func (s *Scrollbar) SetChangedFunc(changedFunc func(position int)) {
	s.changedFunc = changedFunc
	s.track.SetChangedFunc(changedFunc)
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
