package chess

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func quiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}

func assertExpectedBoardState(expectedState string, board *Board) error {
	state := sPrintStateOfBoard(board)
	if spaceFieldsJoin(expectedState) == spaceFieldsJoin(state) {
		return nil
	} else {
		return fmt.Errorf("expected board state to be %v, but got %v", expectedState, state)
	}
}
func sPrintStateOfBoard(b *Board) string {
	var state strings.Builder
	for i, square := range b.Squares {
		if i%8 == 0 {
			state.WriteString("\n")
		}
		occupied, piece := b.GetPieceAtSquare(square.Column, square.Row)
		if occupied {
			state.WriteString(fmt.Sprintf(" %v ", piece.GetAbbreveation()))
		} else {
			state.WriteString(" . ")
		}
	}
	state.WriteString("\n")
	return state.String()
}

func spaceFieldsJoin(str string) string {
	return strings.Join(strings.Fields(str), "")
}

func mapContainsKey(m map[Move]*MoveResult, key Move) bool {
	_, ok := m[key]
	return ok
}
