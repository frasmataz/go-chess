package game

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFENConversion(t *testing.T) {
	type FENtest struct {
		input string
		want  string
	}

	tests := map[string]FENtest{
		"starting_state": {
			input: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			want:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		"en passant": {
			input: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
			want:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		},
		"no castling": {
			input: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b - - 0 1",
			want:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b - - 0 1",
		},
		"castling": {
			input: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w Qk - 0 1",
			want:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w Qk - 0 1",
		},
		"movecounts": {
			input: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w k - 49 75",
			want:  "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w k - 49 75",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			board, err := BoardFromFEN(test.input)
			if err != nil {
				t.Error(err)
			}

			output := board.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(board.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})
	}
}

func TestPositionToIndex(t *testing.T) {
	type positionTest struct {
		input     string
		want      [2]int
		wantError bool
	}

	tests := []positionTest{
		{
			input: "a1",
			want:  [2]int{7, 0},
		},
		{
			input: "h8",
			want:  [2]int{0, 7},
		},
		{
			input: "a8",
			want:  [2]int{0, 0},
		},
		{
			input: "h1",
			want:  [2]int{7, 7},
		},
		{
			input: "f6",
			want:  [2]int{2, 5},
		},
		{
			input: "d2",
			want:  [2]int{6, 3},
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
			output, err := positionToIndex(test.input)
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

func TestIndexToPosition(t *testing.T) {
	type indexTest struct {
		input     [2]int
		want      string
		wantError bool
	}

	tests := []indexTest{
		{
			input: [2]int{7, 0},
			want:  "a1",
		},
		{
			input: [2]int{0, 7},
			want:  "h8",
		},
		{
			input: [2]int{0, 0},
			want:  "a8",
		},
		{
			input: [2]int{7, 7},
			want:  "h1",
		},
		{
			input: [2]int{2, 5},
			want:  "f6",
		},
		{
			input: [2]int{6, 3},
			want:  "d2",
		},
		{
			input:     [2]int{-1, 5},
			wantError: true,
		},
		{
			input:     [2]int{8, 5},
			wantError: true,
		},
		{
			input:     [2]int{2, -1},
			wantError: true,
		},
		{
			input:     [2]int{2, 8},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.input), func(t *testing.T) {
			output, err := indexToPosition(test.input)
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

func TestPrintGameState(t *testing.T) {
	board, err := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Errorf("error testing TestPrintGameState: %v", err)
	}
	t.Log(board.PrintGameState())
}

func TestPawnMoves(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"edges": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []string{
				"a2a4",
				"a7a5",
				"h2h4",
				"h7h5",
			},
			want:      "rnbqkbnr/1pppppp1/8/p6p/P6P/8/1PPPPPP1/RNBQKBNR w KQkq - 0 3",
			wantError: false,
		},
		"capture": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []string{
				"f2f4",
				"c7c5",
				"f4f5",
				"c5c4",
				"f5f6",
				"c4c3",
				"f6e7",
				"c3d2",
			},
			want:      "rnbqkbnr/pp1pPppp/8/8/8/8/PPPpP1PP/RNBQKBNR w KQkq - 0 5",
			wantError: false,
		},
		"white can't move three from start": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []string{
				"f2f5",
			},
			wantError: true,
		},
		"white can't move diagonal without take": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []string{
				"f2e3",
			},
			wantError: true,
		},
		"white can't move forward on occupied square": {
			starting_state: "rnbqkbnr/ppp2ppp/3pp3/8/8/3BP3/PPPP1PPP/RNBQK1NR w KQkq - 0 3",
			moves: []string{
				"d2d3",
			},
			wantError: true,
		},
		"white can't move two forward on occupied space": {
			starting_state: "rnbqkbnr/ppp2ppp/3p4/4p3/2B1P3/8/PPPP1PPP/RNBQK1NR w KQkq - 0 3",
			moves: []string{
				"c2c4",
			},
			wantError: true,
		},
		"white can't jump two forward over occupied square": {
			starting_state: "rnbqkbnr/ppp2ppp/3pp3/8/8/3BP3/PPPP1PPP/RNBQK1NR w KQkq - 0 3",
			moves: []string{
				"d2d4",
			},
			wantError: true,
		},
		"white can't take one forward": {
			starting_state: "rnbqkbnr/ppp1pppp/4P3/8/8/3p4/PPPP1PPP/RNBQKBNR w KQkq - 0 4",
			moves: []string{
				"d2d3",
			},
			wantError: true,
		},
		"white can't take two forward": {
			starting_state: "rnbqkbnr/ppp1pppp/8/4P3/3p4/8/PPPP1PPP/RNBQKBNR w KQkq - 0 3",
			moves: []string{
				"d2d4",
			},
			wantError: true,
		},
		"black can't move three from start": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			moves: []string{
				"f7f4",
			},
			wantError: true,
		},
		"black can't move diagonal without take": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
			moves: []string{
				"f7e6",
			},
			wantError: true,
		},
		"black can't move forward on occupied square": {
			starting_state: "rnbqk1nr/pppp1ppp/3bp3/4P3/8/5P2/PPPP2PP/RNBQKBNR b KQkq - 0 3",
			moves: []string{
				"d7d6",
			},
			wantError: true,
		},
		"black can't move two forward on occupied space": {
			starting_state: "rnbqk1nr/pppp1ppp/4p3/2b1P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
			moves: []string{
				"c7c5",
			},
			wantError: true,
		},
		"black can't jump two forward over occupied square": {
			starting_state: "rnbqk1nr/pppp1ppp/3bp3/4P3/8/5P2/PPPP2PP/RNBQKBNR b KQkq - 0 3",
			moves: []string{
				"d7d5",
			},
			wantError: true,
		},
		"black can't take one forward": {
			starting_state: "rnbqkbnr/ppp1pppp/4P3/8/3p4/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
			moves: []string{
				"e7e6",
			},
			wantError: true,
		},
		"black can't take two forward": {
			starting_state: "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 2",
			moves: []string{
				"e7e5",
			},
			wantError: true,
		},
		// TODO: EN PASSANT IS A LEGITMATE MOVE
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})
	}
}

