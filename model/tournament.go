package model

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/frasmataz/go-chess/db"
	"github.com/google/uuid"
)

type Tournament struct {
	ID        uuid.UUID
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
	t.ID = uuid.New()
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

		for _, matchup := range tournament.Matchups {
			mwg.Add(1)
			go matchup.Run(ctx, &mwg, tournament)
		}

		mwg.Wait()

		tournament.EndTime = time.Now()
		tournament.State = DONE

		log.Printf("Tournament finished: %s", tournament.ID.String())

		err := SaveTournament(tournament)
		if err != nil {
			log.Fatalf("Error saving tournament: %v", err)
		}

		close(tournament.Done)

	}(&t)

	return &t
}

func SaveTournament(tournament *Tournament) error {

	sqlStmt := `
		INSERT INTO tournaments (id, state, start_time, end_time)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.SafeExec(
		sqlStmt,
		tournament.ID,
		tournament.State,
		tournament.StartTime,
		tournament.EndTime,
	)
	if err != nil {
		return err
	}

	return nil

}

func SaveTournamentMatchup(tournament *Tournament, matchup *Matchup) error {

	sqlStmt := `
		INSERT INTO tournament_matchup (
			tournament_id,
			matchup_id
		)
		VALUES (?, ?)
	`

	_, err := db.SafeExec(
		sqlStmt,
		tournament.ID,
		matchup.ID,
	)
	if err != nil {
		return err
	}

	return nil

}
