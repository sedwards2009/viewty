package viewty

import (
	tcell "github.com/gdamore/tcell/v2"
)

type ScrollbarTrack struct {
	*WidgetBase
	thumbPosition   int
	thumbSize       int
	max             int
	width           int
	isHorizontal    bool
	thin            bool
	beforeDrawFunc  func(screen tcell.Screen)
	changedFunc     func(position int)
	cover           *ScrollbarTrackCover
}

func NewScrollbarTrack() *ScrollbarTrack {
	result := &ScrollbarTrack{
		WidgetBase: NewWidgetBase(),
		thumbPosition: 0,
		thumbSize:     10,
		max:           100,
		width:         1,
	}
	result.cover = newScrollbarTrackCover(result)
	return result
}

func (s *ScrollbarTrack) SetBeforeDrawFunc(beforeDrawFunc func(screen tcell.Screen)) {
	s.beforeDrawFunc = beforeDrawFunc
}

func (s *ScrollbarTrack) SetHorizontal(isHorizontal bool) {
	s.isHorizontal = isHorizontal
}

func (s *ScrollbarTrack) SetThin(thin bool) {
	s.thin = thin
}

func (s *ScrollbarTrack) SetChangedFunc(changedFunc func(position int)) {
	s.changedFunc = changedFunc
}

func (s *ScrollbarTrack) SetThumbPosition(position int) int {
	if position < 0 {
		s.thumbPosition = 0
	} else if position > s.max-s.thumbSize {
		s.thumbPosition = s.max - s.thumbSize
	} else {
		s.thumbPosition = position
	}
	return s.thumbPosition
}

func (s *ScrollbarTrack) ThumbPosition() int {
	return s.thumbPosition
}

func (s *ScrollbarTrack) SetThumbSize(size int) {
	if size < 1 {
		s.thumbSize = 1
	} else {
		s.thumbSize = size
	}
}

func (s *ScrollbarTrack) ThumbSize() int {
	return s.thumbSize
}

func (s *ScrollbarTrack) SetMax(max int) {
	if max > 0 {
		s.max = max
	}
}

func (s *ScrollbarTrack) Max() int {
	return s.max
}

func (s *ScrollbarTrack) SetWidth(width int) {
	if width > 0 {
		s.width = width
	}
}

func (s *ScrollbarTrack) Width() int {
	return s.width
}

func (s *ScrollbarTrack) Render(painter Painter) {
	if s.beforeDrawFunc != nil {
		s.beforeDrawFunc(painter.Screen())
	}

	_, _, width, height := s.Position()
	if width < 1 || height < 1 {
		return
	}
	x := 0
	y := 0

	if s.isHorizontal && s.thin {
		s.drawThinHorizontal(painter, x, y, width)
		return
	}

	styles := s.GetStyle("ScrollbarTrack", []string{})
	trackStyle := GetTCellStyle(styles, "trackColor", "trackColor")
	thumbStyle := GetTCellStyle(styles, "trackColor", "thumbColor")

	firstHalfCellRune := '\u2584'
	secondHalfCellRune := '\u2580'

	majorLength := height
	majorPos := y
	if s.isHorizontal {
		majorLength = width
		majorPos = x
		firstHalfCellRune = '\u2590'
		secondHalfCellRune = '\u258c'
	}

	setContent := func(n int, ch rune, style tcell.Style) {
		if s.isHorizontal {
			painter.SetContent(n, y, ch, nil, style)
		} else {
			painter.SetContent(x, n, ch, nil, style)
		}
	}

	// Draw the track
	for i := 0; i < majorLength; i++ {
		setContent(majorPos+i, ' ', trackStyle)
	}

	position := s.thumbPosition
	thumbSize := s.thumbSize
	if thumbSize > s.max {
		thumbSize = s.max
		position = 0
	}

	doubleMajorLength := majorLength * 2
	doubleThumbSizeFloat := float64(doubleMajorLength) * float64(thumbSize) / float64(s.max)
	doubleThumbSize := int(doubleThumbSizeFloat + 0.5) // Round to nearest integer
	if doubleThumbSize < 1 {
		doubleThumbSize = 1
	}

	doubleThumbMajor := doubleMajorLength * position / s.max
	thumbReverseStyle := thumbStyle.Reverse(true)

	if position == s.max-thumbSize {
		// Special case for when the thumb is at the bottom and we don't want to show a gap due to rounding
		doubleThumbSize += 2
	}

	if doubleThumbMajor&1 == 1 {
		// Draw the top of the thumb in the bottom half of the cell
		setContent(majorPos+doubleThumbMajor>>1, firstHalfCellRune, thumbReverseStyle)
		doubleThumbMajor++
		doubleThumbSize--
	}

	if doubleThumbSize&1 == 1 {
		// Draw the bottom part of the thumb in the top half of the cell
		setContent(majorPos+doubleThumbMajor>>1+doubleThumbSize>>1, secondHalfCellRune, thumbReverseStyle)
		doubleThumbSize--
	}

	// Draw the scrollbar thumb.
	thumbPos := doubleThumbMajor >> 1 // Convert back to single height
	for i := 0; i < doubleThumbSize>>1; i++ {
		if thumbPos+i < majorLength {
			setContent(majorPos+thumbPos+i, ' ', thumbStyle)
		}
	}
}

