package conf

import (
	"time"
)

type DrawLevel int

const (
	NONE   DrawLevel = 0
	RESULT DrawLevel = 1
	ALL    DrawLevel = 2
)

type Conf struct {
	GameTimeout   time.Duration
	DrawLevel     DrawLevel
	NumberOfGames int
}

func DefaultConfig() Conf {
	return Conf{
		GameTimeout:   10 * time.Second,
		DrawLevel:     NONE,
		NumberOfGames: 20,
	}
}
