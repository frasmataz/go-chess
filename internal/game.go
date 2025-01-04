package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatih/color"
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

type GameState struct {
	BoardState      [8][8]*Piece // [7][0] is a1, [0][7] is h8 - [0][0] is top-left from white's perspective
	NextPlayer      PlayerColour
	CastlingRights  castlingRights
	EnPassantTarget string
	HalfmoveClock   int
	FullmoveClock   int
	ValidMoves      map[PlayerColour][]move
}

var positionRegex = regexp.MustCompile(`^[a-h]{1}[1-8]{1}$`)

func NewGame() *GameState {
	game, err := NewGameFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		panic(err)
	}

	return game
}

func NewGameFromFEN(fen string) (*GameState, error) {
	// Converts a FEN board state string into a board objectgo run
	// https://www.chess.com/terms/fen-chess

	fenSegments := strings.Split(fen, " ")

	// First segment describes layout of pieces on the board
	rows := strings.Split(fenSegments[0], "/")
	board := GameState{}

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
					board.BoardState[rownum][boardColumn+i] = GetPiece("space")
				}

				boardColumn += n - 1
			} else {
				// If not digit, should be a jpiece
				piece, err := FENToPiece(char)
				if err != nil {
					return nil, err
				}

				board.BoardState[rownum][boardColumn] = piece
			}

			boardColumn++
		}
	}

	// Second segment describes whose turn it is
	if fenSegments[1] == "w" {
		board.NextPlayer = White
	} else if fenSegments[1] == "b" {
		board.NextPlayer = Black
	} else {
		return nil, fmt.Errorf("next player section missing from FEN")
	}

	// Third segment describes castling rights
	if fenSegments[2] == "-" {
		board.CastlingRights = castlingRights{
			blackKingCastle:  false,
			blackQueenCastle: false,
			whiteKingCastle:  false,
			whiteQueenCastle: false,
		}
	} else {
		if strings.Contains(fenSegments[2], "k") {
			board.CastlingRights.blackKingCastle = true
		}

		if strings.Contains(fenSegments[2], "q") {
			board.CastlingRights.blackQueenCastle = true
		}

		if strings.Contains(fenSegments[2], "K") {
			board.CastlingRights.whiteKingCastle = true
		}

		if strings.Contains(fenSegments[2], "Q") {
			board.CastlingRights.whiteQueenCastle = true
		}
	}

	// Fourth segment describes en passant targets
	if fenSegments[3] != "-" {
		if validatePosition(fenSegments[3]) {
			board.EnPassantTarget = fenSegments[3]
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

	board.HalfmoveClock = halfmoves

	// Sixth segment counts halfmoves
	fullmoves, err := strconv.Atoi(fenSegments[5])
	if err != nil {
		return nil, fmt.Errorf("could not parse fullmove clock - got '%v'", fenSegments[5])
	}
	if fullmoves < 0 {
		return nil, fmt.Errorf("fullmove clock cannot be negative - got '%v'", fullmoves)
	}

	board.FullmoveClock = fullmoves
	board.updateValidMoves()

	return &board, nil
}

func validatePosition(position string) bool {
	return positionRegex.MatchString(position)
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

func TryApplyMove(game GameState, uciMoveString string) (*GameState, error) {
	moveFromPos, moveToPos, err := ParseUCIMoveString(uciMoveString)
	if err != nil {
		return nil, fmt.Errorf("error parsing move: %v", err)
	}

	var move *move

	for _, m := range game.ValidMoves[game.NextPlayer] {
		if m.fromPos == moveFromPos && m.toPos == moveToPos {
			move = &m
			break
		}
	}

	if move == nil {
		return nil, fmt.Errorf("invalid move '%v'\nvalid moves for piece:\n%v", uciMoveString, game.ValidMoves)
	}

	if moveCausesCheck(game, *move) {
		return nil, fmt.Errorf("move '%v' causes check\n\nstate: %v\nvalid moves: %v", uciMoveString, game.PrintGameState(), game.ValidMoves)
	}

	newState := ApplyMove(game, *move)
	return newState, nil

}

func ApplyMove(game GameState, move move) *GameState {
	game.setPiece(GetPiece("space"), move.fromPos)
	game.setPiece(move.fromPiece, move.toPos)

	game.HalfmoveClock++
	if move.isThreat || move.fromPiece.Class == Pawn {
		game.HalfmoveClock = 0
	}

	if game.NextPlayer == Black {
		game.FullmoveClock++
	}

	game.NextPlayer = game.NextPlayer.getOpponentColour()
	game.updateCastlingRights()
	game.updateValidMoves()

	return &game
}

func (game *GameState) ToFEN() string {
	// Converts a FEN board state string into a board objectgo run
	// https://www.chess.com/terms/fen-chess

	output := ""

	// First segment describes layout of pieces on the board

	for rownum, row := range game.BoardState {
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
	if game.NextPlayer == Black {
		output += "b"
	} else {
		output += "w"
	}

	output += " "

	// Third segment describes castling rights
	if !(game.CastlingRights.whiteKingCastle || game.CastlingRights.whiteQueenCastle || game.CastlingRights.blackKingCastle || game.CastlingRights.blackQueenCastle) {
		output += "-"
	} else {
		if game.CastlingRights.whiteKingCastle {
			output += "K"
		}
		if game.CastlingRights.whiteQueenCastle {
			output += "Q"
		}
		if game.CastlingRights.blackKingCastle {
			output += "k"
		}
		if game.CastlingRights.blackQueenCastle {
			output += "q"
		}
	}

	output += " "

	// Fourth segment describes en passant targets
	if game.EnPassantTarget == "" {
		output += "-"
	} else {
		output += game.EnPassantTarget
	}

	output += " "

	// Fifth segment counts halfmoves
	output += strconv.Itoa(game.HalfmoveClock)
	output += " "

	// Sixth segment counts halfmoves
	output += strconv.Itoa(game.FullmoveClock)

	return output
}

func (game *GameState) PrintGameState() string {
	blackSquare := color.BgRGB(0, 0, 0)
	whiteSquare := color.BgRGB(60, 50, 80)

	output := "\n\n"

	for ri, row := range game.BoardState {
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

	player := PlayerColour(game.NextPlayer)
	if player == White {
		output += "White to play.\n"
	} else if player == Black {
		output += "Black to play.\n"
	} else {
		output += "Game stopped.\n"
	}

	output += fmt.Sprintf("Halfmoves: %v  Fullmoves: %v\n\n", game.HalfmoveClock, game.FullmoveClock)
	output += fmt.Sprintf(
		"Castling rights:\nWhite: K:%v Q:%v\nBlack: K:%v Q:%v\n\n",
		game.CastlingRights.whiteKingCastle,
		game.CastlingRights.whiteQueenCastle,
		game.CastlingRights.blackKingCastle,
		game.CastlingRights.blackQueenCastle,
	)

	if game.EnPassantTarget != "" {
		output += fmt.Sprintf("En passant target: %v", game.EnPassantTarget)
	}

	return output
}

func (game *GameState) getPiece(position string) (*Piece, error) {
	index, err := positionToIndex(position)
	if err != nil {
		return nil, err
	}

	return game.BoardState[index[0]][index[1]], nil
}

func (game *GameState) getPieceSafe(position string) *Piece {
	index, err := positionToIndex(position)
	if err != nil {
		panic(err)
	}

	return game.BoardState[index[0]][index[1]]
}

func (game *GameState) setPiece(piece *Piece, position string) error {
	index, err := positionToIndex(position)
	if err != nil {
		return err
	}

	game.BoardState[index[0]][index[1]] = piece
	return nil
}

func (game *GameState) isSpaceEmpty(position string) bool {
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

func (game *GameState) isSpacePlayersPiece(position string, player PlayerColour) bool {
	// This is a common operation - extracted to helper function to keep move code cleaner and reduce error checks
	p, err := game.getPiece(position)
	if err != nil {
		return false
	}
	if p.Colour == player {
		return true
	}
	return false
}

func (game *GameState) updateCastlingRights() {
	if game.CastlingRights.whiteKingCastle {
		kingStartSpace := game.getPieceSafe("e1")
		rookStartSpace := game.getPieceSafe("h1")
		if !(kingStartSpace.Class == King && kingStartSpace.Colour == White) || !(rookStartSpace.Class == Rook && rookStartSpace.Colour == White) {
			game.CastlingRights.whiteKingCastle = false
		}
	}

	if game.CastlingRights.whiteQueenCastle {
		kingStartSpace := game.getPieceSafe("e1")
		rookStartSpace := game.getPieceSafe("a1")
		if !(kingStartSpace.Class == King && kingStartSpace.Colour == White) || !(rookStartSpace.Class == Rook && rookStartSpace.Colour == White) {
			game.CastlingRights.whiteQueenCastle = false
		}
	}

	if game.CastlingRights.blackKingCastle {
		kingStartSpace := game.getPieceSafe("e8")
		rookStartSpace := game.getPieceSafe("h8")
		if !(kingStartSpace.Class == King && kingStartSpace.Colour == Black) || !(rookStartSpace.Class == Rook && rookStartSpace.Colour == Black) {
			game.CastlingRights.blackKingCastle = false
		}
	}

	if game.CastlingRights.blackQueenCastle {
		kingStartSpace := game.getPieceSafe("e8")
		rookStartSpace := game.getPieceSafe("a8")
		if !(kingStartSpace.Class == King && kingStartSpace.Colour == Black) || !(rookStartSpace.Class == Rook && rookStartSpace.Colour == Black) {
			game.CastlingRights.blackQueenCastle = false
		}
	}
}

func (game *GameState) updateValidMoves() {
	game.ValidMoves = map[PlayerColour][]move{
		White: GetValidMovesForPlayer(game, White),
		Black: GetValidMovesForPlayer(game, Black),
	}
}

// func removeMovesCausingSelfCheck(game *GameState, moves []*move) bool {
// 	for move := range moves {

// 	}
// }
