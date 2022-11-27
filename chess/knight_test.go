package chess

import (
	"testing"
)

func TestMoveKnight(t *testing.T) {
	defer quiet()()
	board := newBoard()

	_, WN := board.GetPieceAtSquare("G", 1)
	_, err := WN.Move("F", 3, board, false)
	if err != nil {
		t.Errorf("Failed to move the knight from G1 to F3, %v", err.Error())
	}
	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  BN  BR
	bP  bP  bP  bP  bP  bP  bP  bP
	..  ..  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  ..  ..  ..  ..
	..  ..  ..  ..  ..  WN  ..  ..
	wP  wP  wP  wP  wP  wP  wP  wP
	WR  WN  WB  WQ  WK  WB  ..  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}
func TestGetValidKnightMoves(t *testing.T) {
	defer quiet()()
	board := newBoard()
	// From starting position
	_, WN := board.GetPieceAtSquare("G", 1)
	peeks := WN.getValidKnightMoves(board)
	if len(peeks) != 2 {
		t.Errorf("Expected to find 2 valid moves for the knight, but found %v", len(peeks))
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "G", Row: 1}, To: Square{Column: "F", Row: 3}}) {
		t.Errorf("Expected one of the valid moves to be G1-F3, but it wasn't")
	}
	if !mapContainsKey(peeks, Move{From: Square{Column: "G", Row: 1}, To: Square{Column: "H", Row: 3}}) {
		t.Errorf("Expected one of the valid moves to be G1-H3, but it wasn't")
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
		t.Errorf("Expected 0 take action for white knight at G1, got %v", numberOfTakesFound)
	}
	if numberOfGoTosFound != 2 {
		t.Errorf("Expected 2 goto action for white knight at G1, got %v", numberOfGoTosFound)
	}
	if len(potentialTakes) != 0 {
		t.Errorf("Expected 0 potential take for white knight at G1, got %v", len(potentialTakes))
	}
}
func TestMoveKnight_should_be_able_to_remove_check(t *testing.T) {
	defer quiet()()
	board := newBoard()
	// SCENARIO WN should be table to take BR to remove check on WK
	// BR  BN  BB  BQ  BK  BB  ..  ..
	// bP  bP  bP  bP  bP  bP  ..  bP
	// ..  ..  ..  ..  ..  BN  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  WN
	// wP  wP  wP  wP  wP  wP  ..  wP
	// WR  WN  WB  WQ  WK  ..  BR  WR
	// also simulate WB on F1 being taken AND wp on G2 being taken AND bp on g7 being taken

	_, bPG7 := board.GetPieceAtSquare("G", 7)
	bPG7.InPlay = false // simulate taken

	_, WBF1 := board.GetPieceAtSquare("F", 1)
	WBF1.InPlay = false // simulate taken

	_, wpG2 := board.GetPieceAtSquare("G", 2)
	wpG2.InPlay = false // simulate taken

	whiteMove1 := Move{Square{Column: "G", Row: 1}, Square{Column: "H", Row: 3}} // WN to H3
	blackMove1 := Move{Square{Column: "G", Row: 8}, Square{Column: "F", Row: 6}} // BN to F6
	blackMove2 := Move{Square{Column: "H", Row: 8}, Square{Column: "G", Row: 8}} // BR to G8
	blackMove3 := Move{Square{Column: "G", Row: 8}, Square{Column: "G", Row: 1}} // BR to G1 check

	moves := []Move{whiteMove1, blackMove1, blackMove2, blackMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}

	kingInCheck, _ := board.kingIsInCheck(White)
	if !kingInCheck {
		t.Errorf("Expected the king to be in check, but it wasn't")
		return
	}

	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  ..  .. 
	bP  bP  bP  bP  bP  bP  ..  bP 
	..  ..  ..  ..  ..  BN  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  WN 
	wP  wP  wP  wP  wP  wP  ..  wP 
	WR  WN  WB  WQ  WK  ..  BR  WR
	`
	// and assert king not in check...
	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	_, WN := board.GetPieceAtSquare("H", 3)
	_, err := WN.Move("G", 1, board, false)
	if err != nil {
		t.Errorf("Failed to move the knight from H3 to G1 to remove check, %v", err.Error())
	}
	expectedStateOfBoard = `
	BR  BN  BB  BQ  BK  BB  ..  .. 
	bP  bP  bP  bP  bP  bP  ..  bP 
	..  ..  ..  ..  ..  BN  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	wP  wP  wP  wP  wP  wP  ..  wP 
	WR  WN  WB  WQ  WK  ..  WN  WR

	`
	// and assert king not in check...
	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	// finally assert king not in check...
	kingInCheck, _ = board.kingIsInCheck(White)
	if kingInCheck {
		t.Errorf("Expected the king to NOT be in check anymore since WN to BR on G1")
		return
	}
}
