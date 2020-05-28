package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func NewMainWindow() *widgets.QMainWindow {
	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(320, 320)
	window.Resize2(640, 720)
	window.SetWindowTitle("Chess UI")

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

	window.SetCentralWidget(widget)
	return window
}
