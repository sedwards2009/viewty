package viewty

import (
	tcell "github.com/gdamore/tcell/v2"
)

type Painter interface {
	GetContent(x, y int) (primary rune, combining []rune, style tcell.Style, width int)

	// Set the cell contents at x,y
	//
	// Note: Clipping is first applied and then translation before writing the cell to the screen.
	SetContent(x int, y int, primary rune, combining []rune, style tcell.Style)
	Size() (width, height int)
    Translate(tX int, tY int) Painter
    ApplyClipArea(x int, y int, width int, height int) Painter
    IsVisible() bool
    Screen() tcell.Screen
}

type PainterImpl struct {
	screen     tcell.Screen
	clipX      int	// Clip is in terms of the final screen space
	clipY      int
	clipWidth  int
	clipHeight int
	translateX int
	translateY int
}

func NewPainter(screen tcell.Screen) Painter {
	width, height := screen.Size()
	return &PainterImpl{
		screen: screen,
		clipWidth: width,
		clipHeight: height,
	}
}

func (p *PainterImpl) GetContent(x, y int) (primary rune, combining []rune, style tcell.Style, width int) {
	xTranslated := x + p.translateX
	yTranslated := y + p.translateY
	return p.screen.GetContent(xTranslated, yTranslated)
}

func (p *PainterImpl) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	xTranslated := x + p.translateX
	yTranslated := y + p.translateY
	if xTranslated < p.clipX || xTranslated >= (p.clipX + p.clipWidth) {
		return
	}
	if yTranslated < p.clipY || yTranslated >= (p.clipY + p.clipHeight) {
		return
	}
	p.screen.SetContent(xTranslated, yTranslated, primary, combining, style)
}

func (p *PainterImpl) Size() (width, height int) {
	return p.screen.Size()
}

func (p *PainterImpl) IsVisible() bool {
	return p.clipWidth != 0 && p.clipHeight != 0
}

func (p *PainterImpl) Screen() tcell.Screen {
    return p.screen
}

func (p *PainterImpl) Translate(translateX int, translateY int) Painter {
	return &PainterImpl{
		screen:     p.screen,
		clipX:      p.clipX,
		clipY:      p.clipY,
		clipWidth:  p.clipWidth,
		clipHeight: p.clipHeight,
		translateX: translateX + p.translateX,
		translateY: translateY + p.translateY,
	}
}

func (p *PainterImpl) ApplyClipArea(cX int, cY int, width int, height int) Painter {
	newX, newY, newWidth, newHeight := intersectRectangles(
		p.clipX, p.clipY, p.clipWidth, p.clipHeight,
		cX + p.translateX, cY+p.translateY, width, height,
	)
	if newWidth < 0 || newHeight < 0 {
		newWidth = 0
		newHeight = 0
	}
	return &PainterImpl{
		screen:     p.screen,
		clipX:      newX,
		clipY:      newY,
		clipWidth:  newWidth,
		clipHeight: newHeight,
		translateX: p.translateX,
		translateY: p.translateY,
	}
}

// IntersectRectangles computes the intersection of two rectangles.
// Rectangles are defined by their top-left corner (x, y) and dimensions (width, height).
// Returns the intersection rectangle coordinates, or 0,0,0,0 if there is no intersection.
func intersectRectangles(r1X, r1Y, r1Width, r1Height, r2X, r2Y, r2Width, r2Height int) (x, y, width, height int) {
	x = max(r1X, r2X)
	y = max(r1Y, r2Y)
	right := min(r1X+r1Width, r2X+r2Width)
	bottom := min(r1Y+r1Height, r2Y+r2Height)
	width = right - x
	height = bottom - y
	return
}
