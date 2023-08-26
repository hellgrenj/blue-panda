package chess

import (
	"fmt"
	"strings"
)

type Colour int64

const (
	White Colour = iota
	Black
)

func (c Colour) String() string {
	if c == White {
		return "White"
	} else {
		return "Black"
	}
}

type ptype int64

const (
	pawn ptype = iota
	rook
	knight
	bishop
	queen
	king
)

func (t ptype) String() string {
	switch t {
	case pawn:
		return "Pawn"
	case rook:
		return "Rook"
	case knight:
		return "Knight"
	case bishop:
		return "Bishop"
	case queen:
		return "Queen"
	case king:
		return "King"
	default:
		return "Unknown"
	}
}

type Piece struct {
	Type          ptype
	CurrentSquare Square
	Colour        Colour
	InPlay        bool
	hasMoved      bool // pawn 2 squares first move, king castling if king and rook hasn't moved etc...
}

// returns the abbreviation for the piece, this can be used by (e.g) a BoardVisualizer  to visualize the board
func (p *Piece) GetAbbreveation() string {
	switch p.Type {
	case pawn:
		if p.Colour == White {
			return "♙"
		} else {
			return "♟"
		}
	case rook:
		if p.Colour == White {
			return "♖"
		} else {
			return "♜"
		}
	case knight:
		if p.Colour == White {
			return "♘"
		} else {
			return "♞"
		}
	case bishop:
		if p.Colour == White {
			return "♗"
		} else {
			return "♝"
		}
	case queen:
		if p.Colour == White {
			return "♕"
		} else {
			return "♛"
		}
	case king:
		if p.Colour == White {
			return "♔"
		} else {
			return "♚"
		}
	default:
		return "??"
	}
}
func (p *Piece) moveIsStraight(targetColumn string, targetRow int, b *Board) bool {
	if targetColumn == p.CurrentSquare.Column && targetRow != p.CurrentSquare.Row {
		return true
	} else if targetColumn != p.CurrentSquare.Column && targetRow == p.CurrentSquare.Row {
		return true
	} else {
		return false
	}
}
func (p *Piece) moveIsDiagonal(targetColumn string, targetRow int, b *Board) bool {
	if targetColumn != p.CurrentSquare.Column && targetRow != p.CurrentSquare.Row {
		// check if diagonal all the way
		columnDiff := b.getColumnValue(targetColumn) - b.getColumnValue(p.CurrentSquare.Column)
		rowDiff := targetRow - p.CurrentSquare.Row
		// columnDiff absolut value
		if columnDiff < 0 {
			columnDiff = columnDiff * -1
		}
		// rowDiff absolut value
		if rowDiff < 0 {
			rowDiff = rowDiff * -1
		}
		if columnDiff == rowDiff {
			return true
		}
		return false
	} else {
		return false
	}
}
func (p *Piece) moveIsNone(targetColumn string, targetRow int, b *Board) bool {
	if targetColumn == p.CurrentSquare.Column && targetRow == p.CurrentSquare.Row {
		return true
	} else {
		return false
	}
}
func (p *Piece) takeAt(targetColumn string, targetRow int, enemy *Piece, b *Board) {
	fmt.Printf("%v %v is taking %v %v", p.Colour, p.Type, enemy.Colour, enemy.Type)
	square, err := b.getSquare(targetColumn, targetRow)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	p.CurrentSquare = square
	p.hasMoved = true // if first move is a take.. (else it is set in GoTo function)
	enemy.InPlay = false
}
func (p *Piece) couldMoveTo(targetColumn string, targetRow int, b *Board) bool {
	dryRun := true
	if _, err := p.Move(targetColumn, targetRow, b, dryRun); err == nil {
		return true
	} else {
		return false
	}
}
func (p *Piece) goTo(targetColumn string, targetRow int, b *Board) {
	fmt.Printf("piece %v is moving to %v%v\n", p, targetColumn, targetRow)
	square, err := b.getSquare(targetColumn, targetRow)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	p.CurrentSquare = square
	p.hasMoved = true // is also set in TakeAt function
}

