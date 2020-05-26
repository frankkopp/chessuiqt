package views

import (
	"fmt"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/svg"
	"github.com/therecipe/qt/widgets"

	"github.com/frankkopp/FrankyGo/pkg/movegen"
	"github.com/frankkopp/FrankyGo/pkg/position"
	"github.com/frankkopp/FrankyGo/pkg/types"
)

const (
	pieceImagePath = "./assets/qml/"
)

// colors
var (
	border     = gui.NewQPen3(gui.NewQColor6("black"))
	brushWhite = gui.NewQBrush3(gui.NewQColor6("white"), core.Qt__SolidPattern)
	brushBlack = gui.NewQBrush3(gui.NewQColor6("darkgrey"), core.Qt__SolidPattern)
	font       = gui.NewQFont()
	pieces     = make([]*gui.QImage, 0)
	piecesSvg  = make([]*svg.QSvgRenderer, 0)
)

type BoardView struct {
	boardView  *widgets.QGraphicsView
	boardScene *widgets.QGraphicsScene
	position   *position.Position
	movegen    *movegen.Movegen
	dragFlag   bool
	fromSquare types.Square
	toSquare   types.Square
}

// NewBoardView Creates a board view and its data structure.
func NewBoardView(widget *widgets.QWidget) *BoardView {
	bv := &BoardView{}
	bv.newBoardView(widget)
	bv.position = position.NewPosition()
	bv.movegen = movegen.NewMoveGen()
	bv.fromSquare = types.SqNone
	bv.toSquare = types.SqNone
	return bv
}

// View returns the view of the board view
func (b *BoardView) View() *widgets.QGraphicsView {
	return b.boardView
}

func (b *BoardView) newBoardView(widget *widgets.QWidget) {

	// pane for chess board
	b.boardView = widgets.NewQGraphicsView(widget)
	b.boardView.SetStyleSheet("background: yellow")
	b.boardView.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	b.boardView.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)

	// scene to draw the chess board on
	b.boardScene = widgets.NewQGraphicsScene(b.boardView)
	b.boardView.SetScene(b.boardScene)
	b.boardView.SetMinimumSize2(int(b.boardScene.Width()), int(b.boardScene.Height()))

	// read piece images
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

	// squares size
	squareSize := boardSize / 8
	fontSize := int(squareSize * 0.15)

	// squares
	for rank := 8; rank >= 1; rank-- {
		for file := 1; file <= 8; file++ {
			square := (rank-1)*8 + file - 1

			// checkers
			if (rank+file)%2 == 0 {
				rect := b.boardScene.AddRect2(float64(file-1)*squareSize, float64(8-rank)*squareSize, squareSize, squareSize, border, brushBlack)
				rect.SetZValue(0.0)
			} else {
				rect := b.boardScene.AddRect2(float64(file-1)*squareSize, float64(8-rank)*squareSize, squareSize, squareSize, border, brushWhite)
				rect.SetZValue(0.0)
			}

			font.SetPixelSize(fontSize)
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
				// piece size relative to square size
				scaleFactor := squareSize / height * (0.9)
				svgItem.SetScale(scaleFactor)
				// piece pos relative to top left corner of square
				offsetX := height * 0.5 * scaleFactor
				offsetY := width * 0.5 * scaleFactor
				svgItem.SetX((float64(file-1) * squareSize) + (squareSize * 0.5) - offsetX)
				svgItem.SetY((float64(8-rank) * squareSize) + (squareSize * 0.5) - offsetY)
				// make sure pieces are in foreground
				svgItem.SetZValue(0.5)
				svgItem.Show()
				b.boardScene.AddItem(svgItem)
				// mouse event registration
				svgItem.ConnectMousePressEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
					b.fromSquare = getSquare(event, squareSize)
				})
				svgItem.ConnectMouseMoveEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
					b.dragFlag = true
					svgItem.SetZValue(1.0)
					svgItem.SetPos2(event.ScenePos().X()-offsetX, event.ScenePos().Y()-offsetY)
				})
				svgItem.ConnectMouseReleaseEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
					if b.dragFlag {
						b.toSquare = getSquare(event, squareSize)
						fmt.Println("Was drag from", b.fromSquare.String(), "to", b.toSquare.String())
					}
					// reset drag
					b.dragFlag = false
					b.fromSquare = types.SqNone
					b.toSquare = types.SqNone
					svgItem.SetZValue(0.5)
				})
			}
		}
	}
}

func getSquare(e *widgets.QGraphicsSceneMouseEvent, squareSize float64) types.Square {
	return types.SquareOf(types.File(int(e.ScenePos().X()/squareSize)), types.Rank(7-int(e.ScenePos().Y()/squareSize)))
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
