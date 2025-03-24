package internal

import (
	"context"
	"sync"
	"time"

	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/google/uuid"
)

type Tournament struct {
	RunId          uuid.UUID
	StartTime      time.Time
	EndTime        time.Time
	MatchupResults []*MatchupResult
}

var enabledBots = [2]bots.Bot{
	bots.NewRandomBot(),
	bots.NewCheckmateCheckTakeBot(),
}

func RunTournament() Tournament {
	tournament := Tournament{}
	tournament.StartTime = time.Now()
	tournament.RunId = uuid.New()

	cfg := conf.DefaultConfig()
	var matchups []*Matchup

	for _, whiteBot := range enabledBots {
		for _, blackBot := range enabledBots {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.GameTimeout)
			defer cancel()

			matchups = append(matchups, NewMatchup(
				&ctx,
				cfg,
				100,
				whiteBot,
				blackBot,
			))

		}
	}

	var wg sync.WaitGroup
	for _, matchup := range matchups {
		wg.Add(1)
		go matchup.Run(&wg)
	}

	wg.Wait()

	for _, matchup := range matchups {
		tournament.MatchupResults = append(tournament.MatchupResults, &matchup.Result)
	}

	tournament.EndTime = time.Now()
	return tournament
}
