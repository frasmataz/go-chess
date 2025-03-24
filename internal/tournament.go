package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/google/uuid"
)

type Tournament struct {
	id     uuid.UUID
	ctx    *context.Context
	config conf.Conf
	rounds int
	white  bots.Bot
	black  bots.Bot
}

type Result struct {
	game *chess.Game
	err  error
}

func NewTournament(ctx *context.Context, config conf.Conf, rounds int, white bots.Bot, black bots.Bot) (*Tournament, error) {

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Tournament{
		id:     id,
		ctx:    ctx,
		config: config,
		rounds: rounds,
		white:  white,
		black:  black,
	}, nil

}

func (t *Tournament) Run(ctx context.Context) error {

	var results chan Result

	for range t.rounds {
		go func(t *Tournament) {
			game, err := t.simulate(*t.ctx)
			if err != nil {
				results <- Result{err: err}
			}

			results <- Result{game: game}
		}(t)
	}

	for range t.rounds {
		result := <-results
		if result.err != nil {
			return result.err
		}

		printResults(t.config, result.game)
	}

	return nil
}

func (t *Tournament) simulate(ctx context.Context) (*chess.Game, error) {

	game := chess.NewGame()

	for game.Outcome() == chess.NoOutcome {
		select {
		case <-ctx.Done():
			game.Draw(chess.DrawOffer)
			return game, nil
		default:
			if t.config.DrawLevel >= conf.ALL {
				log.Println(game.Position().Board().Draw())
			}

			switch game.Position().Turn() {
			case chess.Black:
				game.Move(t.black.GetMove(game))
			case chess.White:
				game.Move(t.white.GetMove(game))
			}
		}
	}

	return game, nil

}

func printResults(cfg conf.Conf, game *chess.Game) {

	if cfg.DrawLevel >= conf.RESULT {
		log.Println(game.Position().Board().Draw())
	}

	if cfg.DrawLevel >= conf.ALL {
		for i, move := range game.Moves() {
			fmt.Print(move)
			if i%2 == 1 {
				fmt.Println()
			} else {
				fmt.Print(" ")
			}
		}
	}

	if cfg.DrawLevel >= conf.RESULT {
		log.Printf("Game %d complete.  Outcome: %s, by %s", game.Outcome(), game.Method())
		log.Printf("End position: %s", game.Position().String())
	}

}
