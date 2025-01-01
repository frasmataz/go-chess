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
			want: "rnbqkbnr/1pppppp1/8/p6p/P6P/8/1PPPPPP1/RNBQKBNR w KQkq - 0 3",
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
			want: "rnbqkbnr/pp1pPppp/8/8/8/8/PPPpP1PP/RNBQKBNR w KQkq - 0 5",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			game, err := BoardFromFEN(test.starting_state)
			if err != nil {
				t.Error(err)
			}

			for _, move := range test.moves {
				game.ExecuteMove(move)
			}

			output := game.BoardToFEN()
			if !reflect.DeepEqual(output, test.want) {
				fmt.Print(game.PrintGameState())
				t.Errorf("board did not match expected state.\n\nexpected:\n%v\n\ngot:\n%v\n\n", test.want, output)
			}
		})
	}
}
