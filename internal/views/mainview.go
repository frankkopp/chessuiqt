package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func NewMainWindow() *widgets.QMainWindow {
	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(320, 320)
	window.Resize2(640, 640)
	window.SetWindowTitle("Chess UI")

	widget := widgets.NewQWidget(window, 0)
	qvBoxLayout := widgets.NewQVBoxLayout2(widget)
	widget.SetLayout(qvBoxLayout)
	widget.SetStyleSheet("background: red")

	boardView := NewBoardView(widget)

	// scrollArea := widgets.NewQScrollArea(widget)
	// scrollArea.SetSizePolicy2(widgets.QSizePolicy__Preferred, widgets.QSizePolicy__Preferred )
	// scrollArea.SetBackgroundRole(gui.QPalette__Highlight)
	// scrollArea.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	// scrollArea.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	// scrollArea.SetWidget(boardView)

	// add board view to window
	widget.Layout().AddWidget(boardView.View())
	widget.Layout().SetAlignment(boardView.View(), core.Qt__AlignTop)

	window.SetCentralWidget(widget)
	return window
}
