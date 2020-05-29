package views

import (
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/svg"
	"github.com/therecipe/qt/widgets"

	"github.com/frankkopp/FrankyGo/pkg/types"

	"github.com/frankkopp/chessuiqt/internal/board"
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
	piecesSvg  = make([]*svg.QSvgRenderer, 0)
)

func init() {
	// read piece images
	readPieceImages()
}

type BoardView struct {
	boardPane      *widgets.QWidget
	boardView      *widgets.QGraphicsView
	boardScene     *widgets.QGraphicsScene
	fenView        *widgets.QLineEdit
	boardDetail    *widgets.QWidget
	nextPlayerView *widgets.QLabel

	board      *board.Board
	dragFlag   bool
	fromSquare types.Square
	toSquare   types.Square
}

// NewBoardView Creates a board view and its data structure.
func NewBoardView(widget *widgets.QWidget) *BoardView {
	bv := &BoardView{}
	bv.board = board.NewBoard()
	bv.fromSquare = types.SqNone
	bv.toSquare = types.SqNone
	bv.newBoardView(widget)
	return bv
}

// View returns the view of the board view
func (b *BoardView) View() *widgets.QWidget {
	return b.boardPane
}

func (b *BoardView) newBoardView(widget *widgets.QWidget) {

	// text box for fen
	b.fenView = widgets.NewQLineEdit2(b.board.ToFen(), nil)
	b.fenView.SetStyleSheet("background: white; color: black")
	b.fenView.ConnectTextEdited(b.editTextEvent)
	b.fenView.ConnectEditingFinished(b.newFenEvent)

	// next player view
	b.nextPlayerView = widgets.NewQLabel(nil, 0)
	b.nextPlayerView.SetContentsMargins(5, 0, 5, 0)
	b.nextPlayerView.SetText("Next Player: White")
	b.nextPlayerView.SetStyleSheet("background: white; color: black")
	b.nextPlayerView.ConnectMousePressEvent(b.nextPlayerFlipEvent)

	// row for details of fen
	b.boardDetail = widgets.NewQWidget(nil, 0)
	b.boardDetail.SetLayout(widgets.NewQHBoxLayout2(nil))
	b.boardDetail.Layout().SetContentsMargins(0, 0, 0, 0)
	b.boardDetail.Layout().SetAlignment2(b.boardDetail.Layout(), core.Qt__AlignLeft)
	b.boardDetail.Layout().AddWidget(b.nextPlayerView)
	b.boardDetail.Layout().AddItem(widgets.NewQSpacerItem(0, 0, widgets.QSizePolicy__MinimumExpanding, widgets.QSizePolicy__Ignored))

	// scene to draw the chess board on
	b.boardScene = widgets.NewQGraphicsScene(nil)

	// view for chess board
	b.boardView = widgets.NewQGraphicsView(nil)
	b.boardView.SetFrameStyle(int(widgets.QFrame__Sunken) | int(widgets.QFrame__StyledPanel))
	b.boardView.SetStyleSheet("background: yellow")
	b.boardView.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	b.boardView.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOff)
	b.boardView.SetScene(b.boardScene)

	// place the board view, fenView and detail row onto the board pane
	b.boardPane = widgets.NewQWidget(widget, 0)
	b.boardPane.SetStyleSheet("background: blue")
	b.boardPane.SetLayout(widgets.NewQVBoxLayout2(nil))
	b.boardPane.Layout().SetAlignment2(b.boardPane.Layout(), core.Qt__AlignTop)
	b.boardPane.Layout().SetSizeConstraint(widgets.QLayout__SetMinimumSize)
	b.boardPane.Layout().AddWidget(b.boardView)
	b.boardPane.Layout().AddWidget(b.fenView)
	b.boardPane.Layout().AddWidget(b.boardDetail)
	b.boardPane.ConnectKeyPressEvent(b.copyPasteEvent)

	// redraw the chess board on resize
	b.boardView.ConnectResizeEvent(b.resizeEvent)
}

func (b *BoardView) copyPasteEvent(event *gui.QKeyEvent) {
	if event.Matches(gui.QKeySequence__Copy) {
		app.Clipboard().SetText(b.board.ToFen(), gui.QClipboard__Clipboard)
		statusbar.ShowMessage("Copied fen to clipboard", 2000)

	}
	if event.Matches(gui.QKeySequence__Paste) {
		if app.Clipboard().MimeData(gui.QClipboard__Clipboard).HasText() {
			statusbar.ShowMessage("Paste fen OK", 2000)
			err := b.newFen(app.Clipboard().MimeData(gui.QClipboard__Clipboard).Text())
			if err == nil {
				return
			}
		}
		statusbar.ShowMessage("Paste: No valid fen in clipboard", 2000)
	}
}

func (b *BoardView) resizeEvent(event *gui.QResizeEvent) {
	b.drawBoard()
	b.boardView.SetFixedHeight(event.Size().Width() + 2)
}

