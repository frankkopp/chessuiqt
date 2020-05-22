package views

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// colors
var (
	border     = gui.NewQPen3(gui.NewQColor6("black"))
	brushWhite = gui.NewQBrush()
	brushBlack = gui.NewQBrush()
)

func NewBoardView(widget *widgets.QWidget) *widgets.QGraphicsView {

	// setup colors
	brushWhite.SetColor(gui.NewQColor6("white"))
	brushWhite.SetStyle(core.Qt__SolidPattern)
	brushBlack.SetColor(gui.NewQColor6("darkGray"))
	brushBlack.SetStyle(core.Qt__SolidPattern)

	// pane for chess board
	boardView := widgets.NewQGraphicsView(widget)

	// scene to draw the chess board on
	boardScene := widgets.NewQGraphicsScene(boardView)
	boardView.SetScene(boardScene)

	// redraw the chess board on resize
	boardView.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		drawBoard(boardScene, event.Size().Width(), event.Size().Height())
		boardView.SetMinimumHeight(event.Size().Width())
	})
	// boardView.ConnectPaintEvent(func(event *gui.QPaintEvent) {
	// 	fmt.Println("Paint:", event.Type(), time.Now())
	// 	boardView.SetMinimumSize2(event.Size().Width()+4, event.Size().Width()+4)
	// 	drawBoard(boardScene, event.Size().Width(), event.Size().Height())
	// 	boardView.SetScene(boardScene)
	// })

	// boardView.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
	// 	fmt.Println("MouseMove:", event.LocalPos().X(), event.LocalPos().Y())
	// })
	// boardView.ConnectResizeEvent(func(event *gui.QResizeEvent) {
	// 	fmt.Println("Resize:", event.Size().Height(), event.Size().Width(), )
	// })

	boardView.Show()

	return boardView
}

func drawBoard(boardScene *widgets.QGraphicsScene, width int, height int) {

	boardSize := float64(width)

	// clear boardScene and add all children again
	boardScene.Clear()

	// border around board
	boardScene.AddRect2(0, 0, boardSize, boardSize, border, brushWhite)

	// squares
	squareSize := boardSize / 8
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			if (rank+file)%2 != 0 {
				boardScene.AddRect2(float64(file)*squareSize, float64(rank)*squareSize, squareSize, squareSize, border, brushBlack)
			}
		}
	}

	// square numbering

	// pieces according to fen

}
