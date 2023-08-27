package chess

import (
	"fmt"
	"testing"
)

func TestMoveKing_can_move_1_square_in_any_direction(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wp  ..  ..  ..
	// ..  ..  ..  ..  WK  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  ..  WB  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove2 := Move{Square{Column: "E", Row: 1}, Square{Column: "E", Row: 2}}
	whiteMove3 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 3}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WK := board.GetPieceAtSquare("e", 3) // small e => capitalized to E in Move function..
	_, err := WK.Move("F", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the king horizontally from E3 to F3, %v", err.Error())
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  ♙  .  .  .
	.  .  .  .  .  ♔  .  .
	♙  ♙  ♙  ♙  .  ♙  ♙  ♙
	♖  ♘  ♗  ♕  .  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	_, WK = board.GetPieceAtSquare("F", 3)
	_, err = WK.Move("F", 4, board, false)
	if err != nil {
		t.Errorf("Failed to move the king vertically from F3 to F4, %v", err.Error())
	}

	expectedStateOfBoard = `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  ♙  ♔  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  .  ♙  ♙  ♙
	♖  ♘  ♗  ♕  .  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	_, WK = board.GetPieceAtSquare("F", 4)
	_, err = WK.Move("F", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the king back vertically from F4 to F3, %v", err.Error())
	}

	expectedStateOfBoard = `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  ♙  .  .  .
	.  .  .  .  .  ♔  .  .
	♙  ♙  ♙  ♙  .  ♙  ♙  ♙
	♖  ♘  ♗  ♕  .  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	_, WK = board.GetPieceAtSquare("F", 3)
	_, err = WK.Move("E", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the king back horizontally from F3 to E3, %v", err.Error())
	}

	expectedStateOfBoard = `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  ♙  .  .  .
	.  .  .  .  ♔  .  .  .
	♙  ♙  ♙  ♙  .  ♙  ♙  ♙
	♖  ♘  ♗  ♕  .  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}

func TestKingTryRun_can_outrun_a_check(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  ..  BK  BB  BN  BR
	// bP  bP  ..  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// BQ  ..  bP  ..  ..  ..  ..  ..
	// ..  ..  ..  wP  wP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  ..  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove2 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	blackMove1 := Move{Square{Column: "C", Row: 7}, Square{Column: "C", Row: 5}}
	blackMove2 := Move{Square{Column: "D", Row: 8}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2, blackMove1, blackMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	kingInCheck, enemies := board.kingIsInCheck(White)
	if !kingInCheck {
		t.Errorf("Expected the king to be in check, but it wasn't")
		return
	}
	if len(enemies) != 1 {
		t.Errorf("Expected 1 enemy, but got %v", len(enemies))
		return
	}
	if fmt.Sprintf("%v%v", enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row) != "A5" {
		t.Errorf("Expected the enemy to be at A5, but got %v%v", enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row)
		return
	}
	if enemies[0].Type != queen {
		t.Errorf("Expected the enemy to be a queen, but got %v", enemies[0].Type)
		return
	}
	// black queen on A5 is checking the white king on E1
	king := board.getKing(White)
	if err := king.kingTryRun(board); err != nil { // king should be able to outrun the check to E2
		t.Errorf("TryKingRun failed, %v", err.Error())
	}

}
func TestKingTryRun_BLACK_can_outrun_a_check_when_on_edge_of_board(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  BQ  ..  BB  BN  BR
	// bP  bP  bP  bP  bP  ..  bP  bP
	// ..  ..  ..  ..  ..  bP  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  BK
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  wP  wP  wP  ..
	// WR  WN  WB  WQ  WK  WB  WN  WR

	blackMove1 := Move{Square{Column: "F", Row: 7}, Square{Column: "F", Row: 6}}
	blackMove2 := Move{Square{Column: "E", Row: 8}, Square{Column: "F", Row: 7}}
	blackMove3 := Move{Square{Column: "F", Row: 7}, Square{Column: "G", Row: 6}}
	blackMove4 := Move{Square{Column: "G", Row: 6}, Square{Column: "H", Row: 5}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	// now simulate that wp h2 is taken (out of play)
	_, h2Pawn := board.GetPieceAtSquare("H", 2)
	h2Pawn.InPlay = false // taken..
	// now BK is on the edge of the board and in check by WR but it should be able to run to g6, g5, g4

	kingInCheck, enemies := board.kingIsInCheck(Black)
	if !kingInCheck {
		t.Errorf("Expected the king to be in check, but it wasn't")
		return
	}
	if len(enemies) != 1 {
		t.Errorf("Expected 1 enemy, but got %v", len(enemies))
		return
	}
	if fmt.Sprintf("%v%v", enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row) != "H1" {
		t.Errorf("Expected the enemy to be at A5, but got %v%v", enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row)
		return
	}
	if enemies[0].Type != rook {
		t.Errorf("Expected the enemy to be a rook, but got %v", enemies[0].Type)
		return
	}
	// wr on h1 is checking the black king on h5
	king := board.getKing(Black)
	if err := king.kingTryRun(board); err != nil { // king should be able to outrun the check to E2
		t.Errorf("TryKingRun failed, %v", err.Error())
	}

}

func TestKingTryRun_returns_error_if_king_cant_run_anywhere(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  ..  BK  BB  BN  BR
	// bP  bP  ..  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// BQ  ..  bP  ..  ..  ..  ..  ..
	// ..  ..  ..  wP  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  ..  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	blackMove1 := Move{Square{Column: "C", Row: 7}, Square{Column: "C", Row: 5}}
	blackMove2 := Move{Square{Column: "D", Row: 8}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, blackMove1, blackMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	kingInCheck, enemies := board.kingIsInCheck(White)
	if !kingInCheck {
		t.Errorf("Expected the king to be in check, but it wasn't")
		return
	}
	if len(enemies) != 1 {
		t.Errorf("Expected 1 enemy, but got %v", len(enemies))
		return
	}
	if fmt.Sprintf("%v%v", enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row) != "A5" {
		t.Errorf("Expected the enemy to be at A5, but got %v%v", enemies[0].CurrentSquare.Column, enemies[0].CurrentSquare.Row)
		return
	}
	if enemies[0].Type != queen {
		t.Errorf("Expected the enemy to be a queen, but got %v", enemies[0].Type)
		return
	}
	// black queen on A5 is checking the white king on E1 AND King is boxed in, can only move to D2 but that doesnt remove the check
	king := board.getKing(White)
	if err := king.kingTryRun(board); err == nil {
		t.Errorf("TryKingRun should have failed, but it didn't")
	}

}

func TestGetValidKingMovdes(t *testing.T) {
	//defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wp  ..  ..  ..
	// ..  ..  ..  ..  WK  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  ..  WB  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove2 := Move{Square{Column: "E", Row: 1}, Square{Column: "E", Row: 2}}
	whiteMove3 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 3}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WK := board.GetPieceAtSquare("e", 3) // small e => capitalized to E in Move function..
	peeks := WK.getValidKingMoves(board)
	if len(peeks) != 5 {
		t.Errorf("Expected 5 valid moves, but got %v", len(peeks))
		return
	}

	numberOfTakesFound := 0
	numberOfGoTosFound := 0
	potentialTakes := []*Piece{}
	for k := range peeks {
		if peeks[k].Action == Take {
			numberOfTakesFound++
			potentialTakes = append(potentialTakes, peeks[k].Piece)
		}
		if peeks[k].Action == GoTo {
			numberOfGoTosFound++
		}
	}
	if numberOfTakesFound != 0 {
		t.Errorf("Expected 0 take action for white King at E3, got %v", numberOfTakesFound)
	}
	if numberOfGoTosFound != 5 {
		t.Errorf("Expected 5 goto action for white King at E3, got %v", numberOfGoTosFound)
	}
	if len(potentialTakes) != 0 {
		t.Errorf("Expected 0 potential take for white King at E3, got %v", len(potentialTakes))
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "E", Row: 3}, To: Square{Column: "D", Row: 3}}) {
		t.Errorf("Expected one of the valid moves to be E3-D3, but it wasn't")
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "E", Row: 3}, To: Square{Column: "D", Row: 4}}) {
		t.Errorf("Expected one of the valid moves to be E3-D4, but it wasn't")
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "E", Row: 3}, To: Square{Column: "F", Row: 3}}) {
		t.Errorf("Expected one of the valid moves to be E3-F3, but it wasn't")
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "E", Row: 3}, To: Square{Column: "F", Row: 4}}) {
		t.Errorf("Expected one of the valid moves to be E3-F4, but it wasn't")
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "E", Row: 3}, To: Square{Column: "E", Row: 2}}) {
		t.Errorf("Expected one of the valid moves to be E3-E2, but it wasn't")
	}

}

func TestMoveKing_cant_move_into_check(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  ..  BK  BB  BN  BR
	// bP  bP  ..  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// BQ  ..  bP  ..  ..  ..  ..  ..
	// ..  ..  ..  wP  wP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  ..  WK  wP  wP  wP
	// WR  WN  WB  WQ  ..  WB  WN  WR
	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	whiteMove2 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove3 := Move{Square{Column: "E", Row: 1}, Square{Column: "E", Row: 2}}
	blackMove1 := Move{Square{Column: "C", Row: 7}, Square{Column: "C", Row: 5}}
	blackMove2 := Move{Square{Column: "D", Row: 8}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, blackMove1, blackMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WK := board.GetPieceAtSquare("e", 2)
	_, err := WK.Move("E", 1, board, false)

	if err == nil {
		t.Errorf("Expected an error when moving the King into check, but got none")
		return
	}

}

func TestMoveKingBlack_cant_move_into_check(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bp  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  bP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  WB  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  wp  ..  ..  ..  ..
	// wP  wP  wP  ..  wP  wP  wP  wP
	// WR  WN  ..  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 3}}
	whiteMove2 := Move{Square{Column: "C", Row: 1}, Square{Column: "G", Row: 5}}
	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 6}}
	moves := []Move{whiteMove1, whiteMove2, blackMove1}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, BK := board.GetPieceAtSquare("e", 8)
	_, err := BK.Move("E", 7, board, false)

	if err == nil {
		t.Errorf("Expected an error when moving the King into check, but got none")
		return
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  .  ♟  ♟  ♟
	.  .  .  .  ♟  .  .  .
	.  .  .  .  .  .  ♗  .
	.  .  .  .  .  .  .  .
	.  .  .  ♙  .  .  .  .
	♙  ♙  ♙  .  ♙  ♙  ♙  ♙
	♖  ♘  .  ♕  ♔  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}

func TestMoveKingBlack_cant_move_into_check_by_pawn(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  BQ  ..  BB  BN  BR
	// bP  bP  bp  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  bP  BK  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 6}}
	blackMove2 := Move{Square{Column: "E", Row: 8}, Square{Column: "E", Row: 7}}
	blackMove3 := Move{Square{Column: "E", Row: 7}, Square{Column: "F", Row: 6}}
	moves := []Move{whiteMove1, blackMove1, blackMove2, blackMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, BK := board.GetPieceAtSquare("F", 6)
	_, err := BK.Move("F", 5, board, false)

	if err == nil {
		t.Errorf("Expected an error when moving the King into check, but got none")
		return
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  .  ♝  ♞  ♜
	♟  ♟  ♟  ♟  .  ♟  ♟  ♟
	.  .  .  .  ♟  ♚  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  ♙  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  .  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}

func TestMoveKingBlack_cant_move_into_check_by_rook(t *testing.T) {
	board := newBoard()
	// BR  BN  BB  BQ  ..  BB  BN  BR
	// bP  bP  bp  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  bP  BK  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  WR  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 1}, Square{Column: "A", Row: 3}}
	whiteMove3 := Move{Square{Column: "A", Row: 3}, Square{Column: "E", Row: 3}}
	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 6}}
	blackMove2 := Move{Square{Column: "E", Row: 8}, Square{Column: "E", Row: 7}}
	blackMove3 := Move{Square{Column: "E", Row: 7}, Square{Column: "F", Row: 6}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, blackMove1, blackMove2, blackMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, BK := board.GetPieceAtSquare("F", 6)
	_, err := BK.Move("E", 5, board, false)

	if err == nil {
		t.Errorf("Expected an error when moving the King into check, but got none")
		return
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  .  ♝  ♞  ♜
	♟  ♟  ♟  ♟  .  ♟  ♟  ♟
	.  .  .  .  ♟  ♚  .  .
	.  .  .  .  .  .  .  .
	♙  .  .  .  .  .  .  .
	.  .  .  .  ♖  .  .  .
	.  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	.  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}

func TestMoveKing_white_can_castle_kingside(t *testing.T) {
	board := newBoard()

	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bp  bP  bp  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  wP  wP  ..  ..  ..
	// ..  ..  ..  WB  ..  WN  ..  ..
	// wP  wP  wP  ..  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  ..  ..  WR

	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	whiteMove2 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove3 := Move{Square{Column: "F", Row: 1}, Square{Column: "D", Row: 3}}
	whiteMove4 := Move{Square{Column: "G", Row: 1}, Square{Column: "F", Row: 3}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, whiteMove4}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WK := board.GetPieceAtSquare("E", 1)
	// try castling kingside
	_, err := WK.Move("G", 1, board, false)

	if err != nil {
		t.Errorf("Expected no error when castling kingside, but got %v", err.Error())
		return
	}
	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  ♙  ♙  .  .  .
	.  .  .  ♗  .  ♘  .  .
	♙  ♙  ♙  .  .  ♙  ♙  ♙
	♖  ♘  ♗  ♕  .  ♖  ♔  .
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}
func TestMoveKing_white_can_castle_queen_side(t *testing.T) {

	expectedInitState := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  ♙  ♙  .  .  .
	♘  .  .  .  ♗  .  .  .
	♙  ♙  ♙  ♕  .  ♙  ♙  ♙
	♖  .  .  .  ♔  ♗  ♘  ♖
	`
	board := newBoard()

	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	whiteMove2 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove3 := Move{Square{Column: "C", Row: 1}, Square{Column: "E", Row: 3}}
	whiteMove4 := Move{Square{Column: "B", Row: 1}, Square{Column: "A", Row: 3}}
	whiteMove5 := Move{Square{Column: "D", Row: 1}, Square{Column: "D", Row: 2}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, whiteMove4, whiteMove5}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, WK := board.GetPieceAtSquare("E", 1)
	// try castling queenside
	_, err := WK.Move("C", 1, board, false)

	if err != nil {
		t.Errorf("Expected no error when castling queenside, but got %v", err.Error())
		return
	}
	expectedEndState := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  ♙  ♙  .  .  .
	♘  .  .  .  ♗  .  .  .
	♙  ♙  ♙  ♕  .  ♙  ♙  ♙
	.  .  ♔  ♖  .  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedEndState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}

