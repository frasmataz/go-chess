package main

import (
	"log"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
)

var game *chess.Game
var players map[chess.Color]bots.Bot

func main() {
	game = chess.NewGame()

	players := map[chess.Color]bots.Bot{
		chess.White: bots.UserBot{},
		chess.Black: bots.RandomBot{},
	}

	for game.Outcome() == chess.NoOutcome {
		log.Println(game.Position().Board().Draw())

		player := players[game.Position().Turn()]
		nextMove := player.GetMove(game)
		game.Move(nextMove)

	}

	log.Println(game.Position().Board().Draw())
	log.Printf("Game complete.  Outcome: %s, by %s", game.Outcome(), game.Method())
	log.Printf("Moves: %s", game.Moves())
	log.Printf("End position: %s", game.Position().String())

}
