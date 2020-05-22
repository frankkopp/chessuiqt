package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func NewMainWindow() *widgets.QMainWindow {
	// Create main window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(320, 320)
	window.SetWindowTitle("Chess UI")

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	qvBoxLayout := widgets.NewQVBoxLayout()
	widget.SetLayout(qvBoxLayout)
	window.SetCentralWidget(widget)

	// add chess board view
	boardView := NewBoardView(widget)

	// add board view to window
	widget.Layout().AddWidget(boardView)
	widget.Layout().SetAlignment(boardView, core.Qt__AlignTop)
	return window
}
