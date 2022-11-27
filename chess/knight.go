package chess

import (
	"errors"
	"fmt"
)

func moveKnight(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {
	if p.Type != knight {
		return nil, fmt.Errorf("MoveKnight called on piece of type %v", p.Type)
	}
	if !p.InPlay {
		return nil, errors.New("Piece is not in play")
	}
	if !dryRun && !p.MoveIsLegal(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("move is not legal")
	}
	if squareIsPossibleTarget(p, targetRow, targetColumn, b) {
		occupied, pieceAtTarget := b.GetPieceAtSquare(targetColumn, targetRow)
		if occupied && !p.enemyTo(pieceAtTarget) {
			return nil, fmt.Errorf("%v %v cant move to square %v%v, it is occupied by %v %v", p.Colour, p.Type, targetColumn, targetRow, pieceAtTarget.Colour, pieceAtTarget.Type)
		} else if occupied && p.enemyTo(pieceAtTarget) {
			if !dryRun {
				p.takeAt(targetColumn, targetRow, pieceAtTarget, b)
			}
			return &MoveResult{Action: Take, Piece: pieceAtTarget}, nil
		} else { // not occupied
			if !dryRun {
				p.goTo(targetColumn, targetRow, b)
			}
			return &MoveResult{Action: GoTo, Piece: nil}, nil
		}
	}
	return nil, fmt.Errorf("move is not legal")
}

func squareIsPossibleTarget(p *Piece, targetRow int, targetColumn string, b *Board) bool {
	possibleTargetSquares := getPossibleSquaresForKnight(p, b)
	for _, square := range possibleTargetSquares {
		if square.Row == targetRow && square.Column == targetColumn {
			return true
		}
	}
	return false
}
func (p *Piece) getValidKnightMoves(b *Board) map[Move]*MoveResult {

	validMovesWithResult := make(map[Move]*MoveResult)
	possibleTargetSquares := getPossibleSquaresForKnight(p, b)

	for _, s := range possibleTargetSquares {
		result, err := moveKnight(s.Column, s.Row, b, p, true)
		if err == nil {
			validMovesWithResult[Move{From: p.CurrentSquare, To: s}] = result
		}
	}
	return validMovesWithResult
}
func getPossibleSquaresForKnight(p *Piece, b *Board) []Square {
	possibleTargetSquares := []Square{}
	// two right and one up
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  TT  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	columnIndex := b.getColumnIndex(p.CurrentSquare.Column) + 2
	row := p.CurrentSquare.Row + 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

	// two right and one down
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  TT  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 2
	row = p.CurrentSquare.Row - 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// two left and one up
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  TT  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 2
	row = p.CurrentSquare.Row + 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// two left and one down
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  TT  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 2
	row = p.CurrentSquare.Row - 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// two up and one right
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  TT  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
	row = p.CurrentSquare.Row + 2
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// two up and one left
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  TT  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
	row = p.CurrentSquare.Row + 2
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// two down and one right
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  TT  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
	row = p.CurrentSquare.Row - 2
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// two down and one left
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  WN  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  TT  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..

	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
	row = p.CurrentSquare.Row - 2
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

	return possibleTargetSquares
}
