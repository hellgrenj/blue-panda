package chess

import (
	"testing"
)

func TestMoveBishop_can_move_diagonally(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (move pawn out of the way to test the bishop)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
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
	_, WB := board.GetPieceAtSquare("F", 1)
	_, err := WB.Move("B", 5, board, false)
	if err != nil {
		t.Errorf("Failed to move the bishop horizontally from F1 to B5, %v", err.Error())
	}
	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  BN  BR
	bP  bP  bP  bP  bP  bP  bP  bP
	..  ..  ..  ..  ..  ..  ..  ..
	..  WB  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  wP  ..  ..  ..
	..  ..  ..  ..  ..  ..  ..  ..
	wP  wP  wP  wP  ..  wP  wP  wP
	WR  WN  WB  WQ  WK  ..  WN  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}

func TestMoveBishop_can_NOT_move_straight(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (get bishop out to B5)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  WB  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  ..  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove2 := Move{Square{Column: "F", Row: 1}, Square{Column: "B", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WB := board.GetPieceAtSquare("B", 5)
	_, err := WB.Move("B", 4, board, false)
	if err == nil {
		t.Errorf("Failed to move the bishop horizontally from F1 to B5, %v", err.Error())
	}
	expectedErrorMsg := "bishops can only move diagonally"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error message to be '%v', but got '%v'", expectedErrorMsg, err.Error())
	}

}

func TestGetValidBishopMoves(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (get bishop out to B5)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  WB  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  ..  wP  wP  wP
	// WR  WN  WB  WQ  WK  ..  WN  WR
	whiteMove1 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 4}}
	whiteMove2 := Move{Square{Column: "F", Row: 1}, Square{Column: "B", Row: 5}}
	moves := []Move{whiteMove1, whiteMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, WB := board.GetPieceAtSquare("B", 5)
	peeks := WB.getValidBishopMoves(board)
	if len(peeks) != 8 {
		t.Errorf("Expected 8 valid moves for white bishop at A5, got %v", len(peeks))
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
	if numberOfTakesFound != 1 {
		t.Errorf("Expected 1 take action for white pawn at A4, got %v", numberOfTakesFound)
	}
	if numberOfGoTosFound != 7 {
		t.Errorf("Expected 1 goto action for white pawn at A4, got %v", numberOfGoTosFound)
	}
	if len(potentialTakes) != 1 {
		t.Errorf("Expected 1 potential take for white pawn at A4, got %v", len(potentialTakes))
	}
	if potentialTakes[0].GetValue() != 1 {
		t.Errorf("Expected potential take to be worth 1, got %v", potentialTakes[0].GetValue())
	}
	if potentialTakes[0].Colour != Black || potentialTakes[0].Type != pawn ||
		potentialTakes[0].CurrentSquare.Column != "D" || potentialTakes[0].CurrentSquare.Row != 7 {
		t.Errorf("Expected potential take to be a black pawn at E4, got %v", potentialTakes[0])

	}
}
