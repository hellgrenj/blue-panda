package chess

import (
	"testing"
)

func TestKingIsInMate(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (Fools mate)
	// BR  BN  BB  ..  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  bP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  wP  BQ
	// ..  ..  ..  ..  ..  wP  ..  ..
	// wP  wP  wP  wP  wP  ..  ..  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "F", Row: 2}, Square{Column: "F", Row: 3}}
	whiteMove2 := Move{Square{Column: "G", Row: 2}, Square{Column: "G", Row: 4}}
	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 5}}
	blackMove2 := Move{Square{Column: "D", Row: 8}, Square{Column: "H", Row: 4}}
	moves := []Move{whiteMove1, whiteMove2, blackMove1, blackMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	// CHECK IF KING IS IN MATE
	kingIsInMate := board.kingIsInMate(White)
	if !kingIsInMate {
		t.Errorf("Expected king to be in mate")
		return
	}
}
func TestKingIsInMate_Cant_escape_by_taking_if_its_moves_into_new_check(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (King is in mate and cat escape by taking the queen)
	// BR  ..  ..  BK  ..  BB  BN  BR
	// bP  bP  bP  WQ  bP  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  ..  wP  wP  wP  wP
	// ..  ..  ..  WR  WK  WB  WN  WR

	// simulate some taken pieces first
	_, b8 := board.GetPieceAtSquare("B", 8)
	b8.InPlay = false
	_, c8 := board.GetPieceAtSquare("C", 8)
	c8.InPlay = false
	_, d8 := board.GetPieceAtSquare("D", 8)
	d8.InPlay = false

	_, b1 := board.GetPieceAtSquare("B", 1)
	b1.InPlay = false
	_, c1 := board.GetPieceAtSquare("C", 1)
	c1.InPlay = false
	_, d2 := board.GetPieceAtSquare("D", 2)
	d2.InPlay = false

	whiteMove1 := Move{Square{Column: "D", Row: 1}, Square{Column: "D", Row: 2}}
	whiteMove2 := Move{Square{Column: "A", Row: 1}, Square{Column: "D", Row: 1}}
	blackMove1 := Move{Square{Column: "E", Row: 8}, Square{Column: "D", Row: 8}}
	whiteMove3 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 7}}
	moves := []Move{whiteMove1, whiteMove2, blackMove1, whiteMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	expectedStateOfBoard := `
	BR  ..  ..  BK  ..  BB  BN  BR 
	bP  bP  bP  WQ  bP  bP  bP  bP 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	wP  wP  wP  ..  wP  wP  wP  wP 
	..  ..  ..  WR  WK  WB  WN  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	// CHECK IF KING IS IN MATE
	kingIsInMate := board.kingIsInMate(Black)
	if !kingIsInMate {
		t.Errorf("Expected king to be in mate")
		return
	}

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
}

func TestKingIsMate_is_false_if_queen_can_block(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (Discovered bug when bot playing bot)
	// ..  ..  ..  BK  ..  ..  BR  ..
	// ..  BB  ..  bP  ..  ..  ..  bP
	// ..  ..  ..  ..  bP  bP  ..  ..
	// ..  bP  bP  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  wP  ..  ..  wP
	// wP  ..  ..  ..  ..  wP  ..  WR
	// WR  ..  ..  ..  ..  WQ  ..  WK

	//Close enough...
	// BR  BN  ..  BQ  BK  BB  BR  ..
	// bP  BB  bP  bP  bp  bP  ..  bP
	// ..  bp  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  wP  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// wP  wP  wP  wP  wP  wP  ..  wP
	// WR  WN  WB  ..  ..  WQ  ..  WK

	// WB F1 taken
	_, WB := board.GetPieceAtSquare("F", 1)
	WB.InPlay = false // simulate taken
	// WN G1 taken
	_, WN := board.GetPieceAtSquare("G", 1)
	WN.InPlay = false // simulate taken
	// WR H1 taken
	_, WR := board.GetPieceAtSquare("H", 1)
	WR.InPlay = false // simulate taken
	// bp g7 taken
	_, bp := board.GetPieceAtSquare("G", 7)
	bp.InPlay = false // simulate taken
	// bn g8 taken
	_, bn := board.GetPieceAtSquare("G", 8)
	bn.InPlay = false // simulate taken
	// wp h2 taken
	_, wp := board.GetPieceAtSquare("G", 2)
	wp.InPlay = false // simulate taken

	// WK
	_, WK := board.GetPieceAtSquare("E", 1)
	WK.CurrentSquare = Square{Column: "H", Row: 1}

	whiteMove1 := Move{Square{Column: "D", Row: 1}, Square{Column: "F", Row: 1}}

	blackMove1 := Move{Square{Column: "B", Row: 7}, Square{Column: "B", Row: 6}}
	blackMove2 := Move{Square{Column: "C", Row: 8}, Square{Column: "B", Row: 7}}
	blackMove3 := Move{Square{Column: "H", Row: 8}, Square{Column: "G", Row: 8}}
	moves := []Move{whiteMove1, blackMove1, blackMove2, blackMove3}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	expectedStateOfBoard := `
	BR  BN  ..  BQ  BK  BB  BR  .. 
	bP  BB  bP  bP  bP  bP  ..  bP 
	..  bP  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	wP  wP  wP  wP  wP  wP  ..  wP 
	WR  WN  WB  ..  ..  WQ  ..  WK 
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
	// CHECK IF KING IS IN MATE
	kingIsInMate := board.kingIsInMate(White)
	if kingIsInMate {
		t.Errorf("Expected king to NOT be in mate, queen on f1 can block by moving to g2")
		return
	}
}
func TestKingIsInMate_returns_false_if_king_can_run(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (NOT mate but check)
	// BR  BN  BB  ..  BK  BB  BN  BR
	// bP  bP  bP  bP  ..  bP  bP  bP
	// ..  ..  ..  ..  bP  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  wP  BQ
	// ..  ..  ..  ..  wP  wP  ..  ..
	// wP  wP  wP  wP  ..  ..  ..  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR
	whiteMove1 := Move{Square{Column: "F", Row: 2}, Square{Column: "F", Row: 3}}
	whiteMove2 := Move{Square{Column: "G", Row: 2}, Square{Column: "G", Row: 4}}
	whiteMove3 := Move{Square{Column: "E", Row: 2}, Square{Column: "E", Row: 3}}
	blackMove1 := Move{Square{Column: "E", Row: 7}, Square{Column: "E", Row: 5}}
	blackMove2 := Move{Square{Column: "D", Row: 8}, Square{Column: "H", Row: 4}}
	moves := []Move{whiteMove1, whiteMove2, whiteMove3, blackMove1, blackMove2}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}

	// CHECK IF KING IS IN MATE
	kingIsInMate := board.kingIsInMate(White)
	if kingIsInMate {
		t.Errorf("King should not be in mate, can escape to E2")
		return
	}
}

