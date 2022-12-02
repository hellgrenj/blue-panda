package chess

import (
	"strings"
	"testing"
)

func TestMoveRook_should_be_able_to_move_straight_if_not_blocked(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 4}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WRA1 := board.GetPieceAtSquare("A", 1)
	_, err := WRA1.Move("A", 4, board, false)
	if err != nil {
		t.Errorf("Failed to move the rook, %v", err.Error())
	}
	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  BN  BR
	bP  bP  bP  bP  bP  bP  bP  bP
	..  ..  ..  ..  ..  ..  ..  ..
	wP  ..  ..  ..  ..  ..  ..  ..
	WR  ..  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  ..  ..  ..  ..
	..  wP  wP  wP  wP  wP  wP  wP
	..  WN  WB  WQ  WK  WB  WN  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}

func TestMoveRook_can_not_move_to_square_occupied_by_friendly(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 4}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WRA1 := board.GetPieceAtSquare("A", 1)
	_, err := WRA1.Move("A", 5, board, false) // <--- A5 is occupied by a white pawn right now...
	if err == nil {
		t.Errorf("Should not be able to move to a square occupied by a friendly piece")
	}
	expectedErr := "White Rook cant move to square A5, it is occupied by White Pawn"
	if strings.TrimSpace(err.Error()) != expectedErr {
		t.Errorf("Expected error message %v, got %v", expectedErr, err.Error())
	}
}

func TestMoveRook_cant_jump_over_pieces(t *testing.T) {
	defer quiet()()
	board := newBoard()
	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 4}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WRA1 := board.GetPieceAtSquare("A", 1)
	_, err := WRA1.Move("A", 6, board, false) // <--- A6 is free but rook cant jump over the pawn on A5
	if err == nil {
		t.Errorf("Should not be able to jump over pieces")
	}
	expectedErr := "Rook cant move to square A6, cannot jump over other pieces"

	errMsg := strings.ToUpper(strings.TrimSpace(err.Error()))
	if errMsg != strings.ToUpper(expectedErr) {
		t.Errorf("Expected error message %v, got %v", expectedErr, err.Error())
	}
}

func TestMoveRook_can_not_move_diagonally(t *testing.T) {
	defer quiet()()
	board := newBoard()
	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bp  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 4}, Square{Column: "A", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WRA1 := board.GetPieceAtSquare("A", 1)
	_, err := WRA1.Move("B", 2, board, false) // <--- rook cant move diagonally
	if err == nil {
		t.Errorf("Should not be able to move diagonally")
	}
	expectedErr := "rooks can only move straight"
	if strings.TrimSpace(err.Error()) != expectedErr {
		t.Errorf("Expected error message %v, got %v", expectedErr, err.Error())
	}
}

func TestMoveRook_can_move_horizontally(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bp  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// WR  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// ..  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 4}, Square{Column: "A", Row: 5}}
	whitMove3 := Move{Square{Column: "A", Row: 1}, Square{Column: "A", Row: 3}}
	moves := []Move{whiteMove1, whiteMove2, whitMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WRA1 := board.GetPieceAtSquare("A", 3)
	_, err := WRA1.Move("E", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the rook, %v", err.Error())
	}
	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  BN  BR
	bP  bP  bP  bP  bP  bP  bP  bP
	..  ..  ..  ..  ..  ..  ..  ..
	wP  ..  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  WR  ..  ..  ..
	..  wP  wP  wP  wP  wP  wP  wP
	..  WN  WB  WQ  WK  WB  WN  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

}

func TestGetValidRookMoves(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the rook)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bp  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// WR  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// ..  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 4}, Square{Column: "A", Row: 5}}
	whitMove3 := Move{Square{Column: "A", Row: 1}, Square{Column: "A", Row: 3}}
	moves := []Move{whiteMove1, whiteMove2, whitMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WR := board.GetPieceAtSquare("A", 3)

	peeks := WR.getValidRookMoves(board)
	if len(peeks) != 10 {
		t.Errorf("Expected 10 valid moves, got %v", len(peeks))
	}
}
