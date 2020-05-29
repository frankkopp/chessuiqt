package main

import (
	"os"

	"github.com/therecipe/qt/widgets"

	"github.com/frankkopp/chessuiqt/internal/views"
)

func main() {

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetStyle2("fusion")

	// create a window
	window := views.NewMainWindow(app)

	// make the window visible
	window.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}

