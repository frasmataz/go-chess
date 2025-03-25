package main

import (
	"context"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/frasmataz/go-chess/internal"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		tournament := internal.RunTournament()
	}()

	ctx := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				printStatus()
			}
		}
	}(ctx)

	log.Printf("Tournament ID: %s", tournament.RunId)
	log.Printf("Started: %s, Ended %s", tournament.StartTime.String(), tournament.EndTime.String())

	for _, mr := range tournament.MatchupResults {
		log.Printf("Matchup ID: %s", mr.Matchup.ID)
		log.Printf(
			"White: %s, Black: %s",
			reflect.TypeOf(mr.Matchup.White).Name(),
			reflect.TypeOf(mr.Matchup.Black).Name(),
		)
		log.Printf("Score W/B/D/E: %d:%d:%d:%d", mr.WhiteWins, mr.BlackWins, mr.Draws, mr.Errors)
	}
}

func printStatus() {

}
