package views

import (
	"log"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/svg"
	"github.com/therecipe/qt/widgets"
)

const (
	pieceImagePath = "./assets/qml/"
)

// colors
var (
	border     = gui.NewQPen3(gui.NewQColor6("black"))
	brushWhite = gui.NewQBrush()
	brushBlack = gui.NewQBrush()
	font       = gui.NewQFont5(gui.NewQFont())
	pieces     = make([]*gui.QImage, 0)
	piecesSvg  = make([]*svg.QGraphicsSvgItem, 0)
)

func NewBoardView(widget *widgets.QWidget) *widgets.QGraphicsView {

	if len(pieces) == 0 {
		readPiecePixmaps()
	}

	// setup colors
	brushWhite.SetColor(gui.NewQColor6("white"))
	brushWhite.SetStyle(core.Qt__SolidPattern)
	brushBlack.SetColor(gui.NewQColor6("darkGray"))
	brushBlack.SetStyle(core.Qt__SolidPattern)

	// pane for chess board
	boardView := widgets.NewQGraphicsView(widget)
	boardView.SetStyleSheet("background: yellow")
	boardView.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	boardView.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)

	// scene to draw the chess board on
	boardScene := widgets.NewQGraphicsScene(boardView)
	boardView.SetScene(boardScene)
	boardView.SetMinimumSize2(int(boardScene.Width()), int(boardScene.Height()))

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
	if len(boardScene.Children()) > 0 {
		for _, svgItem := range piecesSvg {
			if svgItem != nil {
				boardScene.RemoveItem(svgItem)
			}
		}
		boardScene.Clear()
	}

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

	for i, qImage := range pieces {
		if qImage == nil {
			continue
		}
		scaledImage := qImage.SmoothScaled(80, 80)
		pixmap := gui.NewQPixmap().FromImage(scaledImage, 0)
		p := boardScene.AddPixmap(pixmap)
		p.SetX(float64(i * 50))
		p.SetY(20)
		p.Show()
	}

	for i, svgItem := range piecesSvg {
		if svgItem == nil {
			continue
		}
		svgItem.SetScale(1.5)
		svgItem.SetX(float64(i * 50))
		svgItem.SetY(150)
		boardScene.AddItem(svgItem)
	}
}

func readPiecePixmaps() {
	// (" KPNBRQ- kpnbrq-")
	piecesFileNames := []string{"", "wK", "wP", "wN", "wB", "wR", "wQ", "", "", "bK", "bP", "bN", "bB", "bR", "bQ", ""}
	for _, piece := range piecesFileNames {
		if piece == "" {
			pieces = append(pieces, nil)
			piecesSvg = append(piecesSvg, nil)
			continue
		}
		filePath := pieceImagePath + piece + ".svg"
		image := gui.NewQImage9(filePath, "")
		imageSvg := svg.NewQGraphicsSvgItem2(filePath, nil)
		if image.IsNull() {
			log.Panicf("Can't read piece image from %s", filePath)
		}
		pieces = append(pieces, image)
		piecesSvg = append(piecesSvg, imageSvg)
	}
}
