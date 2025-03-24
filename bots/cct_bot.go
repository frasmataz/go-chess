package bots

import (
	"math/rand"

	"github.com/corentings/chess"
)

type CheckmateCheckTakeBot struct {
}

func NewCheckmateCheckTakeBot() CheckmateCheckTakeBot {
	b := CheckmateCheckTakeBot{}
	return b
}

func (b CheckmateCheckTakeBot) GetMove(game *chess.Game) *chess.Move {

	validMoves := game.ValidMoves()

	// Loop through moves looking for checks and takes.
	// Return immediately if you find a checkmate
	var checks []*chess.Move
	var takes []*chess.Move

	for _, candidate := range game.ValidMoves() {
		if candidate.HasTag(chess.Check) {
			cloneGame := game.Clone()
			cloneGame.Move(candidate)
			if cloneGame.Outcome() != chess.NoOutcome && cloneGame.Method() == chess.Checkmate {
				return candidate
			}

			if cloneGame.Position().Board().Piece(candidate.S2()) != chess.NoPiece {
				takes = append(takes, candidate)
			}

			checks = append(checks, candidate)
		}
	}

	// Return a random check
	if len(checks) > 0 {
		return checks[rand.Intn(len(checks))]
	}

	// Return a random take
	if len(takes) > 0 {
		return takes[rand.Intn(len(takes))]
	}

	return validMoves[rand.Intn(len(validMoves))]

}
