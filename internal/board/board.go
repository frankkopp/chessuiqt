// Package board has data structures and functions for a representation of a FEN string.
// It does not know anything about chess rules and is only used to persist a fen as a data structure
// for use in a view. It translates a FEN string into the data structure and the data structure
// back into a FEN string.
package board

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/frankkopp/FrankyGo/pkg/types"
)

const (
	// StartFen is a string with the fen position for a standard chess game
	StartFen string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

// Board is a simple direct representation of a FEN string. It does not know anything about
// chess rules and is only used to persist a fen as a data structure for use in a view
type Board struct {
	board         [64]types.Piece
	nextPlayer    types.Color
	castling      types.CastlingRights
	enpassant     types.Square
	halfmoveclock int
	movenumber    int
}

// NewBoard creates a new board instance.
// When called without an argument the board will have the start position
// When a fen string is given it will create a position with based on this fen.
// Additional fens/strings are ignored
func NewBoard(fen ...string) *Board {
	if len(fen) == 0 {
		f, _ := NewBoardFen(StartFen)
		return f
	}
	f, _ := NewBoardFen(fen[0])
	return f
}

// NewPositionFen creates a new position with the given fen string
// as board position.
// It returns nil and an error if the fen was invalid.
func NewBoardFen(fen string) (*Board, error) {
	p := &Board{}
	if e := p.FromFen(fen); e != nil {
		return nil, e
	}
	return p, nil
}

// regex for first part of fen (position of pieces)
var regexFenIllegalChars = regexp.MustCompile("[^1-8pPnNbBrRqQkK/]+")

// regex for next player color in fen
var regexWorB = regexp.MustCompile("^[w|b]$")

// regex for castling rights in fen
var regexCastlingRights = regexp.MustCompile("^(K?Q?k?q?|-)$")

// regex for en passant square in fen
var regexEnPassant = regexp.MustCompile("^([a-h][1-8]|-)$")

// setupBoard sets up a board based on a fen. This is basically
// the only way to get a valid Position instance. Internal state
// will be setup as well as all struct data is initialized to 0.
func (b *Board) FromFen(fen string) error {

	// TODO: do not change the board until we know there is no error

	// clear the board
	b.Clear()

	// We will analyse the fen and only require the initial board layout part.
	// All other parts will have defaults. E.g. next player is white, no castling, etc.
	fen = strings.TrimSpace(fen)
	fenParts := strings.Split(fen, " ")

	if len(fenParts) == 0 {
		err := errors.New("fen must not be empty")
		return err
	}

	// make sure only valid chars are used
	match := regexFenIllegalChars.MatchString(fenParts[0])
	if match {
		err := errors.New("fen position contains invalid characters")
		return err
	}

	// a fen start at A8 - which is file==0 and rank==7
	file := 0
	rank := 7

	// loop over fen characters
	for _, c := range fenParts[0] {
		if number, e := strconv.Atoi(string(c)); e == nil { // is number
			file += number
			if file > 8 {
				return fmt.Errorf("too many squares (%d) in rank %d:  %s", file, rank+1, fenParts[0])
			}
		} else if string(c) == "/" { // find rank separator
			if file < 8 {
				return fmt.Errorf("not enough squares (%d) in rank %d:  %s", file, rank+1, fenParts[0])
			}
			if file > 8 {
				return fmt.Errorf("too many squares (%d) in rank %d:  %s", file, rank+1, fenParts[0])
			}
			file = 0
			rank--
			if rank < 0 {
				return fmt.Errorf("too many ranks (%d):  %s", 8-rank, fenParts[0])
			}
		} else { // find piece type
			piece := types.PieceFromChar(string(c))
			if piece == types.PieceNone {
				return fmt.Errorf("invalid piece character '%s' in %s", string(c), fenParts[0])
			}
			if file > 7 {
				return fmt.Errorf("too many squares (%d) in rank %d:  %s", file+1, rank+1, fenParts[0])
			}
			currentSquare := types.SquareOf(types.File(file), types.Rank(rank))
			if !currentSquare.IsValid() {
				return fmt.Errorf("invalid square %d (%s): %s", currentSquare, currentSquare.String(), fenParts[0])
			}
			b.PutPiece(piece, currentSquare)
			file++
		}
	}
	if file != 8 || rank != 0 { // after h1++ we reach a2 - a2 needs to be last current square
		return fmt.Errorf("not reached last square (file=%d, rank=%d) after reading fen", file, rank)
	}

	// set defaults
	b.movenumber = 1
	b.enpassant = types.SqNone

	// everything below is optional as we can apply defaults

	// next player
	if len(fenParts) >= 2 {
		match = regexWorB.MatchString(fenParts[1])
		if !match {
			err := errors.New("fen next player contains invalid characters")
			return err
		}
		switch fenParts[1] {
		case "w":
			b.nextPlayer = types.White
		case "b":
			b.nextPlayer = types.Black
		}
	}

	// castling rights
	if len(fenParts) >= 3 {
		match = regexCastlingRights.MatchString(fenParts[2])
		if !match {
			err := errors.New("fen castling rights contains invalid characters")
			return err
		}
		// are there  rights to be encoded?
		if fenParts[2] != "-" {
			for _, c := range fenParts[2] {
				switch string(c) {
				case "K":
					b.castling.Add(types.CastlingWhiteOO)
				case "Q":
					b.castling.Add(types.CastlingWhiteOOO)
				case "k":
					b.castling.Add(types.CastlingBlackOO)
				case "q":
					b.castling.Add(types.CastlingBlackOOO)
				}
			}
		}
	}

	// enpassant
	if len(fenParts) >= 4 {
		match = regexEnPassant.MatchString(fenParts[3])
		if !match {
			err := errors.New("fen castling rights contains invalid characters")
			return err
		}
		if fenParts[3] != "-" {
			b.enpassant = types.MakeSquare(fenParts[3])
		}
	}

	// half move clock (50 moves rule)
	if len(fenParts) >= 5 {
		if number, e := strconv.Atoi(fenParts[4]); e == nil { // is number
			b.halfmoveclock = number
		} else {
			return e
		}
	}

	// move number
	if len(fenParts) >= 6 {
		if moveNumber, e := strconv.Atoi(fenParts[5]); e == nil { // is number
			if moveNumber == 0 {
				moveNumber = 1
			}
			b.movenumber = moveNumber
		} else {
			return e
		}
	}

	// return without error
	return nil
}

func (b *Board) ToFen() string {
	var fen strings.Builder
	// pieces
	for r := types.Rank1; r <= types.Rank8; r++ {
		emptySquares := 0
		for f := types.FileA; f <= types.FileH; f++ {
			pc := b.board[types.SquareOf(f, types.Rank8-r)]
			if pc == types.PieceNone {
				emptySquares++
			} else {
				if emptySquares > 0 {
					fen.WriteString(strconv.Itoa(emptySquares))
					emptySquares = 0
				}
				fen.WriteString(pc.String())
			}
		}
		if emptySquares > 0 {
			fen.WriteString(strconv.Itoa(emptySquares))
		}
		if r < types.Rank8 {
			fen.WriteString("/")
		}
	}
	// next player
	fen.WriteString(" ")
	fen.WriteString(b.nextPlayer.String())
	// castling
	fen.WriteString(" ")
	fen.WriteString(b.castling.String())
	// en passant
	fen.WriteString(" ")
	fen.WriteString(b.enpassant.String())
	// half move clock
	fen.WriteString(" ")
	fen.WriteString(strconv.Itoa(b.halfmoveclock))
	// full move number
	fen.WriteString(" ")
	fen.WriteString(strconv.Itoa(b.movenumber))
	// return fen string
	return fen.String()
}

func (b *Board) Clear() {
	for sq := 0; sq < 64; sq++ {
		b.board[sq] = types.PieceNone
	}
	b.nextPlayer = types.White
	b.castling = types.CastlingNone
	b.enpassant = types.SqNone
	b.halfmoveclock = 0
	b.movenumber = 0
}

func (b *Board) GetPiece(square types.Square) types.Piece {
	return b.board[square]
}

func (b *Board) MovePiece(from types.Square, to types.Square) types.Piece {
	return b.PutPiece(b.RemovePiece(from), to)
}

func (b *Board) PutPiece(piece types.Piece, to types.Square) types.Piece {
	// update board and return previous piece on square
	tmp := b.board[to]
	b.board[to] = piece
	return tmp

}

func (b *Board) RemovePiece(from types.Square) types.Piece {
	// update board and return previous piece on square
	tmp := b.board[from]
	b.board[from] = types.PieceNone
	return tmp
}

func (b *Board) NextPlayer() types.Color {
	return b.nextPlayer
}

func (b *Board) FlipNextPlayer() {
	b.nextPlayer = b.nextPlayer.Flip()
}

func (b *Board) HalfMoveClock() int {
	return b.halfmoveclock
}

// String returns a string representing the board instance. This
// includes the fen, a board matrix, game phase, material and pos values.
func (b *Board) String() string {
	var os strings.Builder
	os.WriteString(b.StringFen())
	os.WriteString("\n")
	os.WriteString(b.StringBoard())
	os.WriteString("\n")
	os.WriteString(fmt.Sprintf("Next Player    : %s\n", b.nextPlayer.String()))
	return os.String()
}

// StringFen returns a string with the FEN of the current position
func (b *Board) StringFen() string {
	return b.ToFen()
}

// StringBoard returns a visual matrix of the board and pieces
func (b *Board) StringBoard() string {
	var os strings.Builder
	os.WriteString("+---+---+---+---+---+---+---+---+\n")
	for r := types.Rank1; r <= types.Rank8; r++ {
		for f := types.FileA; f <= types.FileH; f++ {
			os.WriteString("| ")
			os.WriteString(b.board[types.SquareOf(f, types.Rank8-r)].UniChar())
			os.WriteString(" ")
		}
		os.WriteString("|\n+---+---+---+---+---+---+---+---+\n")
	}
	return os.String()
}

func (b *Board) SetHalfMoveClock(i int) {
	b.halfmoveclock = i
}
