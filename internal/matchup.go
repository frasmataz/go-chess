package internal

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/frasmataz/go-chess/db"
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

func (m *Matchup) Run(ctx context.Context, wg *sync.WaitGroup, parentTournament *Tournament) {

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

			g.Run(subCtx, m)

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

	err := SaveMatchup(m)
	if err != nil {
		log.Fatalf("Error saving matchup: %v", err)
	}

	err = SaveTournamentMatchup(parentTournament, m)
	if err != nil {
		log.Fatalf("Error saving tournament_matchup: %v", err)
	}

}

func SaveMatchup(matchup *Matchup) error {

	sqlStmt := `
		INSERT INTO matchups (
			id,
			state,
			white,
			black,
			rounds,
			completed,
			black_wins,
			white_wins,
			draws,
			errors
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	db, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("Error saving matchup: %v", err)
	}

	_, err = db.Exec(
		sqlStmt,
		matchup.ID,
		matchup.State,
		reflect.TypeOf(matchup.White).Name(),
		reflect.TypeOf(matchup.Black).Name(),
		matchup.Rounds,
		matchup.Results.Completed,
		matchup.Results.BlackWins,
		matchup.Results.WhiteWins,
		matchup.Results.Draws,
		matchup.Results.Errors,
	)
	if err != nil {
		return err
	}

	return nil

}

func SaveMatchupGame(matchup *Matchup, game *Game) error {

	sqlStmt := `
		INSERT INTO matchup_game (
			matchup_id,
			game_id
		)
		VALUES (?, ?)
	`

	db, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("Error saving matchup_game: %v", err)
	}

	_, err = db.Exec(
		sqlStmt,
		matchup.ID,
		game.ID,
	)
	if err != nil {
		return err
	}

	return nil

}
