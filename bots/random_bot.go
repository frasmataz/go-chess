package bots

import (
	"math/rand"

	"github.com/corentings/chess"
)

type RandomBot struct {
}

func NewRandomBot() RandomBot {
	b := RandomBot{}
	return b
}

func (b RandomBot) GetMove(game *chess.Game) *chess.Move {

	validMoves := game.ValidMoves()
	return validMoves[rand.Intn(len(validMoves))]

}
