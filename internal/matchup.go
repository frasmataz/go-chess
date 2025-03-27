package internal

import (
	"context"
	"sync"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/google/uuid"
)

type Matchup struct {
	config  conf.Conf
	ID      uuid.UUID
	State   State
	Rounds  int
	White   bots.Bot
	Black   bots.Bot
	Games   []Game
	Results MatchupResults
}

type MatchupResults struct {
	Completed int
	BlackWins int
	WhiteWins int
	Draws     int
	Errors    int
}

func NewMatchup(ctx context.Context, config conf.Conf, rounds int, white bots.Bot, black bots.Bot) *Matchup {

	m := Matchup{
		config:  config,
		ID:      uuid.New(),
		State:   INIT,
		Rounds:  rounds,
		White:   white,
		Black:   black,
		Results: MatchupResults{},
	}

	for range rounds {
		m.Games = append(m.Games, NewGame(white, black))
	}

	return &m

}

func (m *Matchup) Run(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	m.State = RUNNING

	// Run all games concurrently
	var game_wg sync.WaitGroup

	for _, game := range m.Games {

		game_wg.Add(1)

		go func(g *Game, game_wg *sync.WaitGroup) {

			defer game_wg.Done()

			subCtx, cancel := context.WithTimeout(ctx, m.config.GameTimeout)
			defer cancel()

			g.Run(subCtx)

			m.Results.Completed++
			switch g.Game.Outcome() {
			case chess.BlackWon:
				m.Results.BlackWins++
			case chess.WhiteWon:
				m.Results.WhiteWins++
			case chess.Draw:
				m.Results.Draws++
			}

		}(&game, &game_wg)

	}

	game_wg.Wait()

	m.State = DONE

}
