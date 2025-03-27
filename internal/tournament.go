package internal

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/google/uuid"
)

type Tournament struct {
	RunId     uuid.UUID
	State     State
	StartTime time.Time
	EndTime   time.Time
	Matchups  []*Matchup
	Done      chan error
}

var enabledBots = [2]bots.Bot{
	bots.NewRandomBot(),
	bots.NewCheckmateCheckTakeBot(),
}

func RunTournament(ctx context.Context) *Tournament {

	cfg := conf.DefaultConfig()

	t := Tournament{}

	t.StartTime = time.Now()
	t.RunId = uuid.New()
	t.Done = make(chan error)

	for _, whiteBot := range enabledBots {
		for _, blackBot := range enabledBots {

			t.Matchups = append(t.Matchups, NewMatchup(
				ctx,
				cfg,
				cfg.NumberOfGames,
				whiteBot,
				blackBot,
			))

		}
	}

	t.State = RUNNING

	go func(tournament *Tournament) {

		var mwg sync.WaitGroup

		for _, matchup := range t.Matchups {
			mwg.Add(1)
			go matchup.Run(ctx, &mwg)
		}

		mwg.Wait()

		tournament.EndTime = time.Now()
		tournament.State = DONE

		log.Printf("Tournament finished: %s", tournament.RunId.String())
		close(tournament.Done)

	}(&t)

	return &t
}
