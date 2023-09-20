package chess

import (
	"fmt"
)

func moveKing(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {
	if p.Type != king {
		return nil, fmt.Errorf("MoveKing called on piece of type %v", p.Type)
	}
	if !dryRun && !p.MoveIsLegal(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("move is not legal")
	}

	currentColumnValue := b.getColumnValue(p.CurrentSquare.Column)
	targetColumnValue := b.getColumnValue(targetColumn)

	// check if castling attempt for either black or white
	if targetRow == p.CurrentSquare.Row && targetColumnValue == currentColumnValue+2 ||
		targetRow == p.CurrentSquare.Row && targetColumnValue == currentColumnValue-2 {
		if err := p.tryCastling(targetColumn, targetRow, b, dryRun); err == nil {
			return &MoveResult{Action: Castling, Piece: nil}, nil
		} else {
			return nil, fmt.Errorf("castling not allowed, err: %v", err)
		}
	}

	if targetColumnValue > currentColumnValue+1 ||
		targetColumnValue < currentColumnValue-1 ||
		targetRow > p.CurrentSquare.Row+1 ||
		targetRow < p.CurrentSquare.Row-1 {
		return nil, fmt.Errorf("king cant move to square %v%v, can only move one square in any direction, except when castling", targetColumn, targetRow)
	}

	enemyAtTargetSquare, enemyPiece := b.targetSquareOccupiedByEnemy(targetColumn, targetRow, p)
	if enemyAtTargetSquare {
		if !dryRun {
			p.takeAt(targetColumn, targetRow, enemyPiece, b)
		}
		return &MoveResult{Action: Take, Piece: enemyPiece}, nil

	} else {
		occupied, pieceAtTarget := b.GetPieceAtSquare(targetColumn, targetRow)
		if occupied {
			return nil, fmt.Errorf("%v %v cant move to square %v%v, it is occupied by %v %v", p.Colour, p.Type, targetColumn, targetRow, pieceAtTarget.Colour, pieceAtTarget.Type)
		} else { // not occupied
			if !dryRun {
				p.goTo(targetColumn, targetRow, b)
			}
			return &MoveResult{Action: GoTo, Piece: nil}, nil
		}
	}

}

type castleSide int64

const (
	kingside castleSide = iota
	queenside
	noCastle
)

func (king *Piece) tryCastling(targetColumn string, targetRow int, b *Board, dryRun bool) error {
	chosenRook, castleSide, err := selectRookAndSideForCastling(king, b, targetColumn)
	if err != nil {
		return err
	}

	// check if king or rook has previously moved
	if chosenRook.hasMoved || king.hasMoved {
		return fmt.Errorf("king or rook has previously moved")
	}
	// the king can not jump over pieces nor castle through check
	// (into check is covered by MoveIsLegal in main func MoveKing)
	if castleSide == kingside {
		_, err := b.checkPathForOccupiedSquaresStraigthRight(chosenRook.CurrentSquare.Column, chosenRook.CurrentSquare.Row, king)
		if err != nil {
			return fmt.Errorf("pieces between king and rook")
		}
		err = b.checkPathForSquaresUnderAttackStraightRight(chosenRook.CurrentSquare.Column, chosenRook.CurrentSquare.Row, king)
		if err != nil {
			return fmt.Errorf("king passes through a square that is attacked by an enemy piece")
		}

	} else { // check queenside castling
		_, err := b.checkPathForOccupiedSquaresStraightLeft(chosenRook.CurrentSquare.Column, chosenRook.CurrentSquare.Row, king)
		if err != nil {
			return fmt.Errorf("pieces between king and rook")
		}
		err = b.checkPathForSquaresUnderAttackStraightLeft(chosenRook.CurrentSquare.Column, chosenRook.CurrentSquare.Row, king)
		if err != nil {
			return fmt.Errorf("king passes through a square that is attacked by an enemy piece")
		}
	}
	// the king is not currently in check
	if isCheck, _ := b.kingIsInCheck(king.Colour); isCheck {
		return fmt.Errorf("king is in check")
	}

	// move king two squares towards rook (kingside or queenside)
	if !dryRun {
		king.goTo(targetColumn, targetRow, b)
	}
	// move rook to other side of king, directly next to it on the opposite side
	if castleSide == kingside { // castle kingside
		if !dryRun {
			// move rook to the left of the king
			chosenRook.goTo(b.getColumnStringByIndex(b.getColumnIndex(king.CurrentSquare.Column)-1), king.CurrentSquare.Row, b)
		}
	} else { // castle queenside
		if !dryRun {
			// move rook to the right of the king
			chosenRook.goTo(b.getColumnStringByIndex(b.getColumnIndex(king.CurrentSquare.Column)+1), king.CurrentSquare.Row, b)
		}
	}
	return nil
}

func selectRookAndSideForCastling(k *Piece, b *Board, targetColumn string) (*Piece, castleSide, error) {
	if k.Type != king {
		return nil, noCastle, fmt.Errorf("Piece is not a King, cant select a rook and side for castling")
	}
	var kingSideRook *Piece
	var kingSideRookFound bool
	var queenSideRook *Piece
	var queenSideRookFound bool

	if k.Colour == White {
		kingSideRookFound, kingSideRook = b.GetPieceAtSquare("H", 1)
		queenSideRookFound, queenSideRook = b.GetPieceAtSquare("A", 1)
	} else {
		kingSideRookFound, kingSideRook = b.GetPieceAtSquare("H", 8)
		queenSideRookFound, queenSideRook = b.GetPieceAtSquare("A", 8)
	}

	targetColumnValue := b.getColumnValue(targetColumn)
	kingSideRookColumn := 8  // regardless of color and if it exists
	queenSideRookColumn := 1 // regardless of color and if it exists
	distanceToKingsideRook := kingSideRookColumn - targetColumnValue
	distanceToQueensideRook := queenSideRookColumn - targetColumnValue
	var chosenRook *Piece

	// math absolute value
	if distanceToKingsideRook < 0 {
		distanceToKingsideRook = distanceToKingsideRook * -1
	}
	if distanceToQueensideRook < 0 {
		distanceToQueensideRook = distanceToQueensideRook * -1
	}

	var castleSide castleSide
	if distanceToKingsideRook < distanceToQueensideRook {
		castleSide = kingside
		if !kingSideRookFound {
			return nil, noCastle, fmt.Errorf("no rook found on kingside")
		}
		chosenRook = kingSideRook
	} else {
		castleSide = queenside
		if !queenSideRookFound {
			return nil, noCastle, fmt.Errorf("no rook found on queenside")
		}
		chosenRook = queenSideRook
	}
	return chosenRook, castleSide, nil
}

func (p *Piece) kingTryRun(b *Board) error {
	if p.Type != king {
		return fmt.Errorf("Piece is not a King")
	}
	if isCheck, _ := b.kingIsInCheck(p.Colour); !isCheck {
		return fmt.Errorf("king is not in check, no need to run")
	}
	validMoves := p.getValidKingMoves(b)
	if len(validMoves) == 0 {
		return fmt.Errorf("king cant run, no valid moves")
	}
	var possibleSquares []Square
	for k := range validMoves {
		possibleSquares = append(possibleSquares, k.To)
	}
	realCurrentSquare := p.CurrentSquare
	for _, square := range possibleSquares {
		_, err := moveKing(square.Column, square.Row, b, p, true)
		if err == nil {

			// check if king in check on pending move
			currentSquare := p.CurrentSquare
			// if enemy on target square simulate take in temp move
			targetSquareOccupied, pieceAtTargetSquare := b.GetPieceAtSquare(square.Column, square.Row)
			if targetSquareOccupied {
				if p.enemyTo(pieceAtTargetSquare) {
					pieceAtTargetSquare.InPlay = false // temp take
				}
			}
			// else
			p.CurrentSquare = square // temp move

			if isCheck, _ := b.kingIsInCheck(p.Colour); !isCheck {
				// fmt.Printf("King can run to %v%v\n", square.Column, square.Row)
				if targetSquareOccupied {
					pieceAtTargetSquare.InPlay = true // reset temp take
				}
				p.CurrentSquare = realCurrentSquare // reset temp move
				return nil
			} else {
				fmt.Printf("King could move to  %v%v, but it is still in check\n", square.Column, square.Row)
			}
			if targetSquareOccupied {
				pieceAtTargetSquare.InPlay = true // reset temp take
			}
			p.CurrentSquare = currentSquare // reset temp move

		} else {
			fmt.Printf("King cant run to %v%v\n", square.Column, square.Row)
		}
		p.CurrentSquare = realCurrentSquare // undo temp move
	}

	return fmt.Errorf("king cant outrun check")
}
func (p *Piece) getValidKingMoves(b *Board) map[Move]*MoveResult {

	validMovesWithResult := make(map[Move]*MoveResult)
	possibleTargetSquares := []Square{}

	// Up right
	columnIndex := b.getColumnIndex(p.CurrentSquare.Column) + 1
	row := p.CurrentSquare.Row + 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// Up left
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
	row = p.CurrentSquare.Row + 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// Down right
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
	row = p.CurrentSquare.Row - 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// Down left
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
	row = p.CurrentSquare.Row - 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// up
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column)
	row = p.CurrentSquare.Row + 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// down
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column)
	row = p.CurrentSquare.Row - 1
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// right
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 1
	row = p.CurrentSquare.Row
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)
	// left
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 1
	row = p.CurrentSquare.Row
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

	// castling queenside
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) - 2
	row = p.CurrentSquare.Row
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

	// castling kingside
	columnIndex = b.getColumnIndex(p.CurrentSquare.Column) + 2
	row = p.CurrentSquare.Row
	possibleTargetSquares = addToListIfValidSquare(p, b, possibleTargetSquares, row, columnIndex)

	for _, s := range possibleTargetSquares {
		result, err := moveKing(s.Column, s.Row, b, p, true)
		if err == nil {
			validMovesWithResult[Move{From: p.CurrentSquare, To: s}] = result
		}
	}
	return validMovesWithResult
}
