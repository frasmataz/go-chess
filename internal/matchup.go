package internal

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/google/uuid"
)

type Matchup struct {
	ID     uuid.UUID
	State  State
	ctx    *context.Context
	config conf.Conf
	Rounds int
	White  bots.Bot
	Black  bots.Bot
	Result MatchupResult
}

type GameResult struct {
	Game *chess.Game
	Err  error
}

type MatchupResult struct {
	Matchup     *Matchup
	BlackWins   int
	WhiteWins   int
	Draws       int
	Errors      int
	GameResults []*GameResult
}

func NewMatchup(ctx *context.Context, config conf.Conf, rounds int, white bots.Bot, black bots.Bot) *Matchup {

	id := uuid.New()

	return &Matchup{
		ID:     id,
		ctx:    ctx,
		config: config,
		Rounds: rounds,
		White:  white,
		Black:  black,
	}

}

func (m *Matchup) Run(wg *sync.WaitGroup) {

	defer wg.Done()
	var matchupResult MatchupResult
	matchupResult.Matchup = m

	gameResults := make(chan *GameResult, 1)

	m.State = RUNNING

	// Run all games concurrently
	for range m.Rounds {
		go func() {
			game, err := m.simulate()
			if err != nil {
				gameResults <- &GameResult{Err: err}
			}

			gameResults <- &GameResult{Game: game}
		}()
	}

	// Process results
	for range m.Rounds {
		result := <-gameResults
		if result.Err != nil {
			matchupResult.Errors++
		}

		switch result.Game.Outcome() {
		case chess.BlackWon:
			matchupResult.BlackWins++
		case chess.WhiteWon:
			matchupResult.WhiteWins++
		case chess.Draw:
			matchupResult.Draws++
		}

		matchupResult.GameResults = append(matchupResult.GameResults, result)
		printResults(m.config, result.Game)
	}

	m.State = DONE
	m.Result = matchupResult

}

func (m *Matchup) simulate() (*chess.Game, error) {

	game := chess.NewGame()

	for game.Outcome() == chess.NoOutcome {
		select {
		case <-(*m.ctx).Done():
			game.Draw(chess.DrawOffer)
			return game, nil
		default:
			if m.config.DrawLevel >= conf.ALL {
				log.Println(game.Position().Board().Draw())
			}

			switch game.Position().Turn() {
			case chess.Black:
				game.Move(m.Black.GetMove(game))
			case chess.White:
				game.Move(m.White.GetMove(game))
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
