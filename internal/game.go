package internal

import (
	"context"
	"log"
	"time"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/db"
	"github.com/google/uuid"
)

type Game struct {
	Game      *chess.Game
	ID        uuid.UUID
	State     State
	StartTime time.Time
	EndTime   time.Time
	White     bots.Bot
	Black     bots.Bot
	Done      chan error
}

func NewGame(white bots.Bot, black bots.Bot) Game {
	return Game{
		Game:  chess.NewGame(),
		ID:    uuid.New(),
		State: INIT,
		White: white,
		Black: black,
		Done:  make(chan error),
	}
}

func (g *Game) Run(ctx context.Context, parentMatchup *Matchup) {

	g.State = RUNNING
	g.StartTime = time.Now()

	for g.Game.Outcome() == chess.NoOutcome {
		select {
		case <-ctx.Done():
			g.Game.Draw(chess.DrawOffer)
			g.Finish(parentMatchup)
			return
		default:
			switch g.Game.Position().Turn() {
			case chess.Black:
				g.Game.Move(g.Black.GetMove(g.Game))
			case chess.White:
				g.Game.Move(g.White.GetMove(g.Game))
			}
		}
	}

	g.Finish(parentMatchup)

}

func (g *Game) Finish(parentMatchup *Matchup) {
	g.State = DONE
	g.EndTime = time.Now()

	err := SaveGame(g)
	if err != nil {
		log.Fatalf("Error saving game: %v", err)
	}

	err = SaveMatchupGame(parentMatchup, g)
	if err != nil {
		log.Fatalf("Error saving matchup_game: %v", err)
	}

	close(g.Done)
}

func SaveGame(game *Game) error {

	sqlStmt := `
		INSERT INTO games (
			id,
			state,
			start_time,
			end_time,
			moves,
			outcome,
			method
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.SafeExec(
		sqlStmt,
		game.ID,
		game.State,
		game.StartTime,
		game.EndTime,
		game.Game.String(),
		game.Game.Outcome().String(),
		game.Game.Method().String(),
	)
	if err != nil {
		return err
	}

	return nil

}
