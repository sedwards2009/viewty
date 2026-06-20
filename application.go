package viewty

import (
	"io"
	"log"
	"os"

	tcell "github.com/gdamore/tcell/v2"
)

type Application struct {
	screen        tcell.Screen
	layers        []Widget
	enableLogging bool
	forceRenderChannel chan bool
	focusWidget Widget
	logFile *os.File
}

var app *Application
var defaultStyleFunc StyleFunc

func NewApplication() *Application {
	app = &Application{
        forceRenderChannel: make(chan bool, 100),
        enableLogging: true,
	}
	// Load the default style
	builder := NewStyleBuilder()
	if err := builder.LoadJSON(defaultStyleJSON); err != nil {
		log.Printf("Error while loading default styles: %v\n", err)
	}
	defaultStyleFunc, _ = builder.Build()

	if app.enableLogging {
		app.logFile = app.setupLogging()
		log.Println("Dinky starting with logging enabled")
	} else {
		// Disable logging by setting output to discard
		log.SetOutput(io.Discard)
	}

	return app
}

func (a *Application) EnableLogging(on bool) {
	a.enableLogging = on
}

func (a *Application) ForceRender() {
	a.forceRenderChannel <- true
}

func (a *Application) AddLayerWidget(widget Widget) {
	a.layers = append(a.layers, widget)
	widget.SetStyleFunc(defaultStyleFunc)
	if a.screen != nil {
    	width, height := a.screen.Size()
    	widget.Reposition(0, 0, width, height)
	}
}

func (a *Application) RemoveLayerWidget(widget Widget) {
	widget.SetStyleFunc(nil)
	for i, w := range a.layers {
		if w == widget {
			a.layers = append(a.layers[:i], a.layers[i+1:]...)
			if a.IsWidgetOnFocusPath(widget) {
				a.focusWidget = nil
			}
			return
		}
	}
}

func (a *Application) Focus(widget Widget) {
	a.focusWidget = widget
}

func (a *Application) HasFocus(widget Widget) bool {
    return a.focusWidget == widget
}

func (a *Application) IsWidgetOnFocusPath(widget Widget) bool {
	if a.focusWidget == nil {
		return false
	}
	currentWidget := a.focusWidget
	for currentWidget != nil {
		if widget == currentWidget {
			return true
		}
		currentWidget = currentWidget.Parent()
	}
	return false
}

func (a *Application) Run() {
	if len(a.layers) == 0 {
		log.Fatalf("No root widget set on the application")
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	a.screen = screen
	a.screen.Init()

	a.screen.EnableMouse()
	a.screen.Clear()

	tcellEvents := make(chan tcell.Event)
	quitChannel := make(chan struct{}, 1)

	go func() {
	    a.screen.ChannelEvents(tcellEvents, quitChannel)
	}()

	a.screen.Show()
	a.rerender()
	for {
    	a.screen.Show()
        select {
            case ev := <-tcellEvents:
          		// Process event
          		switch ev := ev.(type) {
          		case *tcell.EventResize:
         			width, height := ev.Size()
         			log.Printf("Resize(%d, %d)\n", width, height)
                    for _, w := range a.layers {
                        w.Reposition(0, 0, width, height)
                    }

         			a.rerender()
         			a.screen.Sync()

          		case *tcell.EventKey:
         			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
                        quitChannel <- struct{}{}
                        a.screen.Fini()
                        os.Exit(0)
         			} else {
         			    a.handleKeyEvent(ev)
            			a.rerender()
         			}

          		case *tcell.EventMouse:
         			a.handleMouseEvent(ev)
         			a.rerender()
          		}

            case <-a.forceRenderChannel:
     			a.rerender()
        }
	}
}

func (a *Application) rerender() {
	a.screen.Clear()
	for _, w := range a.layers {
	    w.Render(NewPainter(a.screen))
	}
}

func (a *Application) findWidgetPath(hitWidget Widget) []Widget {
	widgetPath := []Widget{}

	ptr := hitWidget.Parent()
	for ptr != nil {
		widgetPath = append(widgetPath, ptr)
		ptr = ptr.Parent()
	}
	return widgetPath
}

func (a *Application) handleMouseEvent(ev *tcell.EventMouse) {
	x, y := ev.Position()
	hitWidget := a.layers[len(a.layers)-1]
	childHitWidget := hitWidget.ChildWidgetAt(x, y)
	if childHitWidget != nil {
		hitWidget = childHitWidget
	}

	if hitWidget == nil {
		return
	}

	widgetPath := a.findWidgetPath(hitWidget)

    mouseEvent := mouseEventImpl{
    	targetWidget: hitWidget,
        sourceEvent: ev,
    }
    var mouseEventInter MouseEvent = &mouseEvent

    // Perform a DOM event style 'capture' phase on each widget on the path.
    mouseEvent.phase = EVENT_PHASE_CAPTURE
    for i, _ := range widgetPath {
    	currentTargetWidget := widgetPath[len(widgetPath) -1 -i]
     	offsetX, offsetY := currentTargetWidget.PointToAbs(0, 0)
        mouseEvent.x = x - offsetX
        mouseEvent.y = y - offsetY
        if currentTargetWidget.HandleMouseEvent(mouseEventInter) {
        	return	// cancelled
        }
    }

    mouseEvent.phase = EVENT_PHASE_TARGET
   	offsetX, offsetY := hitWidget.PointToAbs(0, 0)
    mouseEvent.x = x - offsetX
    mouseEvent.y = y - offsetY
    if hitWidget.HandleMouseEvent(mouseEventInter) {
    	return // cancelled
    }

    // Now the bubble phase
    mouseEvent.phase = EVENT_PHASE_BUBBLE
    for _, currentTargetWidget := range widgetPath {
	   	offsetX, offsetY := currentTargetWidget.PointToAbs(0, 0)
		mouseEvent.x = x - offsetX
		mouseEvent.y = y - offsetY
        if currentTargetWidget.HandleMouseEvent(mouseEventInter) {
        	return // cancelled
        }
    }
}

func (a *Application) handleKeyEvent(ev *tcell.EventKey) {
	hitWidget := a.focusWidget
	if hitWidget == nil {
	    return
	}

	widgetPath := a.findWidgetPath(hitWidget)

    keyEvent := keyEventImpl{
     	targetWidget: hitWidget,
        	sourceEvent: ev,
    }
    var keyEventInter KeyEvent = &keyEvent

    // Perform a DOM event style 'capture' phase on each widget on the path.
    keyEvent.phase = EVENT_PHASE_CAPTURE
    for i, _ := range widgetPath {
     	currentTargetWidget := widgetPath[len(widgetPath) -1 -i]
        if currentTargetWidget.HandleKeyEvent(keyEventInter) {
        	return	// cancelled
        }
    }

    keyEvent.phase = EVENT_PHASE_TARGET
    if hitWidget.HandleKeyEvent(keyEventInter) {
     	return // cancelled
    }

    // Now the bubble phase
    keyEvent.phase = EVENT_PHASE_BUBBLE
    for _, currentTargetWidget := range widgetPath {
        if currentTargetWidget.HandleKeyEvent(keyEventInter) {
        	return // cancelled
        }
    }
}

func (a *Application) setupLogging() *os.File {
	logFile, err := os.OpenFile("viewty.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}
