package board

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	. "github.com/frankkopp/FrankyGo/pkg/types"
	"github.com/stretchr/testify/assert"
)

// make tests run in the projects root directory
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

// Setup the tests
func TestMain(m *testing.M) {
	// setup code here
	code := m.Run()
	os.Exit(code)
}

func TestBoardCreation(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err := NewBoardFen(fen)
	assert.NoError(t, err)
	fmt.Print(p.String())
	assert.Equal(t, White, p.nextPlayer)
	assert.Equal(t, CastlingAny, p.castling)
	assert.Equal(t, SqNone, p.enpassant)
	assert.Equal(t, 0, p.halfmoveclock)
	assert.Equal(t, 1, p.movenumber)
	assert.Equal(t, fen, p.StringFen())

	fmt.Println()

	fen = "r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3 0 14"
	p, err = NewBoardFen(fen)
	assert.NoError(t, err)
	fmt.Print(p.String())
	assert.Equal(t, Black, p.nextPlayer)
	assert.Equal(t, CastlingBlack, p.castling)
	assert.Equal(t, SqE3, p.enpassant)
	assert.Equal(t, 0, p.halfmoveclock)
	assert.Equal(t, 14, p.movenumber)
	assert.Equal(t, fen, p.StringFen())

	// incorrect fens
	fen = "rnbqkbnr/ppppppp2/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err = NewBoardFen(fen)
	fmt.Println(err)
	assert.Error(t, err)
	fen = "rnbqkbnr/pppppppp/8/9/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err = NewBoardFen(fen)
	fmt.Println(err)
	assert.Error(t, err)
	fen = "rnbqkbnr/ppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err = NewBoardFen(fen)
	fmt.Println(err)
	assert.Error(t, err)
	fen = "rnbqkbnr/ppppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err = NewBoardFen(fen)
	fmt.Println(err)
	assert.Error(t, err)
	fen = "rnbqkbnr/pppppppp/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err = NewBoardFen(fen)
	fmt.Println(err)
	assert.Error(t, err)
	fen = "rnbqkbnr/ppppppp?/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err = NewBoardFen(fen)
	fmt.Println(err)
	assert.Error(t, err)
}