// Checks if the target square is valid and that the move doesnt put the own king in check
func (p *Piece) MoveIsLegal(targetColumn string, targetRow int, b *Board) bool {
	_, err := b.getSquare(targetColumn, targetRow)
	if err != nil { // check if target square is on board
		fmt.Printf("error: %v", err)
		return false
	}

	// check if king in check on pending move
	currentSquare := p.CurrentSquare
	s, _ := b.getSquare(targetColumn, targetRow)

	// if enemy on target square simulate take in temp move
	targetSquareOccupied, pieceAtTargetSquare := b.GetPieceAtSquare(s.Column, s.Row)
	if targetSquareOccupied {
		if p.enemyTo(pieceAtTargetSquare) {
			pieceAtTargetSquare.InPlay = false // temp take
		}
	}
	p.CurrentSquare = s // temp move

	if isCheck, _ := b.kingIsInCheck(p.Colour); isCheck {

		if targetSquareOccupied {
			pieceAtTargetSquare.InPlay = true // reset temp take
		}
		p.CurrentSquare = currentSquare // reset temp move
		return false
	}

	if targetSquareOccupied {
		pieceAtTargetSquare.InPlay = true // reset temp take
	}
	p.CurrentSquare = currentSquare // reset temp move

	return true

}

func (p *Piece) enemyTo(piece *Piece) bool {
	return p.Colour != piece.Colour
}

// returns the value of a piece
func (p *Piece) GetValue() int {
	switch p.Type {
	case pawn:
		return 1
	case rook:
		return 5
	case knight:
		return 3
	case bishop:
		return 3
	case queen:
		return 9
	case king:
		return 100
	default:
		return 0
	}
}

func (p *Piece) moveJumpsOverPieces(targetColumn string, targetRow int, b *Board) ([]Square, error) {
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	targetColumnIndex := b.getColumnIndex(targetColumn)
	if p.moveIsStraight(targetColumn, targetRow, b) {
		if currentColumnIndex > targetColumnIndex {
			return b.checkPathStraightLeft(targetColumn, targetRow, p)
		} else {
			if squaresInBetween, foundPieceOnPathError := b.checkPathStraigthRight(targetColumn, targetRow, p); foundPieceOnPathError != nil {
				return squaresInBetween, foundPieceOnPathError
			}
		}
		if p.CurrentSquare.Row > targetRow {
			return b.checkPathStraightDown(targetColumn, targetRow, p)
		} else {
			return b.checkPathStraightUp(targetColumn, targetRow, p)
		}
	} else if p.moveIsDiagonal(targetColumn, targetRow, b) {
		/*
			8
			7  lowerColumn higher Row					higherColumn higherRow
			6
			5
			4
			3   lowerColumn lower Row					higherColumn lowerRow
			2
			1   A  		B  		C  		D	  E  		F  		G 		H
		*/
		if targetColumnIndex > currentColumnIndex { // RIGHT
			if targetRow > p.CurrentSquare.Row { // UP
				return b.checkPathUpRight(targetColumn, targetRow, p)

			} else { // DOWN
				return b.checkPathDownRight(targetColumn, targetRow, p)
			}
		} else { // LEFT
			if targetRow > p.CurrentSquare.Row { // UP
				return b.checkPathUpLeft(targetColumn, targetRow, p)

			} else { //  DOWN
				return b.checkPathDownLeft(targetColumn, targetRow, p)
			}
		}
	}
	return []Square{}, nil
}

// moves the piece and returns a MoveResult or an error
func (p *Piece) Move(targetColumn string, targetRow int, b *Board, dryRun bool) (*MoveResult, error) {
	targetColumn = strings.ToUpper(targetColumn)
	if p.moveIsNone(targetColumn, targetRow, b) {
		return nil, fmt.Errorf("already there")
	}
	if !dryRun && b.kingIsInMate(p.Colour) {
		return nil, fmt.Errorf("%v king is in mate", p.Colour)
	}
	switch p.Type {
	case pawn:
		return movePawn(targetColumn, targetRow, b, p, dryRun)
	case rook:
		return moveRook(targetColumn, targetRow, b, p, dryRun)
	case knight:
		return moveKnight(targetColumn, targetRow, b, p, dryRun)
	case bishop:
		return moveBishop(targetColumn, targetRow, b, p, dryRun)
	case queen:
		return moveQueen(targetColumn, targetRow, b, p, dryRun)
	case king:
		return moveKing(targetColumn, targetRow, b, p, dryRun)
	default:
		return nil, fmt.Errorf("unknown piece type: %v", p.Type)
	}
}
