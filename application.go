package termtronic

import (
	"io"
	"log"
	"os"

	tcell "github.com/gdamore/tcell/v2"
)

type Application struct {
	screen        tcell.Screen
	rootWidget    Widget
	enableLogging bool

	focusWidget Widget
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) EnableLogging(on bool) {
	a.enableLogging = on
}

func (a *Application) SetRootWidget(widget Widget) {
	a.rootWidget = widget
}

func (a *Application) Run() {
	if a.rootWidget == nil {
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

	quit := func() {
		a.screen.Fini()
		os.Exit(0)
	}

	var logFile *os.File
	if a.enableLogging {
		logFile = a.setupLogging()
		defer logFile.Close()
		log.Println("Dinky starting with logging enabled")
	} else {
		// Disable logging by setting output to discard
		log.SetOutput(io.Discard)
	}

	for {
		// Update screen
		a.screen.Show()

		ev := a.screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			width, height := ev.Size()
			log.Printf("Resize(%d, %d)\n", width, height)
			a.rootWidget.Reposition(0, 0, width, height)
			a.rootWidget.Render(NewTranslateScreenWriterAdapter(a.screen))
			a.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		case *tcell.EventMouse:
			a.handleMouseEvent(ev)
			a.rootWidget.Render(NewTranslateScreenWriterAdapter(a.screen))
		}
	}
}

func (a *Application) handleMouseEvent(ev *tcell.EventMouse) {
	x, y := ev.Position()
	hitWidget := a.rootWidget
	childHitWidget := hitWidget.ChildWidgetAt(x, y)
	if childHitWidget != nil {
		hitWidget = childHitWidget
	}

	if hitWidget == nil {
		return
	}

    // Find the complete path from the root widget to this widget.
	widgetPath := []Widget{}

    ptr := hitWidget.Parent()
    for ptr != nil {
    	widgetPath = append(widgetPath, ptr)
     	ptr = ptr.Parent();
    }

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

func (a *Application) setupLogging() *os.File {
	logFile, err := os.OpenFile("termtronic.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	log.SetOutput(logFile)
	return logFile
}
