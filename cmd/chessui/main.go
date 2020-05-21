package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func main() {

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(320, 320)
	window.SetWindowTitle("Chess UI")
	window.SetLayout(widgets.NewQVBoxLayout())

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	qvBoxLayout := widgets.NewQVBoxLayout()
	widget.SetLayout(qvBoxLayout)
	window.SetCentralWidget(widget)

	boardView := board(widget)

	// add board view to window
	widget.Layout().AddWidget(boardView)
	widget.Layout().SetAlignment(boardView, core.Qt__AlignTop)

	// make the window visible
	window.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}

func board(widget *widgets.QWidget) *widgets.QGraphicsView {
	// pane for chess board
	boardView := widgets.NewQGraphicsView(widget)

	// colors
	border := gui.NewQPen3(gui.NewQColor6("black"))
	brushBackground := gui.NewQBrush()
	brushBackground.SetColor(gui.NewQColor6("gray"))
	brushBackground.SetStyle(core.Qt__SolidPattern)
	brushWhite := gui.NewQBrush()
	brushWhite.SetColor(gui.NewQColor6("lightGray"))
	brushWhite.SetStyle(core.Qt__SolidPattern)
	brushBlack := gui.NewQBrush()
	brushBlack.SetColor(gui.NewQColor6("darkGray"))
	brushBlack.SetStyle(core.Qt__SolidPattern)

	boardView.SetScene(drawBoard(boardView, 320, 320, border, brushBackground, brushBlack, brushWhite))

	boardView.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		boardView.SetScene(drawBoard(boardView, event.Size().Width(), event.Size().Height(), border, brushBackground, brushBlack, brushWhite))
	})

	// boardView.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
	// 	fmt.Println("MouseMove:", event.LocalPos().X(), event.LocalPos().Y())
	// })
	// boardView.ConnectResizeEvent(func(event *gui.QResizeEvent) {
	// 	fmt.Println("Resize:", event.Size().Height(), event.Size().Width(), )
	// })
	// boardView.ConnectPaintEvent(func(event *gui.QPaintEvent) {
	// 	fmt.Println("Paint:", event.Type())
	// })

	boardView.Show()

	return boardView
}

func drawBoard(boardView *widgets.QGraphicsView, width int, height int, border *gui.QPen, brushBackground *gui.QBrush, brushBlack *gui.QBrush, brushWhite *gui.QBrush) *widgets.QGraphicsScene {
	boardView.SetMinimumSize2(width, width)
	boardSize := float64(width)

	// board
	boardScene := widgets.NewQGraphicsScene(boardView)
	boardScene.AddRect2(0, 0, boardSize, boardSize, border, brushBackground)

	// squares
	squareSize := boardSize / 8
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			sqColor := brushBlack
			if (rank+file)%2 == 0 {
				sqColor = brushWhite
			}
			boardScene.AddRect2(float64(file)*squareSize, float64(rank)*squareSize, squareSize, squareSize, border, sqColor)
		}
	}

	return boardScene
}
