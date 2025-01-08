package game

import (
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v2/log"
)

type move struct {
	fromPos   string
	toPos     string
	fromPiece *Piece
	toPiece   *Piece
	isThreat  bool
}

func GetValidMovesForPiece(game *GameState, position string) ([]move, error) {
	fromPiece, err := game.getPiece(position)
	if err != nil {
		return nil, err
	}

	possibleMoves := []move{}
	opponentColour := fromPiece.Colour.getOpponentColour()

	switch fromPiece.Class {
	case Pawn:
		// Check one step forward
		yOffset := 1
		if fromPiece.Colour == White {
			yOffset = -1
		}
		target, err := positionRelative(position, 0, yOffset)
		if err != nil {
			log.Warn(err)
		}
		// Can only move there if space is empty
		if game.isSpaceEmpty(target) {
			possibleMoves = append(
				possibleMoves,
				move{
					fromPos:   position,
					toPos:     target,
					fromPiece: game.getPieceSafe(position),
					toPiece:   game.getPieceSafe(target),
					isThreat:  false,
				},
			)

			// Also check two steps forward if still on starting row
			currentIndex, err := positionToIndex(position)
			if err != nil {
				return nil, fmt.Errorf("piece in unexpected position - pos: %v, err: %v", position, err)
			}

			// Get starting row for player
			startingRow := 1
			if fromPiece.Colour == White {
				startingRow = 6
			}

			if currentIndex[0] == startingRow {
				// Get space two ahead
				yOffset := 2
				if fromPiece.Colour == White {
					yOffset = -2
				}

				target, err := positionRelative(position, 0, yOffset)
				if err != nil {
					return nil, fmt.Errorf("pawn is on starting row but can't move forward by 2, this shouldn't happen - pos: '%v', target: '%v'", position, target)
				}

				// Can only move there if space if empty
				if game.isSpaceEmpty(target) {
					possibleMoves = append(
						possibleMoves,
						move{
							fromPos:   position,
							toPos:     target,
							fromPiece: game.getPieceSafe(position),
							toPiece:   game.getPieceSafe(target),
							isThreat:  false,
						},
					)
				}
			}
		}
		// Check left attack
		yOffset = 1
		if fromPiece.Colour == White {
			yOffset = -1
		}

		target, err = positionRelative(position, -1, yOffset)
		if err == nil {
			if game.isSpacePlayersPiece(target, opponentColour) {
				possibleMoves = append(
					possibleMoves,
					move{
						fromPos:   position,
						toPos:     target,
						fromPiece: game.getPieceSafe(position),
						toPiece:   game.getPieceSafe(target),
						isThreat:  true,
					},
				)
			}
		}

		// Check right attack
		yOffset = 1
		if fromPiece.Colour == White {
			yOffset = -1
		}

		target, err = positionRelative(position, 1, yOffset)
		if err == nil {
			if game.isSpacePlayersPiece(target, opponentColour) {
				possibleMoves = append(
					possibleMoves,
					move{
						fromPos:   position,
						toPos:     target,
						fromPiece: game.getPieceSafe(position),
						toPiece:   game.getPieceSafe(target),
						isThreat:  true,
					},
				)
			}
		}
		//TODO: En passant capture
	case Knight:
		target_offsets := [][2]int{
			{1, 2},
			{2, 1},
			{-1, 2},
			{2, -1},
			{1, -2},
			{-2, 1},
			{-1, -2},
			{-2, -1},
		}

		for _, offset := range target_offsets {
			target, err := positionRelative(position, offset[0], offset[1])
			if err == nil {
				if game.isSpaceEmpty(target) || game.isSpacePlayersPiece(target, opponentColour) {
					possibleMoves = append(
						possibleMoves,
						move{
							fromPos:   position,
							toPos:     target,
							fromPiece: game.getPieceSafe(position),
							toPiece:   game.getPieceSafe(target),
							isThreat:  game.isSpacePlayersPiece(target, opponentColour),
						},
					)
				}
			}
		}
	case Bishop:
		direction_offsets := [][2]int{
			{1, 1},
			{1, -1},
			{-1, 1},
			{-1, -1},
		}

		for _, direction_offset := range direction_offsets {
			current_pos := position
			for {
				target, err := positionRelative(current_pos, direction_offset[0], direction_offset[1])
				if err != nil {
					break
				}

				if game.getPieceSafe(target).Colour == fromPiece.Colour {
					break
				}

				newMove := move{
					fromPos:   position,
					toPos:     target,
					fromPiece: game.getPieceSafe(position),
					toPiece:   game.getPieceSafe(target),
					isThreat:  game.isSpacePlayersPiece(target, opponentColour),
				}

				possibleMoves = append(possibleMoves, newMove)

				if newMove.isThreat {
					break
				}

				current_pos = target
			}
		}
	case Rook:
		direction_offsets := [][2]int{
			{1, 0},
			{-1, 0},
			{0, 1},
			{0, -1},
		}

		for _, direction_offset := range direction_offsets {
			current_pos := position
			for {
				target, err := positionRelative(current_pos, direction_offset[0], direction_offset[1])
				if err != nil {
					break
				}

				if game.getPieceSafe(target).Colour == fromPiece.Colour {
					break
				}

				newMove := move{
					fromPos:   position,
					toPos:     target,
					fromPiece: game.getPieceSafe(position),
					toPiece:   game.getPieceSafe(target),
					isThreat:  game.isSpacePlayersPiece(target, opponentColour),
				}

				possibleMoves = append(possibleMoves, newMove)

				if newMove.isThreat {
					break
				}

				current_pos = target
			}
		}
	case Queen:
		direction_offsets := [][2]int{
			{1, 1},
			{1, 0},
			{1, -1},
			{0, 1},
			{0, -1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
		}

		for _, direction_offset := range direction_offsets {
			current_pos := position
			for {
				target, err := positionRelative(current_pos, direction_offset[0], direction_offset[1])
				if err != nil {
					break
				}

				if game.getPieceSafe(target).Colour == fromPiece.Colour {
					break
				}

				newMove := move{
					fromPos:   position,
					toPos:     target,
					fromPiece: game.getPieceSafe(position),
					toPiece:   game.getPieceSafe(target),
					isThreat:  game.isSpacePlayersPiece(target, opponentColour),
				}

				possibleMoves = append(possibleMoves, newMove)

				if newMove.isThreat {
					break
				}

				current_pos = target
			}
		}
	case King:
		target_offsets := [][2]int{
			{1, 1},
			{1, 0},
			{1, -1},
			{0, 1},
			{0, -1},
			{-1, 1},
			{-1, 0},
			{-1, -1},
		}

		for _, offset := range target_offsets {
			target, err := positionRelative(position, offset[0], offset[1])
			if err == nil {
				if game.isSpaceEmpty(target) || game.isSpacePlayersPiece(target, opponentColour) {
					possibleMoves = append(
						possibleMoves,
						move{
							fromPos:   position,
							toPos:     target,
							fromPiece: game.getPieceSafe(position),
							toPiece:   game.getPieceSafe(target),
							isThreat:  game.isSpacePlayersPiece(target, opponentColour),
						},
					)
				}
			}
		}
	}

	return possibleMoves, nil
}

func GetValidMovesForPlayer(game *GameState, player PlayerColour) []move {
	playerMoves := []move{}

	for ri, row := range game.BoardState {
		for ci := range row {
			pos, err := indexToPosition([2]int{ri, ci})
			if err != nil {
				panic(fmt.Errorf("something went wrong iterating board positions for valid moves - err '%v'", err))
			}
			if game.getPieceSafe(pos).Colour == player {
				pieceMoves, err := GetValidMovesForPiece(game, pos)
				if err != nil {
					panic(fmt.Errorf("something went wrong getting moves for piece at position '%v', - err: '%v'", pos, err))
				}
				playerMoves = slices.Concat(playerMoves, pieceMoves)
			}
		}
	}

	return playerMoves
}
