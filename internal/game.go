package internal

import (
	"context"
	"time"

	"github.com/corentings/chess"
	"github.com/frasmataz/go-chess/bots"
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

func (g *Game) Run(ctx context.Context) {

	g.State = RUNNING
	g.StartTime = time.Now()

	for g.Game.Outcome() == chess.NoOutcome {
		select {
		case <-ctx.Done():
			g.Game.Draw(chess.DrawOffer)
			g.Finish()
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

	g.Finish()

}

func (g *Game) Finish() {
	g.State = DONE
	g.EndTime = time.Now()
	close(g.Done)
}