func (s *ScrollbarTrack) HandleMouseEvent(mouseEvent MouseEvent) bool {
	event := mouseEvent.SourceEvent()
	buttons := event.Buttons()

	// Handle left button press (mousedown)
	if buttons == tcell.Button1 {
		relativeX, relativeY := mouseEvent.Position()

		if !s.handleLeftMouse(relativeX, relativeY) {
			return false
		}
		app.AddLayerWidget(s.cover)
		return true
	}

	// Handle scroll wheel up
	if buttons == tcell.WheelUp {
		pos := s.thumbPosition - max(s.thumbSize/2, 1)
		newPos := s.SetThumbPosition(pos)
		if s.changedFunc != nil {
			s.changedFunc(newPos)
		}
		return true
	}
	// Handle scroll wheel down
	if buttons == tcell.WheelDown {
		pos := s.thumbPosition + max(s.thumbSize/2, 1)
		newPos := s.SetThumbPosition(pos)
		if s.changedFunc != nil {
			s.changedFunc(newPos)
		}
		return true
	}

	return false
}

func (s *ScrollbarTrack) handleLeftMouse(relativeX int, relativeY int) bool {
	_, _, width, height := s.Position()
	if s.isHorizontal {
		if relativeY < 0 || relativeY >= height || relativeX < 0 || relativeX >= width {
			return false
		}
		s.setPositionFromMajor(relativeX, width)
	}

	if !s.isHorizontal {
		if relativeX < 0 || relativeX >= width || relativeY < 0 || relativeY >= height {
			return false
		}
		s.setPositionFromMajor(relativeY, height)
	}
	return true
}

// setPositionFromMajor sets the thumb position from a coordinate along the
// scrollbar's major axis (relative to the track's inner rect), centering the
// thumb on the pointer and clamping to the valid range. The pointer may be
// outside the track during a drag; the clamping handles that. The changed
// callback is fired with the resulting position.
func (s *ScrollbarTrack) setPositionFromMajor(eventMajorAxis, majorLength int) {
	if majorLength < 1 {
		return
	}
	newPosition := eventMajorAxis*s.max/majorLength - s.thumbSize/2
	if newPosition < 0 {
		newPosition = 0
	} else if newPosition > s.max-s.thumbSize {
		newPosition = s.max - s.thumbSize
	}
	s.thumbPosition = newPosition
	if s.changedFunc != nil {
		s.changedFunc(newPosition)
	}
}

func (s *ScrollbarTrack) drawThinHorizontal(painter Painter, x, y, width int) {
	const upperHalfRune = '▀'

	styles := s.GetStyle("ScrollbarTrack", []string{})
	trackStyle := GetTCellStyle(styles, "trackColor", "backgroundColor")
	thumbStyle := GetTCellStyle(styles, "thumbColor", "backgroundColor")

	position := s.thumbPosition
	thumbSize := s.thumbSize
	if thumbSize > s.max {
		thumbSize = s.max
		position = 0
	}

	thumbCells := (width*thumbSize + s.max/2) / s.max
	if thumbCells < 1 {
		thumbCells = 1
	}
	if thumbCells > width {
		thumbCells = width
	}

	thumbStart := (width*position + s.max/2) / s.max
	if thumbStart > width-thumbCells {
		thumbStart = width - thumbCells
	}
	if thumbStart < 0 {
		thumbStart = 0
	}

	for i := 0; i < width; i++ {
		style := trackStyle
		if i >= thumbStart && i < thumbStart+thumbCells {
			style = thumbStyle
		}
		painter.SetContent(x+i, y, upperHalfRune, nil, style)
	}
}

type ScrollbarTrackCover struct {
	*WidgetBase
	track *ScrollbarTrack
}

func newScrollbarTrackCover(track *ScrollbarTrack) *ScrollbarTrackCover {
	return &ScrollbarTrackCover{
		WidgetBase: NewWidgetBase(),
		track: track,
	}
}

func (s *ScrollbarTrackCover) HandleMouseEvent(mouseEvent MouseEvent) bool {
	event := mouseEvent.SourceEvent()
	buttons := event.Buttons()
	if buttons == tcell.Button1 {
		trackTopX, trackTopY := s.track.PointToAbs(0, 0)
		x, y := mouseEvent.Position()
		relativeX := x - trackTopX
		relativeY := y - trackTopY
		if s.track.isHorizontal {
			relativeY = 0
		} else {
			relativeX = 0
		}
		s.track.handleLeftMouse(relativeX, relativeY)
	} else {
		app.RemoveLayerWidget(s)
	}
	return false
}
