package viewty


type WidgetBase struct {
	parent    Widget
	name      string
	x         int
	y         int
	width     int
	height    int
	visible   bool
	styleFunc StyleFunc
}

func NewWidgetBase() *WidgetBase {
	return &WidgetBase{visible: true}
}

func (w *WidgetBase) SetParent(parent Widget) {
	w.parent = parent
}

func (w *WidgetBase) Parent() Widget {
	return w.parent
}

func (w *WidgetBase) Name() string {
	return w.name
}

func (w *WidgetBase) SetName(name string) {
	w.name = name
}

func (w *WidgetBase) Reposition(x, y, width, height int) {
	w.x = x
	w.y = y
	w.width = width
	w.height = height
}

func (w *WidgetBase) Position() (x int, y int, width int, height int) {
	return w.x, w.y, w.width, w.height
}

func (w *WidgetBase) ContainsPoint(x int, y int) bool {
	return x >= w.x && x < w.x+w.width && y >= w.y && y < w.y+w.height
}

func (w *WidgetBase) PointToAbs(x int, y int) (ax int, ay int) {
	if w.Parent() == nil {
		return x, y
	}
	return w.Parent().PointToAbs(x+w.x, y+w.y)
}

func (w *WidgetBase) ChildWidgetAt(x int, y int) Widget {
	return nil
}

func (w *WidgetBase) Render(screen Painter) {
}

func (w *WidgetBase) HandleMouseEvent(mouseEvent MouseEvent) bool {
	return false
}

func (w *WidgetBase) HandleKeyEvent(keyEvent KeyEvent) bool {
    return false
}

func (w *WidgetBase) SetVisible(visible bool) {
	w.visible = visible
}

func (w *WidgetBase) IsVisible() bool {
	return w.visible
}

func (w *WidgetBase) GetStyle(widgetType string, class []string) StyleMap {
	parent := w.Parent()
	var parentStyle StyleMap
	if parent == nil {
		parentStyle = make(map[string]any)
	} else {
		parentStyle = parent.GetStyle(widgetType, class)
	}

	if w.styleFunc == nil {
		return parentStyle
	}
	return w.styleFunc(parentStyle, widgetType, class)
}

func (w *WidgetBase) SetStyleFunc(styleFunc StyleFunc) {
	w.styleFunc = styleFunc
}
