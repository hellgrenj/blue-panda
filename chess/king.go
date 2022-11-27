package chess

import "fmt"

func moveKing(targetColumn string, targetRow int, b *Board, p *Piece, dryRun bool) (*MoveResult, error) {
	if p.Type != king {
		return nil, fmt.Errorf("MoveKing called on piece of type %v", p.Type)
	}
	if !dryRun && !p.MoveIsLegal(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("move is not legal")
	}

	currentColumnValue := b.getColumnValue(p.CurrentSquare.Column)
	targetColumnValue := b.getColumnValue(targetColumn)

	if targetColumnValue > currentColumnValue+1 ||
		targetColumnValue < currentColumnValue-1 ||
		targetRow > p.CurrentSquare.Row+1 ||
		targetRow < p.CurrentSquare.Row-1 {
		return nil, fmt.Errorf("king cant move to square %v%v, can only move one square in any direction", targetColumn, targetRow)
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

	for _, s := range possibleTargetSquares {
		result, err := moveKing(s.Column, s.Row, b, p, true)
		if err == nil {
			validMovesWithResult[Move{From: p.CurrentSquare, To: s}] = result
		}
	}
	return validMovesWithResult
}