func TestMoveKing_white_can_castle_queen_side_even_if_kingside_rook_gone(t *testing.T) {

	expectedInitState := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  ♙  ♙  .  .  .
	♘  .  .  .  ♗  .  .  .
	♙  ♙  ♙  ♕  .  ♙  ♙  ♙
	♖  .  .  .  ♔  ♗  ♘  .
	`
	board := newBoard()

	// remove white kingside rook
	_, rook := board.GetPieceAtSquare("H", 1)
	rook.InPlay = false

	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	whiteMove2 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove3 := Move{Square{Column: "C", Row: 1}, Square{Column: "E", Row: 3}}
	whiteMove4 := Move{Square{Column: "B", Row: 1}, Square{Column: "A", Row: 3}}
	whiteMove5 := Move{Square{Column: "D", Row: 1}, Square{Column: "D", Row: 2}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, whiteMove4, whiteMove5}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, WK := board.GetPieceAtSquare("E", 1)
	// try castling queenside
	_, err := WK.Move("C", 1, board, false)

	if err != nil {
		t.Errorf("Expected no error when castling queenside, but got %v", err.Error())
		return
	}
	expectedEndState := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  ♙  ♙  .  .  .
	♘  .  .  .  ♗  .  .  .
	♙  ♙  ♙  ♕  .  ♙  ♙  ♙
	.  .  ♔  ♖  .  ♗  ♘  .
	`

	if err := assertExpectedBoardState(expectedEndState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}
func TestMove_black_can_castle_kingside(t *testing.T) {

	expectedInitState := `
	♜  ♞  ♝  ♛  ♚  .  .  ♜
	♟  ♟  ♟  ♟  .  ♟  ♟  ♟
	.  .  .  ♝  .  .  .  ♞
	.  .  .  .  ♟  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 5}}
	blackMove2 := Move{Square{Column: "F", Row: 8}, Square{Column: "D", Row: 6}}
	blackMove3 := Move{Square{Column: "G", Row: 8}, Square{Column: "H", Row: 6}}
	moves := []Move{blackMove1, blackMove2, blackMove3}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// try castling kingside
	_, err := BK.Move("G", 8, board, false)
	if err != nil {
		t.Errorf("Expected no error when castling kingside (black), but got %v", err.Error())
		return
	}
	expectedEndState := `
	♜  ♞  ♝  ♛  .  ♜  ♚  .
	♟  ♟  ♟  ♟  .  ♟  ♟  ♟
	.  .  .  ♝  .  .  .  ♞
	.  .  .  .  ♟  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	if err := assertExpectedBoardState(expectedEndState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}

func TestMove_black_can_castle_queen_side(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♛  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  ♟  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 7}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// try castling queenside
	_, err := BK.Move("C", 8, board, false)
	if err != nil {
		t.Errorf("Expected no error when castling kingside (black), but got %v", err.Error())
		return
	}
	expectedEndState := `
	.  .  ♚  ♜  .  ♝  ♞  ♜
	♟  ♟  ♟  ♛  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  ♟  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	if err := assertExpectedBoardState(expectedEndState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}

func TestMove_black_cant_castle_queen_side_when_in_check(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  ♟  .  ♟  ♟  ♟  ♟
	♞  .  .  ♛  ♝  .  .  .
	.  .  .  ♟  .  .  .  .
	♕  .  .  .  .  .  .  .
	.  .  ♙  .  .  .  .  .
	♙  ♙  .  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  .  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 6}}
	whiteMove1 := Move{Square{Column: "C", Row: 2}, Square{Column: "C", Row: 3}}
	whiteMove2 := Move{Square{Column: "D", Row: 1}, Square{Column: "A", Row: 4}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4, whiteMove1, whiteMove2}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// try castling queenside
	_, err := BK.Move("C", 8, board, false)
	if err == nil {
		t.Errorf("Expected error when castling queenside (black) because king is in check, but got none")
		return
	}
	if err.Error() != "castling not allowed, err: king is in check" {
		t.Errorf("Expected 'castling not allowed, err: king is in check', but got %v", err.Error())
		return
	}
}

func TestMove_black_cant_castle_queen_side_passing_through_check(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  ♟  .  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  ♕  .  .  .  .
	.  .  ♙  .  .  .  .  .
	♙  ♙  .  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  .  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 6}}
	whiteMove1 := Move{Square{Column: "C", Row: 2}, Square{Column: "C", Row: 3}}
	whiteMove2 := Move{Square{Column: "D", Row: 1}, Square{Column: "A", Row: 4}}
	// move white queen to D4 to attack D8 blocking the black king from castling queenside
	// after we simulate the black pawn at D5 was taken and the black queen at D6 was taken
	whiteMove3 := Move{Square{Column: "A", Row: 4}, Square{Column: "D", Row: 4}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4, whiteMove1, whiteMove2, whiteMove3}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	// simulate black pawn at D5 was taken
	_, pd5 := board.GetPieceAtSquare("D", 5)
	pd5.InPlay = false
	// simulate black queen at D6 was taken
	_, qd6 := board.GetPieceAtSquare("D", 6)
	qd6.InPlay = false

	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// try castling queenside
	_, err := BK.Move("C", 8, board, false)
	if err == nil {
		t.Errorf("Expected error when castling queenside (black) because king is in check, but got none")
		return
	}
	if err.Error() != "castling not allowed, err: king passes through a square that is attacked by an enemy piece" {
		t.Errorf("Expected 'castling not allowed, err: king passes through a square that is attacked by an enemy piece', but got %v", err.Error())
		return
	}
}

func TestMove_black_cant_castle_queen_side_if_king_moved(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♛  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  ♟  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 7}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// first move king one square
	_, err := BK.Move("D", 8, board, false)
	if err != nil {
		t.Errorf("Expected no error when moving king, but got %v", err.Error())
		return
	}
	// try castling queenside
	_, err = BK.Move("B", 8, board, false)
	if err == nil {
		t.Errorf("Expected error when castling queenside (black) because king is in check, but got none")
		return
	}
	if err.Error() != "castling not allowed, err: king or rook has previously moved" {
		t.Errorf("Expected 'castling not allowed, err: king or rook has previously moved', but got %v", err.Error())
		return
	}
}

func TestMove_black_cant_castle_queen_side_if_rook_moved_out_of_position(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♛  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  ♟  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 7}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// first move rook one square
	_, BR := board.GetPieceAtSquare("A", 8)
	_, err := BR.Move("B", 8, board, false)
	if err != nil {
		t.Errorf("Expected no error when moving rook, but got %v", err.Error())
		return
	}
	// try castling queenside
	_, err = BK.Move("C", 8, board, false)
	if err == nil {
		t.Errorf("Expected error when castling queenside (black) because king is in check, but got none")
		return
	}
	if err.Error() != "castling not allowed, err: no rook found on queenside" {
		t.Errorf("Expected 'castling not allowed, err: no rook found on queenside', but got %v", err.Error())
		return
	}
}

func TestMove_black_cant_castle_queen_side_if_rook_moved(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♛  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  ♟  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 7}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// first move rook one square
	_, BR := board.GetPieceAtSquare("A", 8)
	_, err := BR.Move("B", 8, board, false)
	// then move it back so it is in position BUT has already moved
	_, err = BR.Move("A", 8, board, false)
	if err != nil {
		t.Errorf("Expected no error when moving rook, but got %v", err.Error())
		return
	}
	// try castling queenside
	_, err = BK.Move("C", 8, board, false)
	if err == nil {
		t.Errorf("Expected error when castling queenside (black) because king is in check, but got none")
		return
	}
	if err.Error() != "castling not allowed, err: king or rook has previously moved" {
		t.Errorf("Expected 'castling not allowed, err: king or rook has previously moved', but got %v", err.Error())
		return
	}
}

func TestMove_black_cant_castle_queen_side_ending_up_in_check(t *testing.T) {
	expectedInitState := `
	♜  .  .  .  ♚  ♝  ♞  ♜
	♟  ♟  .  .  ♟  ♟  ♟  ♟
	♞  .  .  .  ♝  .  .  .
	.  .  .  .  .  .  .  .
	.  .  ♕  .  .  .  .  .
	.  .  ♙  .  .  .  .  .
	♙  ♙  .  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  .  ♔  ♗  ♘  ♖
	`
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 5}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "E", Row: 6}}
	blackMove3 := Move{Square{Column: "B", Row: 8}, Square{Column: "A", Row: 6}}
	blackMove4 := Move{Square{Column: "D", Row: 8}, Square{Column: "D", Row: 6}}
	whiteMove1 := Move{Square{Column: "C", Row: 2}, Square{Column: "C", Row: 3}}
	whiteMove2 := Move{Square{Column: "D", Row: 1}, Square{Column: "A", Row: 4}}
	// move white queen to C4 to attack C8 blocking the black king from castling queenside
	// after we simulate the black pawn at D5 and C7 was taken and the black queen at D6 was taken
	whiteMove3 := Move{Square{Column: "A", Row: 4}, Square{Column: "C", Row: 4}}
	moves := []Move{blackMove1, blackMove2, blackMove3, blackMove4, whiteMove1, whiteMove2, whiteMove3}
	board := newBoard()
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %s", scenarioPrepError.Error())
		return
	}
	// simulate black pawn at D5 was taken
	_, pd5 := board.GetPieceAtSquare("D", 5)
	pd5.InPlay = false

	// simulate black pawn at C6 was taken
	_, pc5 := board.GetPieceAtSquare("C", 7)
	pc5.InPlay = false

	// simulate black queen at D6 was taken
	_, qd6 := board.GetPieceAtSquare("D", 6)
	qd6.InPlay = false

	if err := assertExpectedBoardState(expectedInitState, board); err != nil {
		t.Errorf("Failed to assert expected board state, %s (Visible whitespace is ignored, something else differs!", err.Error())
	}
	_, BK := board.GetPieceAtSquare("E", 8)
	// try castling queenside
	_, err := BK.Move("C", 8, board, false)
	if err == nil {
		t.Errorf("Expected error when castling queenside (black) because king is in check, but got none")
		return
	}
	if err.Error() != "move is not legal" {
		t.Errorf("Expected 'move is not legal', but got %v", err.Error())
		return
	}
}