func TestKnightMoves(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"valid moves": {
			starting_state: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []string{
				"g1f3",
				"g8f6",
				"f3e5",
				"f6d5",
				"e5d7",
				"d5b6",
				"d7b6",
				"b8c6",
				"b1c3",
				"c6d4",
				"c3a4",
				"d4e2",
				"b6a8",
			},
			want:      "N1bqkb1r/ppp1pppp/8/8/N7/8/PPPPnPPP/R1BQKB1R b KQk - 0 7",
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})
	}
}

func TestBishopMoves(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"valid moves": {
			starting_state: "rnbqkbnr/p1p2p1p/1p1pp1p1/8/8/1P1PP1P1/P1P2P1P/RNBQKBNR w KQkq - 0 5",
			moves: []string{
				"f1e2",
				"f8h6",
				"c1a3",
				"c8d7",
				"e2g4",
				"h6e3",
				"a3d6",
				"d7c6",
				"g4e6",
				"c6h1",
				"d6c7",
			},
			want:      "rn1qk1nr/p1B2p1p/1p2B1p1/8/8/1P1Pb1P1/P1P2P1P/RN1QK1Nb b Qkq - 0 10",
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})

	}
}

func TestRookMoves(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"valid moves": {
			starting_state: "rnbqkbnr/1pppppp1/8/p6p/P6P/8/1PPPPPP1/RNBQKBNR w KQkq - 0 3",
			moves: []string{
				"h1h3",
				"h8h6",
				"a1a3",
				"a8a6",
				"a3d3",
				"a6f6",
				"d3d7",
				"f6f2",
				"d7d6",
				"h6d6",
			},
			want:      "1nbqkbn1/1pp1ppp1/3r4/p6p/P6P/7R/1PPPPrP1/1NBQKBN1 w - - 0 8",
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})

	}
}

