package chess

import (
	"errors"
	"fmt"
	"strings"
)

type Square struct {
	Column string
	Row    int
}

type LastMove struct {
	Piece *Piece
	Move  *Move
}
type Board struct {
	Squares        [64]Square
	WhitePieces    []Piece
	BlackPieces    []Piece
	whitesLastMove LastMove
	blacksLastMove LastMove
	columns        []string
}
type Move struct {
	From Square
	To   Square
}
type MoveResultAction int64

const (
	Take MoveResultAction = iota
	GoTo
)

type MoveResult struct {
	Action MoveResultAction
	Piece  *Piece // taken piece so we can calculate and compare value of moves
}

func newBoard() *Board {
	board := &Board{}
	board.init()
	return board
}
func (b *Board) init() {
	b.columns = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	rows := []int{8, 7, 6, 5, 4, 3, 2, 1}

	for i, row := range rows {
		for j, column := range b.columns {
			b.Squares[i*8+j] = Square{column, row}
		}
	}
	b.placePiecesOnBoard()
}

func (b *Board) getPosition() string {
	var position string
	for s := range b.Squares {
		occupied, piece := b.GetPieceAtSquare(b.Squares[s].Column, b.Squares[s].Row)
		if occupied {
			position += piece.GetAbbreveation()
		} else {
			position += "--"
		}
	}
	return position
}
func (b *Board) targetSquareOccupiedByEnemy(targetColumn string, targetRow int, p *Piece) (bool, *Piece) {
	occupied, piece := b.GetPieceAtSquare(targetColumn, targetRow)
	if occupied && piece.Colour != p.Colour {
		return true, piece
	} else {
		return false, piece
	}
}
func (b *Board) isStaleMate(colour Colour) bool {
	if isCheck, _ := b.kingIsInCheck(colour); isCheck {
		return false
	}
	if b.kingIsInMate(colour) {
		return false
	}
	moves := b.GetAllMovesFor(colour)
	if len(moves) == 0 {
		return true
	} else {
		for m := range moves { // check if any of the moves would NOT result in check
			_, piece := b.GetPieceAtSquare(m.From.Column, m.From.Row)
			if piece.MoveIsLegal(m.To.Column, m.To.Row, b) {
				return false
			}
		}
	}
	// if no valid moves are found, then it is a stale mate
	return true
}
func (b *Board) kingIsInMate(colour Colour) bool {
	king := b.getKing(colour)

	isCheck, enemies := b.kingIsInCheck(colour)
	if !isCheck {
		return false // king is not in check, so not in mate
	}

	cantEscapeErr := king.kingTryRun(b)
	if cantEscapeErr == nil {
		return false // king can run, so not in mate
	}

	if couldAnyPieceBlockCheck(king, enemies, colour, b) {
		return false // a friendly piece can block the check (one or multiple threathening pieces) = no check mate
	}

	if len(enemies) == 1 { // if only one piece is threathening the king, then it might be possible to take it

		if colour == White {
			if CouldAnyTakeAt(b.WhitePieces, enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row, b) {
				return false // a friendly piece can take the ONE treathening piece = no check mate
			}
		} else {
			if CouldAnyTakeAt(b.BlackPieces, enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row, b) {
				return false // a friendly piece can take the ONE treathening piece = no check mate
			}
		}
	}

	return true // king cant run, no piece can take the ONE threathening piece and no piece can block the check = check mate
}

func couldAnyPieceBlockCheck(king *Piece, enemies []Piece, colour Colour, b *Board) bool {
	var numberOfSuccesfulBlocks int
	for _, enemy := range enemies {
		squaresInBetween, err := b.getSquaresBetween(&enemy, king)
		if err != nil {
			continue // not in line of attack (the famous: this should not happen)
		}
		for _, square := range squaresInBetween {
			// check if any pieces can move to a square that would remove the check
			if colour == White {
				if couldAnyBlockCheckAt(b.WhitePieces, square.Column, square.Row, b) {
					numberOfSuccesfulBlocks++
					break // it is enough that one square of the line of attack can be blocked
				}
			} else {
				if couldAnyBlockCheckAt(b.BlackPieces, square.Column, square.Row, b) {
					numberOfSuccesfulBlocks++
					break // it is enough that one square of the line of attack can be blocked
				}
			}
		}
	}
	return numberOfSuccesfulBlocks == len(enemies)
}

