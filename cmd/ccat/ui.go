package main

import (
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app      *tview.Application
	textView *tview.TextView

	setupDone = make(chan struct{})
	uiDone    = make(chan struct{})
)

func uiSetup(title string, w *io.Writer) {
	*argHuman = true
	app = tview.NewApplication()
	textView = tview.NewTextView().SetDynamicColors(true)
	textView.SetWrap(false).SetBorder(true).SetTitle(title)
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' || event.Rune() == 'Q' {
			app.Stop()
			return nil
		}
		return event
	})
	*w = tview.ANSIWriter(textView)
	go func() {
		if err := app.SetRoot(textView, true).Run(); err != nil {
			panic(err)
		}
		close(uiDone)
	}()
	close(setupDone)
}

func uiRun() {
	<-setupDone
	textView.ScrollToBeginning()
	app.Draw()
	<-uiDone
	setupDone = make(chan struct{})
	uiDone = make(chan struct{})
}

func uiWrapProcessFile(process func(_ io.Writer, path string)) func(io.Writer, string) {
	return func(_ io.Writer, path string) {
		var outWriter io.Writer
		uiSetup(path, &outWriter)
		process(outWriter, path)
		uiRun()
	}
}
