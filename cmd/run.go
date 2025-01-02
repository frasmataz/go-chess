package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/log"

	game "github.com/frasmataz/go-chess/internal"
)

func main() {
	fmt.Print("Enter starting FEN, or leave blank for new game: ")

	in := bufio.NewReader(os.Stdin)
	fen, err := in.ReadString('\n')
	if err != nil {
		panic(err)
	}

	gameState := game.NewGame()

	if fen != "" {
		fen = strings.TrimSuffix(fen, "\n")
		gameState, err = game.BoardFromFEN(fen)
		if err != nil {
			panic(err)
		}
	}

	for {
		fmt.Print(gameState.PrintGameState())
		fmt.Println(gameState.BoardToFEN())
		fmt.Print("Enter a move: ")

		var move string
		fmt.Scanln(&move)

		err := gameState.ExecuteMove(move)
		if err != nil {
			log.Error(err)
		}
	}
}