func (b *BoardView) drawBoard() {

	boardSize := float64(b.boardView.Width() - 3)

	// Clear boardScene and add all children again
	b.boardScene.Clear()

	// squares size
	squareSize := boardSize / 8
	fontSize := int(squareSize * 0.15)

	// squares
	for rank := 8; rank >= 1; rank-- {
		for file := 1; file <= 8; file++ {

			square := (rank-1)*8 + file - 1
			upperLeftX := float64(file-1) * squareSize
			upperLeftY := float64(8-rank) * squareSize
			font.SetPixelSize(fontSize)

			// checkers
			bgcolor := brushWhite
			if (rank+file)%2 == 0 {
				bgcolor = brushBlack
			}
			rect := b.boardScene.AddRect2(upperLeftX, upperLeftY, squareSize, squareSize, border, bgcolor)
			rect.SetZValue(0.0)

			// rank number
			if file == 1 {
				t := b.boardScene.AddSimpleText(strconv.Itoa(rank), font)
				t.SetX(upperLeftX + (squareSize * 0.05))
				t.SetY(upperLeftY + (squareSize * 0.05))
				t.SetZValue(0.1)
			}
			// file letter
			if rank == 1 {
				t := b.boardScene.AddSimpleText(string('a'+file-1), font)
				t.SetX(upperLeftX + (squareSize * 0.85))
				t.SetY(upperLeftY + (squareSize * 0.80))
				t.SetZValue(0.1)
			}

			// debug - square index numbers
			font.SetPixelSize(fontSize / 2)
			t := b.boardScene.AddSimpleText(strconv.Itoa(square), font)
			t.SetX(upperLeftX + (squareSize * 0.8))
			t.SetY(upperLeftY + (squareSize * 0.05))

			// pieces for current board
			piece := b.board.GetPiece(types.Square(square))
			if piece != types.PieceNone {
				svgItem := svg.NewQGraphicsSvgItem(nil)
				svgRenderer := piecesSvg[piece]
				svgItem.SetSharedRenderer(svgRenderer)
				// current image size
				h := float64(svgRenderer.DefaultSize().Height())
				w := float64(svgRenderer.DefaultSize().Width())
				// piece size relative to square size
				scaleFactor := squareSize / h * (0.9)
				svgItem.SetScale(scaleFactor)
				h *= scaleFactor
				w *= scaleFactor
				// pos middle of piece relative to top left corner of square
				offsetX := h * 0.5
				offsetY := w * 0.5
				svgItem.SetX(upperLeftX + (squareSize * 0.5) - offsetX)
				svgItem.SetY(upperLeftY + (squareSize * 0.5) - offsetY)
				// make sure pieces are in foreground
				svgItem.SetZValue(0.5)
				svgItem.Show()
				b.boardScene.AddItem(svgItem)

				// mouse event registration
				// mouse press
				svgItem.ConnectMousePressEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
					b.fromSquare = getSquare(event.ScenePos().X(), event.ScenePos().Y(), squareSize)
				})
				// mouse drag
				svgItem.ConnectMouseMoveEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
					b.dragFlag = true
					svgItem.SetZValue(1.0)
					dragX := event.ScenePos().X() - offsetX
					dragY := event.ScenePos().Y() - offsetY
					// avoid dragging out of the view
					if dragX < 0 {
						dragX = 0
					}
					if dragX+h > boardSize {
						dragX = boardSize - h
					}
					if dragY < 0 {
						dragY = 0
					}
					if dragY+w > boardSize {
						dragY = boardSize - w
					}
					// move the piece image
					svgItem.SetPos2(dragX, dragY)
				})
				// mouse release
				svgItem.ConnectMouseReleaseEvent(func(event *widgets.QGraphicsSceneMouseEvent) {
					if b.dragFlag {
						b.toSquare = getSquare(svgItem.X()+(w/2), svgItem.Y()+(h/2), squareSize)
						b.board.MovePiece(b.fromSquare, b.toSquare)
						b.drawBoard()
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
	b.fenView.SetText(b.board.ToFen())
	if b.board.NextPlayer() == types.White {
		b.nextPlayerView.SetText("Next Player: White")
		b.nextPlayerView.SetStyleSheet("background: white; color: black")
	} else {
		b.nextPlayerView.SetText("Next Player: Black")
		b.nextPlayerView.SetStyleSheet("background: black; color: white")
	}
}

func (b *BoardView) editTextEvent(text string) {
	_, err := board.NewBoardFen(text)
	if err != nil {
		b.fenView.SetToolTip(err.Error())
		b.fenView.SetStyleSheet("background: white; color: red")
		return
	}
	b.fenView.SetToolTip("fen ok")
	b.fenView.SetStyleSheet("background: white; color: black")
}

func (b *BoardView) newFenEvent() {
	_ = b.newFen(b.fenView.Text())
}

func (b *BoardView) newFen(fen string) error {
	_, err := board.NewBoardFen(fen)
	if err == nil {
		_ = b.board.FromFen(fen)
		b.drawBoard()
		return nil
	}
	return err
}

func (b *BoardView) nextPlayerFlipEvent(event *gui.QMouseEvent) {
	b.board.FlipNextPlayer()
	b.drawBoard()
}

func getSquare(x float64, y float64, squareSize float64) types.Square {
	return types.SquareOf(types.File(x/squareSize), types.Rank(7-int(y/squareSize)))
}

// init piece images
func readPieceImages() {
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
