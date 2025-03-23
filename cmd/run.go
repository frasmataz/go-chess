package main

import (
	"fmt"
	"log"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
)

var game *chess.Game
var players map[chess.Color]bots.Bot

func main() {
	game = chess.NewGame()

	players := map[chess.Color]bots.Bot{
		chess.White: bots.RandomBot{},
		chess.Black: bots.CheckmateCheckTakeBot{},
	}

	for game.Outcome() == chess.NoOutcome {
		log.Println(game.Position().Board().Draw())

		player := players[game.Position().Turn()]
		nextMove := player.GetMove(game)
		game.Move(nextMove)
	}

	log.Println(game.Position().Board().Draw())

	for i, move := range game.Moves() {
		fmt.Print(move)
		if i%2 == 1 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
	}
	log.Printf("Game complete.  Outcome: %s, by %s", game.Outcome(), game.Method())
	log.Printf("End position: %s", game.Position().String())

}
