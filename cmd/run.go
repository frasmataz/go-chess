package main

import (
	"context"
	"log"

	"github.com/frasmataz/go-chess/bots"
	"github.com/frasmataz/go-chess/conf"
	"github.com/frasmataz/go-chess/internal"
)

func main() {
	cfg := conf.DefaultConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GameTimeout)
	defer cancel()

	t, err := internal.NewTournament(&ctx, cfg, 100, bots.NewRandomBot(1), bots.NewCheckmateCheckTakeBot(1))
	if err != nil {
		log.Fatal(err)
	}

	t.Run(ctx)
}
