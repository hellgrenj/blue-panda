package chess

import (
	"errors"
	"fmt"
)

func moveRook(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {
	if p.Type != rook {
		return nil, fmt.Errorf("MoveRook called on piece of type %v", p.Type)
	}
	if !p.InPlay {
		return nil, errors.New("Piece is not in play")
	}
	if !dryRun && !p.MoveIsLegal(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("move is not legal")
	}
	if !p.moveIsStraight(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("rooks can only move straight")
	}
	if _, err := p.moveJumpsOverPieces(targetColumn, targetRow, b); err != nil {
		return nil, err
	}
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
func (p *Piece) getValidRookMoves(b *Board) map[Move]*MoveResult {

	validMovesWithResult := make(map[Move]*MoveResult)
	possibleTargetSquares := []Square{}
	// horizontal and vertical * 7 squares (those out of the board will be filtered out by addToListIfValidSquare)
	// Up
	for i := 1; i < 8; i++ {
		row := p.CurrentSquare.Row + i
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, b.getColumnIndex(p.CurrentSquare.Column))
	}
	// Down
	for i := 1; i < 8; i++ {
		row := p.CurrentSquare.Row - i
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, b.getColumnIndex(p.CurrentSquare.Column))
	}
	// Right
	for i := 1; i < 8; i++ {
		columnIndex := b.getColumnIndex(p.CurrentSquare.Column) + i
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, p.CurrentSquare.Row, columnIndex)
	}
	// Left
	for i := 1; i < 8; i++ {
		columnIndex := b.getColumnIndex(p.CurrentSquare.Column) - i
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, p.CurrentSquare.Row, columnIndex)
	}

	for _, s := range possibleTargetSquares {
		result, err := moveRook(s.Column, s.Row, b, p, true)
		if err == nil {
			validMovesWithResult[Move{From: p.CurrentSquare, To: s}] = result
		}
	}
	return validMovesWithResult
}
