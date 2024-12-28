package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"

	game "github.com/frasmataz/go-chess/internal"
)

func main() {
	game := game.NewGame()

	for {
		fmt.Print(game.PrintGameState())
		fmt.Print("Enter a move: ")

		var input string
		fmt.Scanln(&input)

		err := game.ExecuteMove(input)
		if err != nil {
			log.Error(err)
		}
	}
}
