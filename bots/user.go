package bots

import (
	"fmt"

	"github.com/corentings/chess"
)

type UserBot struct {
}

func (b UserBot) GetMove(game *chess.Game) *chess.Move {
	fmt.Print("Enter a move: ")
	var movestring string
	fmt.Scanln(&movestring)

	alg := chess.AlgebraicNotation{}

	nextMove, err := alg.Decode(game.Position(), movestring)
	if err != nil {
		fmt.Printf("Invalid move: %s", movestring)
	}

	return nextMove
}