func (b *Board) getSquaresBetween(attacker *Piece, target *Piece) ([]Square, error) {
	if !attacker.couldMoveTo(target.CurrentSquare.Column, target.CurrentSquare.Row, b) {
		return nil, fmt.Errorf("target square %v%v is not in line of attack", target.CurrentSquare.Column, target.CurrentSquare.Row)
	}
	squaresInBetween, err := attacker.moveJumpsOverPieces(target.CurrentSquare.Column, target.CurrentSquare.Row, b)
	if err != nil {
		return nil, err
	}
	return squaresInBetween, nil
}

func CouldAnyTakeAt(pieces []Piece, targetColumn string, targetRow int, b *Board) bool {
	for _, piece := range pieces {
		if piece.couldMoveTo(targetColumn, targetRow, b) {
			_, ref := b.GetPieceAtSquare(piece.CurrentSquare.Column, piece.CurrentSquare.Row)
			realCurrentSquare := ref.CurrentSquare
			targetSquareOccupied, pieceAtTargetSquare := b.GetPieceAtSquare(targetColumn, targetRow)
			if targetSquareOccupied {
				if ref.enemyTo(pieceAtTargetSquare) {
					pieceAtTargetSquare.InPlay = false // temp take
				}
			}
			ref.CurrentSquare = Square{Column: targetColumn, Row: targetRow} // temp move to see if this would leave or put own king in check

			if isInCheck, _ := b.kingIsInCheck(piece.Colour); !isInCheck {
				if targetSquareOccupied {
					pieceAtTargetSquare.InPlay = true // undo temp take
				}
				ref.CurrentSquare = realCurrentSquare // undo temp move
				return true
			}
			if targetSquareOccupied {
				pieceAtTargetSquare.InPlay = true // undo temp take
			}
			ref.CurrentSquare = realCurrentSquare // undo temp move			return true
		}
	}
	return false
}
func couldAnyBlockCheckAt(pieces []Piece, targetColumn string, targetRow int, b *Board) bool {
	for _, piece := range pieces {
		if piece.couldMoveTo(targetColumn, targetRow, b) {
			_, ref := b.GetPieceAtSquare(piece.CurrentSquare.Column, piece.CurrentSquare.Row)
			realCurrentSquare := ref.CurrentSquare
			ref.CurrentSquare = Square{Column: targetColumn, Row: targetRow} // temp move before checking if king would still be in check..
			if isInCheck, _ := b.kingIsInCheck(piece.Colour); !isInCheck {
				ref.CurrentSquare = realCurrentSquare // undo temp move
				return true
			}
			ref.CurrentSquare = realCurrentSquare // undo temp move
		}
	}
	return false
}
func (b *Board) kingIsInCheck(colour Colour) (bool, []Piece) {
	king := b.getKing(colour)
	dryRun := true
	isCheck := false
	enemies := make([]Piece, 0)

	if colour == White {
		for _, piece := range b.BlackPieces {
			if piece.InPlay {
				if _, err := piece.Move(king.CurrentSquare.Column, king.CurrentSquare.Row, b, dryRun); err == nil {
					// fmt.Printf("%v %v can move to %v%v\n", piece.Colour, piece.Type, king.CurrentSquare.Column, king.CurrentSquare.Row)
					isCheck = true
					enemies = append(enemies, piece)
				}
			}
		}
	} else {
		for _, piece := range b.WhitePieces {
			if piece.InPlay {
				if _, moveFailedErr := piece.Move(king.CurrentSquare.Column, king.CurrentSquare.Row, b, dryRun); moveFailedErr == nil {
					//fmt.Printf("%v %v is in check by %v %v at %v %v\n", king.Colour, king.Type, piece.Colour, piece.Type, piece.CurrentSquare.Column, piece.CurrentSquare.Row)
					isCheck = true
					enemies = append(enemies, piece)
				}
			}
		}
	}
	return isCheck, enemies
}
func (b *Board) checkPathForOccupiedSquaresStraightLeft(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from current position, check if any pieces in the way
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	var squaresInBetween []Square
	for i := currentColumnIndex - 1; i > columnIndex; i-- { // always at least one square between
		occupied, _ := b.GetPieceAtSquare(b.columns[i], targetRow)
		if occupied {
			return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
		}
		squaresInBetween = append(squaresInBetween, Square{Column: b.columns[i], Row: targetRow})
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForSquaresUnderAttackStraightLeft(targetColumn string, targetRow int, p *Piece) error {
	// starting from current position, check if any square is under attack
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	for i := currentColumnIndex - 1; i > columnIndex; i-- { // always at least one square between
		var activeEnemies []Piece
		if p.Colour == White {
			activeEnemies = b.BlackPieces
		} else {
			activeEnemies = b.WhitePieces
		}
		for _, enemy := range activeEnemies {
			if enemy.InPlay {
				if enemy.couldMoveTo(b.columns[i], targetRow, b) {
					return fmt.Errorf("enemy %v can go to %v%v", p.Type, targetColumn, targetRow)
				}
			}
		}
	}
	return nil
}
func (b *Board) checkPathForOccupiedSquaresStraigthRight(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from current position, check if any pieces in the way
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	var squaresInBetween []Square
	for i := currentColumnIndex + 1; i < columnIndex; i++ { // always at least one square between
		occupied, _ := b.GetPieceAtSquare(b.columns[i], targetRow)
		if occupied {
			return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
		}
		squaresInBetween = append(squaresInBetween, Square{Column: b.columns[i], Row: targetRow})
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForSquaresUnderAttackStraightRight(targetColumn string, targetRow int, p *Piece) error {
	// starting from current position, check if any square is under attack
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	for i := currentColumnIndex + 1; i < columnIndex; i++ { // always at least one square between
		var activeEnemies []Piece
		if p.Colour == White {
			activeEnemies = b.BlackPieces
		} else {
			activeEnemies = b.WhitePieces
		}
		for _, enemy := range activeEnemies {
			if enemy.InPlay {
				if enemy.couldMoveTo(b.columns[i], targetRow, b) {
					return fmt.Errorf("enemy %v can go to %v%v", p.Type, targetColumn, targetRow)
				}
			}
		}
	}
	return nil
}
func (b *Board) checkPathForOccupiedSquaresStraightDown(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from current position, check if any pieces in the way
	var squaresInBetween []Square
	for i := p.CurrentSquare.Row - 1; i > targetRow; i-- { // always at least one square between
		occupied, _ := b.GetPieceAtSquare(targetColumn, i)
		if occupied {
			return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
		}
		squaresInBetween = append(squaresInBetween, Square{Column: targetColumn, Row: i})
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForOccupiedSquaresStraightUp(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from current position, check if any pieces in the way
	var squaresInBetween []Square
	for i := p.CurrentSquare.Row + 1; i < targetRow; i++ { // always at least one square between
		occupied, _ := b.GetPieceAtSquare(targetColumn, i)
		if occupied {
			return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
		}
		squaresInBetween = append(squaresInBetween, Square{Column: targetColumn, Row: i})
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForOccupiedSquaresUpRight(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from TARGET position moving backwards, check if any pieces in the way
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	checkRow := targetRow
	var squaresInBetween []Square
	// move one column to the left and one row down at a time and check if square occupied (from target to starting position)
	for i := columnIndex; i > currentColumnIndex; i-- {
		checkRow = checkRow - 1
		index := i - 1 // skip first (target) column, only check if we are jumping over pieces
		if checkRow > p.CurrentSquare.Row {
			occupied, by := b.GetPieceAtSquare(b.columns[index], checkRow)
			if occupied {
				return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces %v", p.Type, targetColumn, targetRow, by)
			}
			squaresInBetween = append(squaresInBetween, Square{Column: b.columns[index], Row: checkRow})
		}
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForOccupiedSquaresDownRight(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from TARGET position moving backwards, check if any pieces in the way
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	checkRow := targetRow
	var squaresInBetween []Square
	// move one column to the left and one row up at a time and check if square occupied (from target to starting position)
	for i := columnIndex; i > currentColumnIndex; i-- {
		checkRow = checkRow + 1
		index := i - 1
		if checkRow < p.CurrentSquare.Row {
			occupied, _ := b.GetPieceAtSquare(b.columns[index], checkRow)
			if occupied {
				return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
			}
			squaresInBetween = append(squaresInBetween, Square{Column: b.columns[index], Row: checkRow})
		}
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForOccupiedSquaresUpLeft(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from TARGET position moving backwards, check if any pieces in the way
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	checkRow := targetRow
	var squaresInBetween []Square
	// move one column to the right and one row down at a time and check if square occupied (from target to starting position)
	for i := columnIndex; i < currentColumnIndex; i++ {
		checkRow = checkRow - 1
		index := i + 1
		if checkRow > p.CurrentSquare.Row {
			occupied, _ := b.GetPieceAtSquare(b.columns[index], checkRow)
			if occupied {

				return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
			}
			squaresInBetween = append(squaresInBetween, Square{Column: b.columns[index], Row: checkRow})
		}
	}
	return squaresInBetween, nil
}
func (b *Board) checkPathForOccupiedSquaresDownLeft(targetColumn string, targetRow int, p *Piece) ([]Square, error) {
	// starting from TARGET position moving backwards, check if any pieces in the way
	currentColumnIndex := b.getColumnIndex(p.CurrentSquare.Column)
	columnIndex := b.getColumnIndex(targetColumn)
	checkRow := targetRow
	var squaresInBetween []Square
	// move one column to the right and one row up at a time and check if square occupied (from target to starting position)
	for i := columnIndex; i < currentColumnIndex; i++ {
		checkRow = checkRow + 1
		index := i + 1
		if checkRow < p.CurrentSquare.Row {
			occupied, _ := b.GetPieceAtSquare(b.columns[index], checkRow)
			if occupied {
				return nil, fmt.Errorf("%v cant move to square %v%v, cannot jump over other pieces", p.Type, targetColumn, targetRow)
			}
			squaresInBetween = append(squaresInBetween, Square{Column: b.columns[index], Row: checkRow})
		}
	}
	return squaresInBetween, nil
}
func (b *Board) getSquare(column string, row int) (Square, error) {

	for _, square := range b.Squares {
		if square.Column == column && square.Row == row {
			return square, nil
		}
	}
	return Square{}, errors.New("Square not found")
}
func (b *Board) getKing(colour Colour) *Piece {
	if colour == White {
		for i, piece := range b.WhitePieces {
			if piece.Type == king {
				return &b.WhitePieces[i]
			}
		}
	} else {
		for i, piece := range b.BlackPieces {
			if piece.Type == king {
				return &b.BlackPieces[i]
			}
		}
	}
	return &Piece{}
}
func (b *Board) getColumnValue(targetColumn string) int {
	var columnValue int
	for i, c := range b.columns {
		if c == targetColumn {
			columnValue = (i + 1) // A = 1, H = 8
		}
	}
	return columnValue
}
func (b *Board) getColumnIndex(targetColumn string) int {
	var columnValue int
	for i, c := range b.columns {
		if c == targetColumn {
			columnValue = i // A = 0, H = 7
		}
	}
	return columnValue
}
func (b *Board) getColumnStringByIndex(columnIndex int) string {
	return b.columns[columnIndex]
}
func (board *Board) placePiecesOnBoard() {

	isAlive := true
	hasMoved := false
	for _, square := range board.Squares {

		if square.Row == 1 && (square.Column == "A" || square.Column == "H") { // place rooks
			// place a white rook on the first square of the A and H columns
			board.WhitePieces = append(board.WhitePieces, Piece{rook, square, White, isAlive, hasMoved})
		} else if square.Row == 8 && (square.Column == "A" || square.Column == "H") {
			// place a black rook on the eighth square of the A and H columns
			board.BlackPieces = append(board.BlackPieces, Piece{rook, square, Black, isAlive, hasMoved})
		} else if square.Row == 1 && (square.Column == "B" || square.Column == "G") { // place knight
			// place a white knight on the first square of the B and G columns
			board.WhitePieces = append(board.WhitePieces, Piece{knight, square, White, isAlive, hasMoved})
		} else if square.Row == 8 && (square.Column == "B" || square.Column == "G") {
			// place a black knight on the eighth square of the B and G columns
			board.BlackPieces = append(board.BlackPieces, Piece{knight, square, Black, isAlive, hasMoved})
		} else if square.Row == 1 && (square.Column == "C" || square.Column == "F") { // place bishops
			// place a white bishop on the first square of the C and F columns
			board.WhitePieces = append(board.WhitePieces, Piece{bishop, square, White, isAlive, hasMoved})
		} else if square.Row == 8 && (square.Column == "C" || square.Column == "F") {
			// place a black bishop on the eighth square of the C and F columns
			board.BlackPieces = append(board.BlackPieces, Piece{bishop, square, Black, isAlive, hasMoved})
		} else if square.Row == 1 && square.Column == "D" { // place queen
			// place a white queen on the first square of the D column
			board.WhitePieces = append(board.WhitePieces, Piece{queen, square, White, isAlive, hasMoved})
		} else if square.Row == 8 && square.Column == "D" {
			// place a black queen on the eighth square of the D column
			board.BlackPieces = append(board.BlackPieces, Piece{queen, square, Black, isAlive, hasMoved})
		} else if square.Row == 1 && square.Column == "E" { // place king
			// place a white king on the first square of the E column
			board.WhitePieces = append(board.WhitePieces, Piece{king, square, White, isAlive, hasMoved})
		} else if square.Row == 8 && square.Column == "E" {
			// place a black king on the eighth square of the E column
			board.BlackPieces = append(board.BlackPieces, Piece{king, square, Black, isAlive, hasMoved})
		} else if square.Row == 2 { // place pawns
			// place a white pawn on the second square of every column
			board.WhitePieces = append(board.WhitePieces, Piece{pawn, square, White, isAlive, hasMoved})
		} else if square.Row == 7 {
			// place a black pawn on the seventh square of every column
			board.BlackPieces = append(board.BlackPieces, Piece{pawn, square, Black, isAlive, hasMoved})
		}
	}
}

// returns the piece at the given square
func (b *Board) GetPieceAtSquare(column string, row int) (bool, *Piece) {
	for i, piece := range b.WhitePieces {
		if piece.CurrentSquare.Column == strings.ToUpper(column) && piece.CurrentSquare.Row == row && piece.InPlay {
			return true, &b.WhitePieces[i]
		}
	}
	for i, piece := range b.BlackPieces {
		if piece.CurrentSquare.Column == strings.ToUpper(column) && piece.CurrentSquare.Row == row && piece.InPlay {
			return true, &b.BlackPieces[i]
		}
	}
	return false, &Piece{}
}
func addToListIfValidSquare(p *Piece, b *Board, validSquares []Square, row int, columnIndex int) []Square {
	if columnIndex < len(b.columns) && columnIndex >= 0 &&
		row <= 8 && row >= 0 {
		s, err := b.getSquare(b.columns[columnIndex], row)
		if err == nil {
			validSquares = append(validSquares, s)
		}
	}
	return validSquares
}
func (b *Board) getMovesFor(p *Piece) map[Move]*MoveResult {
	moves := map[Move]*MoveResult{}
	switch p.Type {
	case pawn:
		for move, result := range p.getValidPawnMoves(b) {
			moves[move] = result
		}
	case knight:
		for move, result := range p.getValidKnightMoves(b) {
			moves[move] = result
		}
	case bishop:
		for move, result := range p.getValidBishopMoves(b) {
			moves[move] = result
		}
	case rook:
		for move, result := range p.getValidRookMoves(b) {
			moves[move] = result
		}
	case queen:
		for move, result := range p.getValidQueenMoves(b) {
			moves[move] = result
		}
	case king:
		for move, result := range p.getValidKingMoves(b) {
			moves[move] = result
		}

	}
	return moves
}

// returns a map of all valid moves for a player/colour
func (b *Board) GetAllMovesFor(Colour Colour) map[Move]*MoveResult {
	moves := map[Move]*MoveResult{}
	if Colour == White {
		for _, piece := range b.WhitePieces {
			m := b.getMovesFor(&piece)
			for k := range m {
				moves[k] = m[k]
			}

		}
	} else {
		for _, piece := range b.BlackPieces {
			m := b.getMovesFor(&piece)
			for k := range m {
				moves[k] = m[k]
			}
		}
	}
	return moves
}
