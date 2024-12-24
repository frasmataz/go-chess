package game

import (
	"reflect"
	"testing"
)

func TestBoardFromFEN(t *testing.T) {
	type FENtest struct {
		input string
		want  gameState
	}

	tests := map[string]FENtest{
		"starting_state": {
			input: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want: gameState{
				boardState: [8][8]Piece{
					{GetPiece("blackRook"), GetPiece("blackKnight"), GetPiece("blackBishop"), GetPiece("blackQueen"), GetPiece("blackKing"), GetPiece("blackBishop"), GetPiece("blackKnight"), GetPiece("blackRook")},
					{GetPiece("blackPawn"), GetPiece("blackPawn"), GetPiece("blackPawn"), GetPiece("blackPawn"), GetPiece("blackPawn"), GetPiece("blackPawn"), GetPiece("blackPawn"), GetPiece("blackPawn")},
					{GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space")},
					{GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space")},
					{GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space")},
					{GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space"), GetPiece("space")},
					{GetPiece("whitePawn"), GetPiece("whitePawn"), GetPiece("whitePawn"), GetPiece("whitePawn"), GetPiece("whitePawn"), GetPiece("whitePawn"), GetPiece("whitePawn"), GetPiece("whitePawn")},
					{GetPiece("whiteRook"), GetPiece("whiteKnight"), GetPiece("whiteBishop"), GetPiece("whiteQueen"), GetPiece("whiteKing"), GetPiece("whiteBishop"), GetPiece("whiteKnight"), GetPiece("whiteRook")},
				},
				nextPlayer: White,
				castlingRights: castlingRights{
					blackKingCastle:  true,
					blackQueenCastle: true,
					whiteKingCastle:  true,
					whiteQueenCastle: true,
				},
				enPassantTarget: nil,
				halfmoveClock:   0,
				fullmoveClock:   1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			board, err := BoardFromFEN(test.input)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(*board, test.want) {
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, *board)
			}
		})
	}
}

func TestPositionToIndex(t *testing.T) {
	type positionTest struct {
		input     string
		want      *[2]int
		wantError bool
	}

	tests := []positionTest{
		{
			input: "a1",
			want:  &[2]int{7, 0},
		},
		{
			input: "h8",
			want:  &[2]int{0, 7},
		},
		{
			input: "a8",
			want:  &[2]int{0, 0},
		},
		{
			input: "h1",
			want:  &[2]int{7, 7},
		},
		{
			input: "f6",
			want:  &[2]int{2, 5},
		},
		{
			input: "d2",
			want:  &[2]int{6, 3},
		},
		{
			input:     "A1",
			wantError: true,
		},
		{
			input:     "m1",
			wantError: true,
		},
		{
			input:     "a9",
			wantError: true,
		},
		{
			input:     "a",
			wantError: true,
		},
		{
			input:     "a ",
			wantError: true,
		},
		{
			input:     "a",
			wantError: true,
		},
		{
			input:     "a ",
			wantError: true,
		},
		{
			input:     "1",
			wantError: true,
		},
		{
			input:     " 1",
			wantError: true,
		},
		{
			input:     "",
			wantError: true,
		},
		{
			input:     "test",
			wantError: true,
		},
		{
			input:     "#!",
			wantError: true,
		},
		{
			input:     " a1",
			wantError: true,
		},
		{
			input:     "a1 ",
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output, err := PositionToIndex(test.input)
			if err != nil {
				if !test.wantError {
					t.Errorf("PositionToIndex threw unexpected error: %v", err)
				}
				return
			}

			if test.wantError {
				t.Errorf("PositionToIndex did not throw error as expected - input: '%v', got '%v'", test.input, output)
			}

			if !reflect.DeepEqual(test.want, output) {
				t.Errorf("PositionToIndex returned unexpected output - input '%v', want '%v', got '%v'", test.input, test.want, output)
			}
		})
	}
}
