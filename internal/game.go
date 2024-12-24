package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Status int

const (
	Normal    Status = 0
	Check     Status = 1
	Checkmate Status = 2
)

type castlingRights struct {
	blackKingCastle  bool
	blackQueenCastle bool
	whiteKingCastle  bool
	whiteQueenCastle bool
}

type gameState struct {
	boardState      [8][8]Piece // [7][0] is a1, [0][7] is h8 - [0][0] is top-left from white's perspective
	nextPlayer      PlayerColour
	castlingRights  castlingRights
	enPassantTarget *[2]int
	halfmoveClock   int
	fullmoveClock   int
}

func PositionToIndex(position string) (*[2]int, error) {
	// Converts standard board position (eg. "e4") to indeces in the [8][8]Piece board state
	positionRegex := `^[a-h]{1}[1-8]{1}$`
	re := regexp.MustCompile(positionRegex)
	if !re.MatchString(position) {
		return nil, fmt.Errorf("position did not match regex '%v' - got '%v'", positionRegex, position)
	}

	colIndex := int(position[0]) - 97 // 97 is int value of rune 'a'

	if colIndex < 0 || colIndex > 7 {
		return nil, fmt.Errorf("parsed row out-of-bounds, should be in range 0-7 - input '%v', parsed to row %v", position, colIndex)
	}

	rowRawInt, err := strconv.Atoi(string(rune(position[1]))) // FIXME: this can't be the best way
	if err != nil {
		return nil, err
	}

	rowIndex := 8 - int(rowRawInt) // board row indeces are top-down, position strings are bottom-up

	if rowIndex < 0 || rowIndex > 7 {
		return nil, fmt.Errorf("parsed column out-of-bounds, should be in range 0-7 - input '%v', parsed to column %v", position, rowIndex)
	}

	return &[2]int{rowIndex, colIndex}, nil
}

func BoardFromFEN(fen string) (*gameState, error) {
	// Converts a FEN board state string into a board objectgo run

	// https://www.chess.com/terms/fen-chess
	fenSegments := strings.Split(fen, " ")

	// First segment describes layout of pieces on the board
	rows := strings.Split(fenSegments[0], "/")
	board := gameState{}

	for rownum, row := range rows {
		boardColumn := 0

		for _, char := range row {
			// Digits indicate a number of consecutive spaces
			if unicode.IsDigit(char) {
				n, err := strconv.Atoi(string(char))
				if err != nil {
					return nil, fmt.Errorf("error parsing space digit: %v", err)
				}

				if n > 8 || n < 1 {
					return nil, fmt.Errorf("space digit out of range - must be between 1 and 8 inclusive, got '%v'", char)
				}

				for i := 0; i < n; i++ {
					board.boardState[rownum][boardColumn+i] = Pieces["space"]
				}
			} else {
				// If not digit, should be a piece
				piece, err := FENToPiece(char)
				if err != nil {
					return nil, err
				}

				board.boardState[rownum][boardColumn] = piece
			}

			boardColumn++
		}
	}

	// Second segment describes whose turn it is
	if fenSegments[1] == "w" {
		board.nextPlayer = White
	} else if fenSegments[1] == "b" {
		board.nextPlayer = Black
	} else {
		return nil, fmt.Errorf("next player section missing from FEN")
	}

	// Third segment describes castling rights
	if fenSegments[2] == "-" {
		board.castlingRights = castlingRights{
			blackKingCastle:  false,
			blackQueenCastle: false,
			whiteKingCastle:  false,
			whiteQueenCastle: false,
		}
	} else {
		if strings.Contains(fenSegments[2], "k") {
			board.castlingRights.blackKingCastle = true
		}

		if strings.Contains(fenSegments[2], "q") {
			board.castlingRights.blackQueenCastle = true
		}

		if strings.Contains(fenSegments[2], "K") {
			board.castlingRights.whiteKingCastle = true
		}

		if strings.Contains(fenSegments[2], "Q") {
			board.castlingRights.whiteQueenCastle = true
		}
	}

	// Fourth segment describes en passant targets
	if fenSegments[3] != "-" {
		enPassantTarget, err := PositionToIndex(fenSegments[3])
		if err != nil {
			return nil, err
		}
		board.enPassantTarget = enPassantTarget
	}

	// Fifth segment counts halfmoves
	halfmoves, err := strconv.Atoi(fenSegments[4])
	if err != nil {
		return nil, fmt.Errorf("could not parse halfmove clock - got '%v'", fenSegments[4])
	}
	if halfmoves < 0 {
		return nil, fmt.Errorf("halfmove clock cannot be negative - got '%v'", halfmoves)
	}

	board.halfmoveClock = halfmoves

	// Sixth segment counts halfmoves
	fullmoves, err := strconv.Atoi(fenSegments[5])
	if err != nil {
		return nil, fmt.Errorf("could not parse fullmove clock - got '%v'", fenSegments[5])
	}
	if fullmoves < 0 {
		return nil, fmt.Errorf("fullmove clock cannot be negative - got '%v'", fullmoves)
	}

	board.fullmoveClock = fullmoves

	return &board, nil
}
