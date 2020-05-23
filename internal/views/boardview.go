package views

import (
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// colors
var (
	border     = gui.NewQPen3(gui.NewQColor6("black"))
	brushWhite = gui.NewQBrush()
	brushBlack = gui.NewQBrush()
	font       = gui.NewQFont5(gui.NewQFont())
)

func NewBoardView(widget *widgets.QWidget) *widgets.QGraphicsView {

	// setup colors
	brushWhite.SetColor(gui.NewQColor6("white"))
	brushWhite.SetStyle(core.Qt__SolidPattern)
	brushBlack.SetColor(gui.NewQColor6("darkGray"))
	brushBlack.SetStyle(core.Qt__SolidPattern)

	// pane for chess board
	boardView := widgets.NewQGraphicsView(widget)
	// boardView.SetStyleSheet("background: yellow")
	boardView.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	boardView.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)

	// scene to draw the chess board on
	boardScene := widgets.NewQGraphicsScene(boardView)
	boardView.SetScene(boardScene)

	// redraw the chess board on resize
	boardView.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		drawBoard(boardScene, event.Size().Width(), event.Size().Height())
		boardView.SetFixedHeight(event.Size().Width() + 2)
	})

	boardView.Show()

	return boardView
}

func drawBoard(boardScene *widgets.QGraphicsScene, width int, height int) {

	boardSize := float64(width) - 1

	// clear boardScene and add all children again
	boardScene.Clear()

	// border around board
	boardScene.AddRect2(0, 0, boardSize, boardSize, border, brushWhite)

	// squares size
	squareSize := boardSize / 8
	font.SetPixelSize(int(squareSize * 0.2))

	// squares
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			if (rank+file)%2 != 0 {
				boardScene.AddRect2(float64(file)*squareSize, float64(rank)*squareSize, squareSize, squareSize, border, brushBlack)
			}
			if file == 0 {
				t := boardScene.AddSimpleText(strconv.Itoa(8-rank), font)
				t.SetX((float64(file) * squareSize) + (squareSize * 0.1))
				t.SetY((float64(rank) * squareSize) + (squareSize * 0.1))
			}
			if rank == 7 {
				t := boardScene.AddSimpleText(string('a'+file), font)
				t.SetX((float64(file) * squareSize) + (squareSize * 0.75))
				t.SetY((float64(rank) * squareSize) + (squareSize * 0.75))
			}
		}
	}

	// square numbering

	// pieces according to fen

}
