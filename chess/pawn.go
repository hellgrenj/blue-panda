package chess

import (
	"errors"
	"fmt"
)

func movePawn(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {

	if p.Type != pawn {
		return nil, fmt.Errorf("MovePawn called on piece of type %v", p.Type)
	}

	if !p.InPlay {
		return nil, errors.New("Piece is not in play")
	}
	if !dryRun && !p.MoveIsLegal(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("illegal move")
	}

	if p.moveIsStraight(targetColumn, targetRow, b) {
		if b.getColumnIndex(p.CurrentSquare.Column) != b.getColumnIndex(targetColumn) &&
			p.CurrentSquare.Row == targetRow {
			return nil, fmt.Errorf("pawns cant move sideways")
		}
		if p.Colour == White && targetRow < p.CurrentSquare.Row {
			return nil, fmt.Errorf("pawns can only move forward")
		}
		if p.Colour == Black && targetRow > p.CurrentSquare.Row {
			return nil, fmt.Errorf("pawns can only move forward")
		}
		if p.hasMoved && p.Colour == White && targetRow > (p.CurrentSquare.Row+1) {
			err := fmt.Errorf("pawn cant move to square %v%v, can only move two squares forward on first move and then one square", targetColumn, targetRow)
			return nil, err
		}
		if !p.hasMoved && p.Colour == White && targetRow > (p.CurrentSquare.Row+2) {
			err := fmt.Errorf("pawn cant move to square %v%v, can only move two squares forward on first move and then one square", targetColumn, targetRow)
			return nil, err
		}
		if p.hasMoved && p.Colour == Black && targetRow < (p.CurrentSquare.Row-1) {
			err := fmt.Errorf("pawn cant move to square %v%v, can only move two squares forward on first move and then one square", targetColumn, targetRow)
			return nil, err
		}
		if !p.hasMoved && p.Colour == Black && targetRow < (p.CurrentSquare.Row-2) {
			err := fmt.Errorf("pawn cant move to square %v%v, can only move two squares forward on first move and then one square", targetColumn, targetRow)
			return nil, err
		}

		if _, err := p.moveJumpsOverPieces(targetColumn, targetRow, b); err != nil {
			return nil, err
		}

		occupied, byPiece := b.GetPieceAtSquare(targetColumn, targetRow)
		if occupied {
			err := fmt.Errorf("square %v %v is occupied by %v", targetColumn, targetRow, byPiece)
			return nil, err
		} else {
			if !dryRun {
				p.goTo(targetColumn, targetRow, b)
				p.tryPromoteToQueen()
			}
			return &MoveResult{Action: GoTo, Piece: nil}, nil
		}
	} else if p.moveIsDiagonal(targetColumn, targetRow, b) {

		if targetRow > (p.CurrentSquare.Row+1) || targetRow < (p.CurrentSquare.Row-1) {
			err := fmt.Errorf("%v cant move to square %v%v, can only move one square diagonally", p.Type, targetColumn, targetRow)
			return nil, err
		}
		if p.Colour == White && targetRow < p.CurrentSquare.Row {
			return nil, fmt.Errorf("pawns can only move (and take) forward")
		}
		if p.Colour == Black && targetRow > p.CurrentSquare.Row {
			return nil, fmt.Errorf("pawns can only move (and take) forward")
		}
		currentColumnValue := b.getColumnValue(p.CurrentSquare.Column)
		targetColumnValue := b.getColumnValue(targetColumn)

		if targetColumnValue == currentColumnValue+1 || targetColumnValue == currentColumnValue-1 {
			enemyAtTargetSquare, enemyPiece := b.targetSquareOccupiedByEnemy(targetColumn, targetRow, p)
			if enemyAtTargetSquare {
				if !dryRun {
					p.takeAt(targetColumn, targetRow, enemyPiece, b)
					p.tryPromoteToQueen()
				}
				return &MoveResult{Action: Take, Piece: enemyPiece}, nil

			} else { // diagonal move, but no enemy piece at target square
				err := fmt.Errorf("piece %v cant move to %v %v, can only move diagonally when taking", p, targetColumn, targetRow)
				return nil, err
			}
		} else {
			err := fmt.Errorf("pawn can only move 1 step diagonally")
			return nil, err
		}

	}
	return nil, fmt.Errorf("unknown error")
}
func (p *Piece) tryPromoteToQueen() error {
	if p.Type != pawn {
		return fmt.Errorf("can only promote pawns")
	}
	if p.Colour == White && p.CurrentSquare.Row == 8 {
		p.Type = queen
		return nil
	}
	if p.Colour == Black && p.CurrentSquare.Row == 1 {
		p.Type = queen
		return nil
	}
	return fmt.Errorf("pawn cant be promoted")
}
func (p *Piece) getValidPawnMoves(b *Board) map[Move]*MoveResult {

	validMovesWithResult := make(map[Move]*MoveResult)

	possibleTargetSquares := []Square{}
	if p.Colour == White {
		if !p.hasMoved {
			// Up two squares!
			columnIndex := b.getColumnIndex(p.CurrentSquare.Column)
			row := p.CurrentSquare.Row + 2
			possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
		}
		// Up one square
		columnIndex := b.getColumnIndex(p.CurrentSquare.Column)
		row := p.CurrentSquare.Row + 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

		// Up right one square
		columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
		row = p.CurrentSquare.Row + 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

		// Up left one square
		columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
		row = p.CurrentSquare.Row + 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	} else {
		if !p.hasMoved {
			// Down two squares!
			columnIndex := b.getColumnIndex(p.CurrentSquare.Column)
			row := p.CurrentSquare.Row - 2
			possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
		}
		// Down one square
		columnIndex := b.getColumnIndex(p.CurrentSquare.Column)
		row := p.CurrentSquare.Row - 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

		// Down right one square
		columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
		row = p.CurrentSquare.Row - 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

		// Down left one square
		columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
		row = p.CurrentSquare.Row - 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	}

	for _, s := range possibleTargetSquares {
		result, err := movePawn(s.Column, s.Row, b, p, true)
		if err == nil {
			validMovesWithResult[Move{From: p.CurrentSquare, To: s}] = result
		}
	}
	return validMovesWithResult
}
