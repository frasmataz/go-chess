package game

type PlayerColour int

const (
	Undefined PlayerColour = iota
	Black
	White
)

func (p *PlayerColour) getOpponentColour() PlayerColour {
	if *p == White {
		return Black
	} else if *p == Black {
		return White
	}

	return Undefined
}
