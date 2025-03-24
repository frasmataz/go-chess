package bots

import (
	"math/rand"
	"time"

	"github.com/corentings/chess"
)

type RandomBot struct {
	seed int
}

func NewRandomBot(seed int) RandomBot {
	b := RandomBot{}
	b.seed = seed
	return b
}

func (b RandomBot) GetMove(game *chess.Game) *chess.Move {

	source := rand.NewSource(time.Now().UnixNano() + int64(b.seed))
	generator := rand.New(source)

	validMoves := game.ValidMoves()
	return validMoves[generator.Intn(len(validMoves))]

}
