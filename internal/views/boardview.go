package views

import (
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/svg"
	"github.com/therecipe/qt/widgets"

	"github.com/frankkopp/FrankyGo/pkg/position"
	"github.com/frankkopp/FrankyGo/pkg/types"
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
	piecesSvg  = make([]*svg.QSvgRenderer, 0)
)

type BoardView struct {
	boardView  *widgets.QGraphicsView
	boardScene *widgets.QGraphicsScene
	position   *position.Position
}

// NewBoardView Creates a board view and its data structure.
func NewBoardView(widget *widgets.QWidget) *BoardView {
	bv := &BoardView{}
	bv.newBoardView(widget)
	bv.position = position.NewPosition()
	return bv
}

// View returns the view of the board view
func (b *BoardView) View() *widgets.QGraphicsView {
	return b.boardView
}

func (b *BoardView) newBoardView(widget *widgets.QWidget) {

	// setup colors
	brushWhite.SetColor(gui.NewQColor6("white"))
	brushWhite.SetStyle(core.Qt__SolidPattern)
	brushBlack.SetColor(gui.NewQColor6("darkGray"))
	brushBlack.SetStyle(core.Qt__SolidPattern)

	// pane for chess board
	b.boardView = widgets.NewQGraphicsView(widget)
	b.boardView.SetStyleSheet("background: yellow")
	b.boardView.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	b.boardView.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)

	// scene to draw the chess board on
	b.boardScene = widgets.NewQGraphicsScene(b.boardView)
	b.boardView.SetScene(b.boardScene)
	b.boardView.SetMinimumSize2(int(b.boardScene.Width()), int(b.boardScene.Height()))

	if len(pieces) == 0 {
		b.readPiecePixmaps()
	}

	// redraw the chess board on resize
	b.boardView.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		b.drawBoard(event.Size().Width(), event.Size().Height())
		b.boardView.SetFixedHeight(event.Size().Width() + 2)
	})

	b.boardView.Show()
}

func (b *BoardView) drawBoard(width int, height int) {

	boardSize := float64(width) - 1

	// Clear boardScene and add all children again
	b.boardScene.Clear()

	// border around board
	b.boardScene.AddRect2(0, 0, boardSize, boardSize, border, brushWhite)

	// squares size
	squareSize := boardSize / 8
	fontSize := int(squareSize * 0.15)

	// squares
	for rank := 8; rank >= 1; rank-- {
		for file := 1; file <= 8; file++ {
			font.SetPixelSize(fontSize)
			// checkers
			if (rank+file)%2 == 0 {
				rect := b.boardScene.AddRect2(float64(file-1)*squareSize, float64(8-rank)*squareSize, squareSize, squareSize, border, brushBlack)
				rect.SetZValue(0.0)
			}
			// rank number
			if file == 1 {
				t := b.boardScene.AddSimpleText(strconv.Itoa(rank), font)
				t.SetX((float64(file-1) * squareSize) + (squareSize * 0.05))
				t.SetY((float64(8-rank) * squareSize) + (squareSize * 0.05))
			}
			// file letter
			if rank == 1 {
				t := b.boardScene.AddSimpleText(string('a'+file-1), font)
				t.SetX((float64(file-1) * squareSize) + (squareSize * 0.85))
				t.SetY((float64(8-rank) * squareSize) + (squareSize * 0.80))
			}
			// debug - square index numbers
			font.SetPixelSize(fontSize / 2)
			square := (rank-1)*8 + file - 1
			t := b.boardScene.AddSimpleText(strconv.Itoa(square), font)
			t.SetX((float64(file-1) * squareSize) + (squareSize * 0.8))
			t.SetY((float64(8-rank) * squareSize) + (squareSize * 0.05))

			// pieces for current position
			piece := b.position.GetPiece(types.Square(square))
			if piece != types.PieceNone {
				svgItem := svg.NewQGraphicsSvgItem(nil)
				svgRenderer := piecesSvg[piece]
				height := float64(svgRenderer.DefaultSize().Height())
				width := float64(svgRenderer.DefaultSize().Width())
				svgItem.SetSharedRenderer(svgRenderer)

				scaleFactor := squareSize / height * (0.9)
				svgItem.SetScale(scaleFactor)

				offsetX := height * 0.5 * scaleFactor
				offsetY := width * 0.5 * scaleFactor
				svgItem.SetX((float64(file-1) * squareSize) + (squareSize * 0.5) - offsetX)
				svgItem.SetY((float64(8-rank) * squareSize) + (squareSize * 0.5) - offsetY)
				svgItem.SetZValue(1.0)
				svgItem.Show()
				b.boardScene.AddItem(svgItem)
			}
		}
	}
}

func (b *BoardView) readPiecePixmaps() {
	// (" KPNBRQ- kpnbrq-")
	piecesFileNames := []string{"", "wK", "wP", "wN", "wB", "wR", "wQ", "", "", "bK", "bP", "bN", "bB", "bR", "bQ", ""}
	for _, piece := range piecesFileNames {
		if piece == "" {
			piecesSvg = append(piecesSvg, nil)
			continue
		}
		filePath := pieceImagePath + piece + ".svg"
		renderer := svg.NewQSvgRenderer(nil)
		renderer.Load(filePath)
		piecesSvg = append(piecesSvg, renderer)
	}
}
