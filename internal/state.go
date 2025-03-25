package internal

type State int

const (
	INIT    State = 0
	RUNNING State = 1
	DONE    State = 2
)
