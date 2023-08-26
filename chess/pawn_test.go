package chess

import (
	"strings"
	"testing"
)

func TestMovePawn_should_be_able_to_move_A2_A3(t *testing.T) {
	defer quiet()()
	board := newBoard()
	_, a2Pawn := board.GetPieceAtSquare("A", 2)
	_, err := a2Pawn.Move("A", 3, board, false)
	if err != nil {
		t.Errorf("Expected pawn to move to A3 successfully")
		expectedErrorMessage := "square A 7 is occupied by &{Pawn bp {A 7} Black true false}"
		if strings.TrimSpace(err.Error()) != expectedErrorMessage {
			t.Errorf("Expected error message %v, got %v", expectedErrorMessage, err.Error())
		}
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	♙  .  .  .  .  .  .  .
	.  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}

func TestMovePawn_cant_take_backwards(t *testing.T) {
	defer quiet()()
	board := newBoard()
	// CREATE START SCENARIO (white pawn cant take black pawn backwards)
	// bR  bN  bB  bQ  bK  bB  bN  bR
	// bp  bp  bp  bp  bp  bp  bp  bp
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  wp  ..  ..  ..  ..
	// ..  ..  ..  ..  bp  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  wp  wp  wp  wp  wp  wp  wp
	// wR  wN  wB  wQ  wK  wB  wN  wR

	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	whiteMove2 := Move{Square{Column: "D", Row: 4}, Square{Column: "D", Row: 5}}

	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 5}}
	blackMove2 := Move{Square{Column: "E", Row: 5}, Square{Column: "E", Row: 4}}
	moves := []Move{whiteMove1, blackMove1, whiteMove2, blackMove2}
	scenarioPrepError := prepScenario(moves, board)

	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, wpD5 := board.GetPieceAtSquare("D", 5)
	_, err := wpD5.Move("E", 4, board, false)
	if err == nil {
		t.Errorf("Expected this move D5-E4 to fail, cannot take backwards")
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  ♟  ♟  ♟  .  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  .  .  ♙  .  .  .  .
	.  .  .  .  ♟  .  .  .
	.  .  .  .  .  .  .  .
	♙  ♙  ♙  .  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}
func TestMovePawn_should_NOT_be_able_to_move_A2_A7(t *testing.T) {
	defer quiet()()
	board := newBoard()
	_, a2Pawn := board.GetPieceAtSquare("A", 2)
	_, err := a2Pawn.Move("A", 7, board, false)
	if err == nil {
		t.Errorf("Expected this move A2-A7 to fail")
	}
	if err != nil {

		expectedErrorMessage := "pawn cant move to square A7, can only move two squares forward on first move and then one square"
		if strings.TrimSpace(err.Error()) != expectedErrorMessage {
			t.Errorf("Expected error message %v, got %v", expectedErrorMessage, err.Error())
		}
	}
}
func TestMovePawn_should_NOT_able_to_move_diagonally_if_not_when_taking(t *testing.T) {
	defer quiet()()
	board := newBoard()
	_, a2Pawn := board.GetPieceAtSquare("A", 2)
	_, err := a2Pawn.Move("B", 3, board, false)
	if err == nil {
		t.Errorf("Expected this move A2-B3 to fail, cannot move diagonally and not take")
	}
	if err != nil {

		expectedErrorMessage := "piece &{Pawn {A 2} White true false} cant move to B 3, can only move diagonally when taking"
		if strings.TrimSpace(err.Error()) != expectedErrorMessage {
			t.Errorf("Expected error message %v, got %v", expectedErrorMessage, err.Error())
		}
	}
}
func TestMovePawn_should_be_able_to_move_diagonally_when_taking(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (white pawn takes black pawn)
	// bR  bN  bB  bQ  bK  bB  bN  bR
	// bp  ..  bp  bp  bp  bp  bp  bp
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  bp  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wp  wp  wp  wp  wp  wp  wp
	// wR  wN  wB  wQ  wK  wB  wN  wR

	// wp from A2 to A4
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	// then bp from B7 to B5
	blackMove1 := Move{Square{Column: "B", Row: 7}, Square{Column: "B", Row: 5}}
	moves := []Move{whiteMove1, blackMove1}
	scenarioPrepError := prepScenario(moves, board)

	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, a4Pawn := board.GetPieceAtSquare("A", 4)
	_, err := a4Pawn.Move("B", 5, board, false)
	if err != nil {
		t.Errorf("Expected pawn to move to B3 successfully")
		expectedErrorMessage := "square B 3 is occupied by &{Pawn bp {B 3} Black true false}"
		if strings.TrimSpace(err.Error()) != expectedErrorMessage {
			t.Errorf("Expected error message %v, got %v", expectedErrorMessage, err.Error())
		}
	}

	expectedStateOfBoard := `
	♜  ♞  ♝  ♛  ♚  ♝  ♞  ♜
	♟  .  ♟  ♟  ♟  ♟  ♟  ♟
	.  .  .  .  .  .  .  .
	.  ♙  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  .  .  .  .  .  .  .
	.  ♙  ♙  ♙  ♙  ♙  ♙  ♙
	♖  ♘  ♗  ♕  ♔  ♗  ♘  ♖
	`
	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}
