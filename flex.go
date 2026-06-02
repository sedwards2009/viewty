package termtronic


type flexItem struct {
	widget Widget
	fixed int
	proportion int
}

type Flex struct {
	*WidgetBase
	vertical bool
	gapSize int

	items []flexItem
}

func NewVFlex() *Flex {
	return &Flex{
		WidgetBase: NewWidgetBase(),
		vertical: true,
	}
}

func NewHFlex() *Flex {
	flex := NewVFlex()
	flex.vertical = false
	return flex
}

func (f *Flex) SetGapSize(gapSize int) {
	f.gapSize = gapSize
}

func (f *Flex) AddWidget(widget Widget, fixed int, proportion int) {
	item := flexItem{
		widget: widget,
		fixed: fixed,
		proportion: proportion,
	}

    f.items = append(f.items, item)
    widget.SetParent(f)
}

func (f *Flex) Reposition(x int, y int, width int, height int) {
	f.WidgetBase.Reposition(x, y, width, height)

	if f.vertical {
		yList, heightList := layout(height, f.gapSize, f.items)
		for i, item := range f.items {
			item.widget.Reposition(0, yList[i], width, heightList[i])
		}
	} else {
		xList, widthList := layout(width, f.gapSize, f.items)
		for i, item := range f.items {
			item.widget.Reposition(xList[i], 0, widthList[i], height)
		}
	}
}

func layout(size int, gapSize int, items []flexItem) ([]int, []int) {
	// Compute the amount of space left to distribute
	fixedSize := 0
	proportionDenominator := 0
	for _, item := range items {
		fixedSize += item.fixed
		proportionDenominator += item.proportion
	}

	remaining := size - fixedSize - (gapSize * (len(items) -1))

	xList := []int{}
	widthList := []int{}
	x := 0
	for _, item := range items {
		xList = append(xList, x)
		width := item.fixed + item.proportion * remaining / proportionDenominator
		widthList = append(widthList, width)

		x += width
		x += gapSize
	}
	return xList, widthList
}

func (f *Flex) Render(painter Painter) {
	for _, item := range f.items {
		x, y, width, height := item.widget.Position()
		clippedPainter := painter.Translate(x, y).ApplyClipArea(0, 0, width, height)
		if clippedPainter.IsVisible() {}
			item.widget.Render(clippedPainter)
		}
	}
}

func (f *Flex) ChildWidgetAt(x int, y int) Widget {
	myX, myY, _, _ := f.Position()
	childX := x - myX
	childY := y - myY
	for _, item := range f.items {
		if item.widget.ContainsPoint(childX, childY) {
			childWidget := item.widget.ChildWidgetAt(childX, childY)
			if childWidget != nil {
				return childWidget
			}
			return item.widget
		}
	}
	return nil
}
