package main

import (
	"log"
	"reflect"

	"github.com/frasmataz/go-chess/internal"
)

func main() {
	tournament := internal.RunTournament()

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
