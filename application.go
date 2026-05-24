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
}

func NewApplication() *Application {
	rootWidget := NewBox()

	var White = tcell.NewHexColor(0xf3f3f3).TrueColor()
    var Blue = tcell.NewHexColor(0x007ace).TrueColor()
	rootWidget.SetBackgroundStyle(tcell.StyleDefault.Foreground(White).Background(Blue))
	return &Application{
		rootWidget: rootWidget,
	}
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
			a.rootWidget.Render(a.screen)
			a.screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
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