func TestMovePawn_Black_cant_move_backwards(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (black pawn cant move backwards)
	// bR  bN  bB  bQ  bK  bB  bN  bR
	// ..  bp  bp  bp  bp  bp  bp  bp
	// bp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  wp  wp  wp  wp  wp  wp  wp
	// wR  wN  wB  wQ  wK  wB  wN  wR

	// black pawn from A7 to A6
	blackMove := Move{Square{Column: "A", Row: 7}, Square{Column: "A", Row: 6}}
	moves := []Move{blackMove}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}

	_, bpA6 := board.GetPieceAtSquare("A", 6)
	_, err := bpA6.Move("A", 7, board, false)
	if err == nil {
		t.Errorf("Expected this move BLACK PAWN A6-A7 to fail, black pawn cant move backwards")
	}
	expectedErr := "pawns can only move forward"
	if strings.TrimSpace(err.Error()) != expectedErr {
		t.Errorf("Expected error message %v, got %v", expectedErr, err.Error())
	}
}

func TestMovePawn_White_cant_move_backwards(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (white pawn cant move backwards)
	// bR  bN  bB  bQ  bK  bB  bN  bR
	// bp  bp  bp  bp  bp  bp  bp  bp
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  wp  wp  wp  wp  wp  wp  wp
	// wR  wN  wB  wQ  wK  wB  wN  wR
	whiteMove := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 3}}
	moves := []Move{whiteMove}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, wpA3 := board.GetPieceAtSquare("A", 3)
	_, err := wpA3.Move("A", 2, board, false)
	if err == nil {
		t.Errorf("Expected this move WHITE PAWN A3-A2 to fail, white pawn cant move backwards")
	}
	expectedErr := "pawns can only move forward"
	if strings.TrimSpace(err.Error()) != expectedErr {
		t.Errorf("Expected error message %v, got %v", expectedErr, err.Error())
	}
}
func TestMovePawn_should_only_be_able_to_move_two_squares_first_move(t *testing.T) {
	defer quiet()()
	board := newBoard()
	_, a2Pawn := board.GetPieceAtSquare("A", 2)
	_, err := a2Pawn.Move("A", 4, board, false)
	if err != nil {
		t.Errorf("Expected pawn to move to A4 successfully on first move.. 2 squares")
	}
	_, a4Pawn := board.GetPieceAtSquare("A", 4)
	_, err2 := a4Pawn.Move("A", 6, board, false)
	expectedErrorMessage := "pawn cant move to square A6, can only move two squares forward on first move and then one square"
	if strings.TrimSpace(err2.Error()) != expectedErrorMessage {
		t.Errorf("Expected error message %v, got %v", expectedErrorMessage, err2.Error())
	}
	if err2 == nil {
		t.Errorf("Expected pawn to NOT move to A6 successfully on second move.. 2 squares")
	}
}
func TestMovePawn_should_be_illegal_if_king_is_in_check_and_it_doesnt_remove_check(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (Black King has to take, no other move is legal)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  bP  bP  bP  WR  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// ..  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	whiteMove2 := Move{Square{Column: "A", Row: 1}, Square{Column: "A", Row: 3}}
	whiteMove3 := Move{Square{Column: "A", Row: 3}, Square{Column: "E", Row: 3}}
	whiteMove4 := Move{Square{Column: "E", Row: 3}, Square{Column: "E", Row: 7}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, whiteMove4}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, bpA7 := board.GetPieceAtSquare("A", 7)
	_, err := bpA7.Move("A", 6, board, false)
	if err == nil {
		t.Errorf("Expected this move BLACK PAWN A7-A6 to fail, black king is in check by white rook at E7")
	}
	expectedErr := "illegal move"
	if strings.TrimSpace(err.Error()) != expectedErr {
		t.Errorf("Expected error message %v, got %v", expectedErr, err.Error())
	}
}

func TestGetValidPawnMoves(t *testing.T) {
	defer quiet()()
	board := newBoard()

	// CREATE START SCENARIO (white pawn takes black pawn)
	// bR  bN  bB  bQ  bK  bB  bN  bR
	// bp  ..  bp  bp  bp  bp  bp  bp
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  bp  ..  ..  ..  ..  ..  ..
	// wp  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wp  wp  wp  wp  wp  wp  wp
	// wR  wN  wB  wQ  wK  wB  wN  wR

	// wp from A2 to A4
	whiteMove1 := Move{Square{Column: "A", Row: 2}, Square{Column: "A", Row: 4}}
	// then bp from B7 to B5
	blackMove1 := Move{Square{Column: "B", Row: 7}, Square{Column: "B", Row: 5}}
	moves := []Move{whiteMove1, blackMove1}
	scenarioPrepError := prepScenario(moves, board)

	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	_, wPA4 := board.GetPieceAtSquare("A", 4)
	peeks := wPA4.getValidPawnMoves(board)
	if len(peeks) != 2 {
		t.Errorf("Expected 2 valid moves for white pawn at A4, got %v", len(peeks))
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
	if numberOfGoTosFound != 1 {
		t.Errorf("Expected 1 goto action for white pawn at A4, got %v", numberOfGoTosFound)
	}
	if len(potentialTakes) != 1 {
		t.Errorf("Expected 1 potential take for white pawn at A4, got %v", len(potentialTakes))
	}
	if potentialTakes[0].GetValue() != 1 {
		t.Errorf("Expected potential take to be worth 1, got %v", potentialTakes[0].GetValue())
	}
}
