package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
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

	scrollArea := widgets.NewQScrollArea(nil)
	scrollArea.SetBackgroundRole(gui.QPalette__Highlight)
	scrollArea.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	scrollArea.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	// scrollArea.SetWidget(widget)

	// add chess board view
	boardView := NewBoardView(widget)

	// add board view to window
	widget.Layout().AddWidget(boardView)
	widget.Layout().SetAlignment(boardView, core.Qt__AlignTop)

	window.SetCentralWidget(widget)
	return window
}
