package chess

import (
	"testing"
)

func TestMoveQueen(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the queen)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wp  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	moves := []Move{whiteMove1}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WQ := board.GetPieceAtSquare("D", 1)
	_, err := WQ.Move("G", 4, board, false)
	if err != nil {
		t.Errorf("Failed to move the queen from D1 to G4, %v", err.Error())
	}

	_, WQ = board.GetPieceAtSquare("G", 4)
	_, err = WQ.Move("H", 4, board, false)
	if err != nil {
		t.Errorf("Failed to move the queen from G4 to H4, %v", err.Error())
	}
	_, WQ = board.GetPieceAtSquare("H", 4)
	_, err = WQ.Move("H", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the queen from H4 to H3, %v", err.Error())
	}
	_, WQ = board.GetPieceAtSquare("H", 3)
	_, err = WQ.Move("G", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the queen from H3 to G3, %v", err.Error())
	}

	_, WQ = board.GetPieceAtSquare("G", 3)
	_, err = WQ.Move("G", 6, board, false)
	if err != nil {
		t.Errorf("Failed to move the queen from G3 to G6, %v", err.Error())
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  ♕  .
	.  .  .  .  .  .  .  .
	.  .  .  .  ♙  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  ♙  .  ♙  ♙  ♙
	♖  ♘  ♗  .  ♔  ♗  ♘  ♖
	`
	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}
func TestGetValidQueenMoves(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (move pawn out of the way to test the queen)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wp  ..  WQ  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	moves := []Move{whiteMove1}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WQ := board.GetPieceAtSquare("D", 1)
	peeks := WQ.getValidQueenMoves(board)
	if len(peeks) != 4 {
		t.Errorf("Expected 4 peeks, got %v", len(peeks))
	}
	WQ.Move("G", 4, board, false)
	peeks = WQ.getValidQueenMoves(board)
	if len(peeks) != 14 {
		t.Errorf("Expected 14 peeks, got %v", len(peeks))
	}
}
