package game

import (
	"fmt"
	"regexp"
	"strconv"
)

var positionRegex = regexp.MustCompile(`^[a-h]{1}[1-8]{1}$`)
var uciMoveRegex = regexp.MustCompile(`^[a-h]{1}[1-8]{1}[a-h]{1}[1-8]{1}$`)

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

func validatePosition(position string) bool {
	return positionRegex.MatchString(position)
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

func ParseUCIMoveString(uciMoveString string) (string, string, error) {
	if !uciMoveRegex.MatchString(uciMoveString) {
		return "", "", fmt.Errorf("move did not match regex %v - got '%v'", uciMoveRegex, uciMoveString)
	}

	fromPos := uciMoveString[0:2]
	toPos := uciMoveString[2:4]

	return fromPos, toPos, nil
}
