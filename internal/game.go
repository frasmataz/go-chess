package game

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2/log"
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

func indexToPosition(index [2]int) (string, error) {
	if index[0] < 0 || index[0] > 7 || index[1] < 0 || index[1] > 7 {
		return "", fmt.Errorf("input indeces must both be between 0 - 7 inclusive - got '%v'", index)
	}

	return fmt.Sprintf("%c%s", rune(index[1]+97), strconv.Itoa(8-index[0])), nil
}

func positionRelative(position string, offsetX int, offsetY int) (string, error) {
	startingIndex, err := positionToIndex(position)
	if err != nil {
		return "", err
	}

	output, err := indexToPosition([2]int{startingIndex[0] + offsetY, startingIndex[1] + offsetX})
	if err != nil {
		return "", err
	}

	return output, nil
}

func (game *gameState) getPiece(position string) (Piece, error) {
	index, err := positionToIndex(position)
	if err != nil {
		return Piece{}, err
	}

	return game.boardState[index[0]][index[1]], nil
}

func (game *gameState) setPiece(piece Piece, position string) error {
	index, err := positionToIndex(position)
	if err != nil {
		return err
	}

	game.boardState[index[0]][index[1]] = piece
	return nil
}

func NewGame() *gameState {
	game, err := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		panic(err)
	}

	return game
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

func (game *gameState) BoardToFEN() string {
	// Converts a FEN board state string into a board objectgo run
	// https://www.chess.com/terms/fen-chess

	output := ""

	// First segment describes layout of pieces on the board

	for rownum, row := range game.boardState {
		rowString := ""
		spaceCount := 0
		for _, piece := range row {
			if piece.Class == Space {
				spaceCount++
			} else {
				if spaceCount != 0 {
					rowString += strconv.Itoa(spaceCount)
					spaceCount = 0
				}
				rowString += string(piece.FENSymbol)
			}

		}
		if spaceCount != 0 {
			rowString += strconv.Itoa(spaceCount)
			spaceCount = 0
		}
		if rownum < 7 {
			rowString += "/"
		} else {
			rowString += " "
		}
		output += rowString
	}

	// Second segment describes whose turn it is
	if game.nextPlayer == Black {
		output += "b"
	} else {
		output += "w"
	}

	output += " "

	// Third segment describes castling rights
	if !(game.castlingRights.whiteKingCastle || game.castlingRights.whiteQueenCastle || game.castlingRights.blackKingCastle || game.castlingRights.blackQueenCastle) {
		output += "-"
	} else {
		if game.castlingRights.whiteKingCastle {
			output += "K"
		}
		if game.castlingRights.whiteQueenCastle {
			output += "Q"
		}
		if game.castlingRights.blackKingCastle {
			output += "k"
		}
		if game.castlingRights.blackQueenCastle {
			output += "q"
		}
	}

	output += " "

	// Fourth segment describes en passant targets
	if game.enPassantTarget == "" {
		output += "-"
	} else {
		output += game.enPassantTarget
	}

	output += " "

	// Fifth segment counts halfmoves
	output += strconv.Itoa(game.halfmoveClock)
	output += " "

	// Sixth segment counts halfmoves
	output += strconv.Itoa(game.fullmoveClock)

	return output
}

func (b *gameState) PrintGameState() string {
	blackSquare := color.BgRGB(0, 0, 0)
	whiteSquare := color.BgRGB(60, 50, 80)

	output := "\n\n"

	for ri, row := range b.boardState {
		rowString := fmt.Sprintf("%s ", strconv.Itoa(8-ri))
		for ci, column := range row {
			squareColour := blackSquare
			if ri%2 == ci%2 {
				squareColour = whiteSquare
			}
			rowString += squareColour.Sprint(string(column.Symbol) + " ")
		}

		rowString += "\n"
		output += rowString
	}

	output += "  a b c d e f g h \n\n"

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

func (game *gameState) isSpaceValidAndEmpty(position string) bool {
	// This is a common operation - extracted to helper function to keep move code cleaner and reduce error checks
	p, err := game.getPiece(position)
	if err != nil {
		return false
	}
	if p.Class == Space {
		return true
	}
	return false
}

func (game *gameState) isSpaceValidAndOpponentPiece(position string, opponentColor PlayerColour) bool {
	// This is a common operation - extracted to helper function to keep move code cleaner and reduce error checks
	p, err := game.getPiece(position)
	if err != nil {
		return false
	}
	if p.Colour == opponentColor {
		return true
	}
	return false
}

func (game *gameState) GetValidMovesForPiece(position string) ([]string, error) {
	piece, err := game.getPiece(position)
	if err != nil {
		return nil, err
	}

	output := []string{}

	var opponentColour PlayerColour
	if piece.Colour == White {
		opponentColour = Black
	} else if piece.Colour == Black {
		opponentColour = White
	}

	switch piece.Class {
	case Pawn:
		// Check one step forward
		yOffset := 1
		if piece.Colour == White {
			yOffset = -1
		}
		target, err := positionRelative(position, 0, yOffset)
		if err != nil {
			log.Warn(err)
		}
		// Can only move there if space is empty
		if game.isSpaceValidAndEmpty(target) {
			output = append(output, target)

			// Also check two steps forward if still on starting row
			currentIndex, err := positionToIndex(position)
			if err != nil {
				return nil, fmt.Errorf("piece in unexpected position - pos: %v, err: %v", position, err)
			}

			// Get starting row for player
			startingRow := 1
			if piece.Colour == White {
				startingRow = 6
			}

			if currentIndex[0] == startingRow {
				// Get space two ahead
				yOffset := 2
				if piece.Colour == White {
					yOffset = -2
				}

				target, err := positionRelative(position, 0, yOffset)
				if err != nil {
					return nil, fmt.Errorf("pawn is on starting row but can't move forward by 2, this shouldn't happen - pos: '%v', target: '%v'", position, target)
				}

				// Can only move there if space if empty
				if game.isSpaceValidAndEmpty(target) {
					output = append(output, target)
				}
			}
		}
		// Check left attack
		yOffset = 1
		if piece.Colour == White {
			yOffset = -1
		}

		target, err = positionRelative(position, -1, yOffset)
		if err == nil {
			if game.isSpaceValidAndOpponentPiece(target, opponentColour) {
				output = append(output, target)
			}
		}

		// Check right attack
		yOffset = 1
		if piece.Colour == White {
			yOffset = -1
		}

		target, err = positionRelative(position, 1, yOffset)
		if err == nil {
			if game.isSpaceValidAndOpponentPiece(target, opponentColour) {
				output = append(output, target)
			}
		}

		//TODO: En passant capture
	}

	return output, nil
}

func (game *gameState) ExecuteMove(move string) error {
	moveRegex := `^[a-h]{1}[1-8]{1}[a-h]{1}[1-8]{1}$`
	re := regexp.MustCompile(moveRegex)
	if !re.MatchString(move) {
		return fmt.Errorf("move did not match regex %v - got '%v'", moveRegex, move)
	}

	moveFrom := move[0:2]
	moveTo := move[2:4]

	piece, err := game.getPiece(moveFrom)
	if err != nil {
		return err
	}

	validMoves, err := game.GetValidMovesForPiece(moveFrom)
	if err != nil {
		log.Error(err)
	}

	if slices.Contains(validMoves, moveTo) {
		game.setPiece(Pieces["space"], moveFrom)
		game.setPiece(piece, moveTo)
	} else {
		log.Warn("Invalid move")
		log.Warnf("%v", validMoves)
	}

	return nil
}
