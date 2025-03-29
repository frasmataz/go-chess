package conf

import (
	"time"
)

type DrawLevel int

type Conf struct {
	GameTimeout   time.Duration
	NumberOfGames int
	HTTPAddr      string
}

func DefaultConfig() Conf {
	return Conf{
		GameTimeout:   60 * time.Second,
		NumberOfGames: 100,
		HTTPAddr:      "0.0.0.0:8080",
	}
}
