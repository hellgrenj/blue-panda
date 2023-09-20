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
		return movePawnStraight(targetColumn, targetRow, b, p, dryRun)
	} else if p.moveIsDiagonal(targetColumn, targetRow, b) {
		return movePawnDiagonally(targetColumn, targetRow, b, p, dryRun)
	}
	return nil, fmt.Errorf("unknown error")
}

func movePawnStraight(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {
	if p.Type != pawn {
		return nil, fmt.Errorf("movePawnStraight called on piece of type %v", p.Type)
	}
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
}

func movePawnDiagonally(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {
	if p.Type != pawn {
		return nil, fmt.Errorf("movePawnDiagonally called on piece of type %v", p.Type)
	}
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

			if enemyTaken, err := p.tryEnPassant(b, dryRun); err == nil {
				return &MoveResult{Action: EnPassant, Piece: enemyTaken}, nil
			} else {

				finalErr := fmt.Errorf("piece %v cant move to %v %v, can only move diagonally when taking (en passant not possible, reason: %v)", p, targetColumn, targetRow, err)
				return nil, finalErr
			}
		}
	} else {
		err := fmt.Errorf("pawn can only move 1 step diagonally")
		return nil, err
	}
}
func (p *Piece) tryEnPassant(b *Board, dryRun bool) (*Piece, error) {
	if p.Colour == White {
		return tryEnPassantWhite(p, b, dryRun)
	} else {
		return tryEnPassantBlack(p, b, dryRun)
	}
}

func tryEnPassantWhite(p *Piece, b *Board, dryRun bool) (*Piece, error) {
	if p.Type != pawn {
		return nil, fmt.Errorf("en passant is only for pawns")
	}

	if b.blacksLastMove.Piece == nil && b.blacksLastMove.Move == nil {
		return nil, fmt.Errorf("black has not moved")
	}
	// if oppents last move was pawn two squares forward..
	if b.blacksLastMove.Piece.Type == pawn && b.blacksLastMove.Move.From.Row == 7 && b.blacksLastMove.Move.To.Row == 5 {
		// ...and we are next to the square it moved to..
		columnIndexDiff := b.getColumnIndex(p.CurrentSquare.Column) - b.getColumnIndex(b.blacksLastMove.Move.To.Column)
		if columnIndexDiff < 0 {
			columnIndexDiff = columnIndexDiff * -1
		}
		if columnIndexDiff == 1 && p.CurrentSquare.Row == 5 && b.blacksLastMove.Move.To.Row == 5 {
			if !dryRun {
				p.goTo(b.blacksLastMove.Move.From.Column, b.blacksLastMove.Move.From.Row-1, b) // move to square behind enemy pawn
				enemyFound, enemy := b.GetPieceAtSquare(b.blacksLastMove.Move.To.Column, b.blacksLastMove.Move.To.Row)
				if !enemyFound {
					return nil, fmt.Errorf("error: enemy not found, expected enemy pawn at %v%v", b.whitesLastMove.Move.To.Column, b.whitesLastMove.Move.To.Row)
				}
				enemy.InPlay = false

			}
			return b.blacksLastMove.Piece, nil
		}
		return &Piece{}, fmt.Errorf("pawn not next to oppents pawn on row 5")
	} else {
		return &Piece{}, fmt.Errorf("oppents last move was not pawn two squares forward")
	}
}
func tryEnPassantBlack(p *Piece, b *Board, dryRun bool) (*Piece, error) {
	if p.Type != pawn {
		return nil, fmt.Errorf("en passant is only for pawns")
	}

	if b.whitesLastMove.Piece == nil && b.whitesLastMove.Move == nil {
		return nil, fmt.Errorf("white has not moved")
	}
	// if oppents last move was pawn two squares forward..
	if b.whitesLastMove.Piece.Type == pawn && b.whitesLastMove.Move.From.Row == 2 && b.whitesLastMove.Move.To.Row == 4 {
		// ...and we are next to the square it moved to..
		columnIndexDiff := b.getColumnIndex(p.CurrentSquare.Column) - b.getColumnIndex(b.whitesLastMove.Move.To.Column)
		if columnIndexDiff < 0 {
			columnIndexDiff = columnIndexDiff * -1
		}
		if columnIndexDiff == 1 && p.CurrentSquare.Row == 4 && b.whitesLastMove.Move.To.Row == 4 {
			if !dryRun {
				p.goTo(b.whitesLastMove.Move.From.Column, b.whitesLastMove.Move.From.Row+1, b) // move to square behind enemy pawn
				enemyFound, enemy := b.GetPieceAtSquare(b.whitesLastMove.Move.To.Column, b.whitesLastMove.Move.To.Row)
				if !enemyFound {
					return nil, fmt.Errorf("error: enemy not found, expected enemy pawn at %v%v", b.whitesLastMove.Move.To.Column, b.whitesLastMove.Move.To.Row)
				}
				enemy.InPlay = false

			}
			return b.whitesLastMove.Piece, nil
		}
		return &Piece{}, fmt.Errorf("pawn not next to oppents pawn on row 4")
	} else {
		return &Piece{}, fmt.Errorf("oppents last move was not pawn two squares forward")
	}
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

		// Up right one square (covers en passant as well)
		columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
		row = p.CurrentSquare.Row + 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

		// Up left one square (covers en passant as well)
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

		// Down right one square (covers en passant as well)
		columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
		row = p.CurrentSquare.Row - 1
		possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

		// Down left one square (covers en passant as well)
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
