package game

import (
	"fmt"
	"regexp"
)

type move struct {
	from string
	to   string
}

func UCIToMove(uciMoveString string) (*move, error) {
	moveRegex := `^[a-h]{1}[1-8]{1}[a-h]{1}[1-8]{1}$`
	re := regexp.MustCompile(moveRegex)
	if !re.MatchString(uciMoveString) {
		return nil, fmt.Errorf("move did not match regex %v - got '%v'", moveRegex, uciMoveString)
	}

	move := move{
		from: uciMoveString[0:2],
		to:   uciMoveString[2:4],
	}

	return &move, nil
}
