package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/frasmataz/go-chess/api"
	"github.com/frasmataz/go-chess/db"
	"github.com/frasmataz/go-chess/model"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api.StartServer(ctx)
	//runTournament()

}

func runTournament() {

	sim_ctx, sim_cancel := context.WithCancel(context.Background())
	defer sim_cancel()

	err := db.InitDB()
	defer db.CloseDB()

	if err != nil {
		log.Fatalf("Error initialising DB: %v", err)
	}

	t := model.RunTournament(sim_ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				printStatus(t)
			case <-ctx.Done():
				return
			}
		}
	}(sim_ctx)

	err = <-t.Done
	if err != nil {
		log.Fatalf("error running tournament: %v", err)
	}

	printEndResults(t)

}

func printStatus(t *model.Tournament) {
	for _, mu := range t.Matchups {
		log.Printf("--- Matchup %s ----", mu.ID)
		log.Printf(
			"--- WHITE %s | %s BLACK ",
			reflect.TypeOf(mu.White).Name(),
			reflect.TypeOf(mu.Black).Name(),
		)
		log.Printf(
			"--- %d : %d ",
			mu.Results.WhiteWins,
			mu.Results.BlackWins,
		)
		log.Printf(
			"--- %.1f%% complete | %d draws, %d errors",
			(float32(mu.Results.Completed)/float32(mu.Rounds))*100.0,
			mu.Results.Draws,
			mu.Results.Errors,
		)
		log.Println("-----------------------------------------------------")
	}
	log.Println("")
}

func printEndResults(t *model.Tournament) {

	log.Printf("Tournament ID: %s", t.ID)
	log.Printf("Started: %s, Ended %s", t.StartTime.String(), t.EndTime.String())

	for _, mu := range t.Matchups {
		log.Printf("Matchup ID: %s", mu.ID)
		log.Printf(
			"White: %s, Black: %s",
			reflect.TypeOf(mu.White).Name(),
			reflect.TypeOf(mu.Black).Name(),
		)
		log.Printf("Score W/B/D/E: %d:%d:%d:%d", mu.Results.WhiteWins, mu.Results.BlackWins, mu.Results.Draws, mu.Results.Errors)
	}

}