func TestQueenMoves(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"valid moves": {
			starting_state: "rnbqkbnr/pp1p1ppp/2p1p3/8/8/2P1P3/PP1P1PPP/RNBQKBNR w KQkq - 0 3",
			moves: []string{
				"d1g4",
				"d8a5",
				"g4g7",
				"a5a2",
				"g7g3",
				"a2c4",
				"g3b8",
				"c4d5",
				"b8a8",
			},
			want:      "Q1b1kbnr/pp1p1p1p/2p1p3/3q4/8/2P1P3/1P1P1PPP/RNB1KBNR b KQk - 0 7",
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})

	}
}

func TestKingMoves(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"valid moves": {
			starting_state: "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
			moves: []string{
				"e1e2",
				"e8e7",
				"e2d3",
				"e7f6",
				"d3c4",
				"f6g5",
				"c4b5",
				"g5f4",
				"b5a4",
				"f4e4",
			},
			want:      "rnbq1bnr/pppp1ppp/8/4p3/K3k3/8/PPPP1PPP/RNBQ1BNR w - - 0 7",
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})

	}
}

func TestSelfCheckPrevention(t *testing.T) {
	type FENtest struct {
		starting_state string
		moves          []string
		want           string
		wantError      bool
	}

	tests := map[string]FENtest{
		"pawn/w/1": {
			starting_state: "8/8/4k3/3p4/8/3K4/8/8 w - - 0 2",
			moves: []string{
				"d3c4",
			},
			wantError: true,
		},
		"pawn/w/2": {
			starting_state: "8/8/4k3/3p4/8/3K4/8/8 w - - 0 2",
			moves: []string{
				"d3e4",
			},
			wantError: true,
		},
		"pawn/w/3": {
			starting_state: "8/8/4k3/3p4/8/3K4/8/8 w - - 0 2",
			moves: []string{
				"d3d4",
			},
			want: "8/8/4k3/3p4/3K4/8/8/8 b - - 1 2",
		},
		"pawn/b/1": {
			starting_state: "8/8/4k3/8/4P3/3K4/8/8 b - - 0 2",
			moves: []string{
				"e6f5",
			},
			wantError: true,
		},
		"pawn/b/2": {
			starting_state: "8/8/4k3/8/4P3/3K4/8/8 b - - 0 2",
			moves: []string{
				"e6d5",
			},
			wantError: true,
		},
		"pawn/b/3": {
			starting_state: "8/8/4k3/8/4P3/3K4/8/8 b - - 0 2",
			moves: []string{
				"e6e5",
			},
			want: "8/8/8/4k3/4P3/3K4/8/8 w - - 1 3",
		},
		"bishop/w/1": {
			starting_state: "8/8/4k3/2b5/8/3K4/8/8 w - - 0 2",
			moves: []string{
				"d3d4",
			},
			wantError: true,
		},
		"bishop/w/2": {
			starting_state: "8/8/4k3/2b5/8/3K4/8/8 w - - 0 2",
			moves: []string{
				"d3e3",
			},
			wantError: true,
		},
		"bishop/w/3": {
			starting_state: "8/8/4k3/2b5/8/3K4/8/8 w - - 0 2",
			moves: []string{
				"d3c4",
			},
			want: "8/8/4k3/2b5/2K5/8/8/8 b - - 1 2",
		},
		"bishop/b/1": {
			starting_state: "8/8/4k3/8/2K2B2/8/8/8 b - - 0 2",
			moves: []string{
				"e6d6",
			},
			wantError: true,
		},
		"bishop/b/2": {
			starting_state: "8/8/4k3/8/2K2B2/8/8/8 b - - 0 2",
			moves: []string{
				"e6e5",
			},
			wantError: true,
		},
		"bishop/b/3": {
			starting_state: "8/8/4k3/8/2K2B2/8/8/8 b - - 0 2",
			moves: []string{
				"e6f5",
			},
			want: "8/8/8/5k2/2K2B2/8/8/8 w - - 1 3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				err := game.ExecuteMove(move)
				if err != nil {
					if !test.wantError {
						t.Errorf("unexpected error executing move - test: '%v', move: '%v', err: '%v'", test, move, err)
					}
					return
				}
			}

			if test.wantError {
				t.Errorf("did not throw error as expected - test: '%v''", test)
				return
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})

	}
}
