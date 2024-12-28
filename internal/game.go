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
	enPassantTarget string
	halfmoveClock   int
	fullmoveClock   int
}

func validatePosition(position string) bool {
	positionRegex := `^[a-h]{1}[1-8]{1}$`
	re := regexp.MustCompile(positionRegex)
	return re.MatchString(position)
}

func positionToIndex(position string) ([2]int, error) {
	// Converts standard board position (eg. "e4") to indeces in the [8][8]Piece board state

	if !validatePosition(position) {
		return [2]int{-1, -1}, fmt.Errorf("position string invalid - got '%v'", position)
	}

	colIndex := int(position[0]) - 97 // 97 is int value of rune 'a'
	if colIndex < 0 || colIndex > 7 {
		return [2]int{-1, -1}, fmt.Errorf("parsed row out-of-bounds, should be in range 0-7 - input '%v', parsed to row %v", position, colIndex)
	}

	rowRawInt, err := strconv.Atoi(string(rune(position[1]))) // FIXME: this can't be the best way
	if err != nil {
		return [2]int{-1, -1}, err
	}

	rowIndex := 8 - int(rowRawInt) // board row indeces are top-down, position strings are bottom-up
	if rowIndex < 0 || rowIndex > 7 {
		return [2]int{-1, -1}, fmt.Errorf("parsed column out-of-bounds, should be in range 0-7 - input '%v', parsed to column %v", position, rowIndex)
	}

	return [2]int{rowIndex, colIndex}, nil
}

func IndexToPosition(index [2]int) (string, error) {
	if index[0] < 0 || index[0] > 7 || index[1] < 0 || index[1] > 7 {
		return "", fmt.Errorf("input indeces must both be between 0 - 7 inclusive - got '%v'", index)
	}

	return fmt.Sprintf("%c%c", rune(index[1]+97), rune(8-index[0])), nil
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

				boardColumn += n - 1
			} else {
				// If not digit, should be a jpiece
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
		if validatePosition(fenSegments[3]) {
			board.enPassantTarget = fenSegments[3]
		}
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

func (b *gameState) PrintGameState() string {
	output := "\n\n"

	for _, row := range b.boardState {
		rowString := ""
		for _, column := range row {
			rowString += string(column.Symbol)
		}

		rowString += "\n"
		output += rowString
	}

	output += "\n\n"

	player := PlayerColour(b.nextPlayer)
	if player == White {
		output += "White to play.\n"
	} else if player == Black {
		output += "Black to play.\n"
	} else {
		output += "Game stopped.\n"
	}

	output += fmt.Sprintf("Halfmoves: %v  Fullmoves: %v\n\n", b.halfmoveClock, b.fullmoveClock)
	output += fmt.Sprintf(
		"Castling rights:\nWhite: K:%v Q:%v\nBlack: K:%v Q:%v\n\n",
		b.castlingRights.whiteKingCastle,
		b.castlingRights.whiteQueenCastle,
		b.castlingRights.blackKingCastle,
		b.castlingRights.blackQueenCastle,
	)

	if b.enPassantTarget != "" {
		output += fmt.Sprintf("En passant target: %v", b.enPassantTarget)
	}

	return output
}