func TestKingIsInMate_returns_false_transient_bug(t *testing.T) {

	board := newBoard()
	// CREATE START SCENARIO (transient bug scenario)
	// ..  ..  ..  ..  BK  BB  ..  ..
	// ..  ..  bP  ..  ..  ..  bP  ..
	// bP  bP  ..  bP  ..  ..  ..  BR
	// ..  ..  ..  ..  bP  BB  ..  bP
	// ..  ..  ..  ..  ..  BQ  ..  ..
	// ..  wP  wP  WK  ..  ..  ..  ..
	// wP  ..  ..  wP  wP  ..  wP  wP
	// WR  WN  ..  ..  ..  WB  ..  WR

	_, BR := board.GetPieceAtSquare("A", 8)
	BR.InPlay = false // simulate taken
	_, aBp := board.GetPieceAtSquare("A", 7)
	aBp.CurrentSquare = Square{Column: "A", Row: 6}
	_, bWp := board.GetPieceAtSquare("B", 2)
	bWp.CurrentSquare = Square{Column: "B", Row: 3}
	_, bBN := board.GetPieceAtSquare("B", 8)
	bBN.InPlay = false // simulate taken
	_, bbP := board.GetPieceAtSquare("B", 7)
	bbP.CurrentSquare = Square{Column: "B", Row: 6}
	_, cWB := board.GetPieceAtSquare("C", 1)
	cWB.InPlay = false // simulate taken
	_, cwP := board.GetPieceAtSquare("C", 2)
	cwP.CurrentSquare = Square{Column: "C", Row: 3}
	_, cBB := board.GetPieceAtSquare("C", 8)
	cBB.CurrentSquare = Square{Column: "F", Row: 5}
	_, WQ := board.GetPieceAtSquare("D", 1)
	WQ.InPlay = false // simulate taken
	_, WK := board.GetPieceAtSquare("E", 1)
	WK.CurrentSquare = Square{Column: "D", Row: 3}
	_, BQ := board.GetPieceAtSquare("D", 8)
	BQ.CurrentSquare = Square{Column: "F", Row: 4}
	_, dbP := board.GetPieceAtSquare("D", 7)
	dbP.CurrentSquare = Square{Column: "D", Row: 6}
	_, eBp := board.GetPieceAtSquare("E", 7)
	eBp.CurrentSquare = Square{Column: "E", Row: 5}
	_, fwP := board.GetPieceAtSquare("F", 2)
	fwP.InPlay = false // simulate taken
	_, fbP := board.GetPieceAtSquare("F", 7)
	fbP.InPlay = false // simulate taken
	_, gWN := board.GetPieceAtSquare("G", 1)
	gWN.InPlay = false // simulate taken
	_, gBN := board.GetPieceAtSquare("G", 8)
	gBN.InPlay = false // simulate taken
	_, hbP := board.GetPieceAtSquare("H", 7)
	hbP.CurrentSquare = Square{Column: "H", Row: 5}
	_, hBR := board.GetPieceAtSquare("H", 8)
	hBR.CurrentSquare = Square{Column: "H", Row: 6}

	expectedStateOfBoard := `
	..  ..  ..  ..  BK  BB  ..  ..
	..  ..  bP  ..  ..  ..  bP  ..
	bP  bP  ..  bP  ..  ..  ..  BR
	..  ..  ..  ..  bP  BB  ..  bP
	..  ..  ..  ..  ..  BQ  ..  ..
	..  wP  wP  WK  ..  ..  ..  ..
	wP  ..  ..  wP  wP  ..  wP  wP
	WR  WN  ..  ..  ..  WB  ..  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	// CHECK IF KING IS IN MATE
	kingIsInMate := board.kingIsInMate(White)
	if kingIsInMate {
		t.Errorf("Expected king to NOT be in mate, E2 Pawn can move to E4 and block the black bishop")
		return
	}

}
func TestKingIsInMate_returns_false_if_king_can_run_BLACK(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (NOT mate but check - vs Bot game 1 scenario.. bot gave up with check mate.. but its not)
	// "Black King is in check by White Bishop at H 5" ... but can run to D7
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// ..  ..  bP  ..  bP  ..  bP  bP
	// bP  bP  ..  bP  ..  bP  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  WB
	// ..  ..  ..  wP  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  wP  ..
	// wP  wP  wP  ..  wP  wP  ..  wP
	// WR  WN  WB  WQ  WK  ..  WN  WR

	whiteMove1 := Move{Square{Column: "D", Row: 2}, Square{Column: "D", Row: 4}}
	blackMove1 := Move{Square{Column: "D", Row: 7}, Square{Column: "D", Row: 6}}
	whiteMove2 := Move{Square{Column: "G", Row: 2}, Square{Column: "G", Row: 3}}
	blackMove2 := Move{Square{Column: "A", Row: 7}, Square{Column: "A", Row: 6}}
	whiteMove3 := Move{Square{Column: "F", Row: 1}, Square{Column: "H", Row: 3}}
	blackMove3 := Move{Square{Column: "F", Row: 7}, Square{Column: "F", Row: 6}}
	whiteMove4 := Move{Square{Column: "H", Row: 3}, Square{Column: "G", Row: 4}}
	blackMove4 := Move{Square{Column: "B", Row: 7}, Square{Column: "B", Row: 6}}
	whiteMove5 := Move{Square{Column: "G", Row: 4}, Square{Column: "H", Row: 5}}

	moves := []Move{whiteMove1, blackMove1, whiteMove2, blackMove2, whiteMove3, blackMove3, whiteMove4, blackMove4, whiteMove5}
	scenarioPrepError := prepScenario(moves, board)
	if scenarioPrepError != nil {
		t.Errorf("Failed to prep the board, %v", scenarioPrepError.Error())
		return
	}
	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  BN  BR 
	..  ..  bP  ..  bP  ..  bP  bP 
	bP  bP  ..  bP  ..  bP  ..  .. 
	..  ..  ..  ..  ..  ..  ..  WB 
	..  ..  ..  wP  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  wP  .. 
	wP  wP  wP  ..  wP  wP  ..  wP 
	WR  WN  WB  WQ  WK  ..  WN  WR
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}
	// CHECK IF KING IS IN MATE
	kingIsInMate := board.kingIsInMate(Black)
	if kingIsInMate {
		t.Errorf("King should not be in mate, can escape to D7")
		return
	}
}

func TestKingIsInCheck_returns_false_if_not_in_check(t *testing.T) {
	board := newBoard()
	// CREATE START SCENARIO (NOT check)
	// BR  BN  BB  BQ  BK  BB  BN  BR
	// bP  ..  bP  ..  bp  bP  bP  bP
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  bP  ..  ..  ..  ..
	// wp  bP  ..  ..  ..  ..  ..  ..
	// ..  ..  ..  ..  ..  ..  ..  ..
	// ..  wP  wP  wP  wP  wP  wP  wP
	// WR  WN  WB  WQ  WK  WB  WN  WR

	_, wp := board.GetPieceAtSquare("A", 2)
	wp.CurrentSquare = Square{Column: "A", Row: 4}

	_, bp := board.GetPieceAtSquare("B", 7)
	bp.CurrentSquare = Square{Column: "B", Row: 4}

	_, bp2 := board.GetPieceAtSquare("D", 7)
	bp2.CurrentSquare = Square{Column: "D", Row: 5}

	expectedStateOfBoard := `
	BR  BN  BB  BQ  BK  BB  BN  BR 
	bP  ..  bP  ..  bP  bP  bP  bP 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  ..  ..  bP  ..  ..  ..  .. 
	wP  bP  ..  ..  ..  ..  ..  .. 
	..  ..  ..  ..  ..  ..  ..  .. 
	..  wP  wP  wP  wP  wP  wP  wP 
	WR  WN  WB  WQ  WK  WB  WN  WR 
	`

	if err := assertExpectedBoardState(expectedStateOfBoard, board); err != nil {
		t.Errorf("Failed to assert expected board state, %v (Visible whitespace is ignored, something else differs!", err.Error())
	}

	if kinInCheck, _ := board.kingIsInCheck(Black); kinInCheck {
		t.Errorf("King should not be in check!, pawn at A4 cant reach it...	")
		return
	}
}
