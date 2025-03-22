package bots

import "github.com/corentings/chess"

type Bot interface {
	GetMove(*chess.Game) *chess.Move
}
