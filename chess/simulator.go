package chess

import (
	"errors"
	"strconv"
)

func prepScenario(moves []Move, b *Board) error {
	for _, move := range moves {
		found, piece := b.GetPieceAtSquare(move.From.Column, move.From.Row)
		if !found {
			return errors.New("No piece found at " + move.From.Column + strconv.Itoa(move.From.Row))
		}
		_, err := piece.Move(move.To.Column, move.To.Row, b, false)
		if err != nil {
			return err
		}
	}
	return nil
}
