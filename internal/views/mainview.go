package views

import (
	"runtime"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var (
	app       *widgets.QApplication
	statusbar *widgets.QStatusBar
)

func NewMainWindow(qapp *widgets.QApplication) *widgets.QMainWindow {
	// make a application back reference available
	app = qapp

	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(320, 320)
	window.Resize2(640, 720)
	window.SetWindowTitle("Chess UI")

	// Statusbar
	statusbar = widgets.NewQStatusBar(window)
	window.SetStatusBar(statusbar)
	statusbar.ShowMessage(core.QCoreApplication_ApplicationDirPath(), 0)

	widget := widgets.NewQWidget(window, 0)
	qvBoxLayout := widgets.NewQVBoxLayout2(widget)
	widget.SetLayout(qvBoxLayout)
	widget.SetStyleSheet("background: red")

	boardPane := NewBoardView(widget)

	// scrollArea := widgets.NewQScrollArea(widget)
	// scrollArea.SetSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Preferred )
	// scrollArea.SetBackgroundRole(gui.QPalette__Highlight)
	// scrollArea.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	// scrollArea.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	// scrollArea.SetWidget(boardPane)

	// add board view to window
	widget.Layout().AddWidget(boardPane.View())
	widget.Layout().SetAlignment(boardPane.View(), core.Qt__AlignTop)

	if runtime.GOOS == "darwin" {
		window.SetUnifiedTitleAndToolBarOnMac(false)
	}

	window.SetCentralWidget(widget)
	return window
}
