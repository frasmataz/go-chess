package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/corentings/chess"
)

var game *chess.Game

func main() {
	fmt.Print("Enter starting FEN, or leave blank for new game: ")

	in := bufio.NewReader(os.Stdin)
	fenstring, err := in.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading string: %s", err)
	}

	game = chess.NewGame()

	fenstring = strings.TrimSuffix(fenstring, "\n")
	if fenstring != "" {
		fen, err := chess.FEN(fenstring)
		if err != nil {
			log.Fatalf("Invalid FEN: %s", fenstring)
		}

		game = chess.NewGame(fen)
	}

	for game.Outcome() == chess.NoOutcome {
		log.Println(game.Position().Board().Draw())
		log.Print("Enter a move: ")

		var movestring string
		fmt.Scanln(&movestring)

		err := game.MoveStr(movestring)
		if err != nil {
			log.Printf("Invalid move: %s", movestring)
		}
	}
}
