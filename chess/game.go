package chess

import (
	"errors"
	"fmt"
)

type Player interface {
	PickMove(g *Game) (*Move, error)
}

type BoardVisualizer interface {
	VisualizeState(b *Board)
}
type Result struct {
	Draw   bool
	Winner Colour
	Reason string
}

type Game struct {
	white              Player
	black              Player
	Board              *Board
	History            []Move
	finished           bool
	NextToMove         Colour
	result             Result
	boardVisualizer    BoardVisualizer
	fiftyRuleCounter   int
	positions          map[string]int
	numberOfWhiteMoves int
	numberOfBlackMoves int
}

// creates and returns a new game
func NewGame(white Player, black Player, stateVisualizer BoardVisualizer) *Game {
	game := &Game{}
	game.white = white
	game.black = black
	game.boardVisualizer = stateVisualizer
	game.finished = false
	game.History = make([]Move, 0)
	game.Board = newBoard()
	game.NextToMove = White
	game.fiftyRuleCounter = 0
	game.positions = make(map[string]int)
	game.numberOfWhiteMoves = 0
	game.numberOfBlackMoves = 0
	return game
}

// starts the game
func (g *Game) Start() Result {
	g.boardVisualizer.VisualizeState(g.Board)
	g.nextMove()
	diff := g.numberOfBlackMoves - g.numberOfWhiteMoves
	if diff < 0 {
		diff = -diff
	}
	if diff > 1 {
		panic(fmt.Sprintf("white number of moves %v, black number of moves %v. expected diff of max 1", g.numberOfWhiteMoves, g.numberOfBlackMoves))
	}
	return g.result
}

func isGameOver(g *Game) bool {
	if g.Board.kingIsInMate(Black) {
		fmt.Println("mate!")
		g.boardVisualizer.VisualizeState(g.Board)
		g.result = Result{Draw: false, Winner: White, Reason: "Black is in mate"}
		return true
	}

	if g.Board.kingIsInMate(White) {
		fmt.Println("mate!")
		g.boardVisualizer.VisualizeState(g.Board)
		g.result = Result{Draw: false, Winner: Black, Reason: "White is in mate"}
		return true
	}

	if g.Board.isStaleMate(White) || g.Board.isStaleMate(Black) {
		fmt.Println("Stale mate!")
		g.boardVisualizer.VisualizeState(g.Board)
		g.result = Result{Draw: true, Winner: White, Reason: "Stale mate"}
		return true
	}

	if g.fiftyRuleCounter >= 50 {
		fmt.Println("50 move rule!")
		g.boardVisualizer.VisualizeState(g.Board)
		g.result = Result{Draw: true, Winner: White, Reason: "50 move rule"}
		return true
	}

	// else if three fold repetition check
	return threefoldRepetitionCheck(g)
}
func threefoldRepetitionCheck(g *Game) bool {
	for k := range g.positions {
		if g.positions[k] >= 3 {
			fmt.Println("3-fold repetition!")
			g.boardVisualizer.VisualizeState(g.Board)
			g.result = Result{Draw: true, Winner: White, Reason: "3-fold repetition"}
			return true
		}
	}
	return false
}
func (g *Game) nextMove() {

	for {

		if g.NextToMove == White { // white to move
			move, pickErr := g.white.PickMove(g)
			if pickErr != nil {
				fmt.Printf("Error picking move: %v", pickErr)
			} else {
				_, err := g.move(*move, White)
				if err != nil {
					fmt.Printf("Error making move: %v\n", err)
				} else {
					g.boardVisualizer.VisualizeState(g.Board)
					fmt.Printf("White moved from %v%v to %v%v", move.From.Column, move.From.Row, move.To.Column, move.To.Row)
					g.NextToMove = Black
				}
			}
		} else { // black to move
			move, pickErr := g.black.PickMove(g)
			if pickErr != nil {
				fmt.Printf("Error picking move: %v", pickErr)
			} else {
				_, err := g.move(*move, Black)
				if err != nil {
					fmt.Printf("Error making move: %v\n", err)
				} else {
					g.boardVisualizer.VisualizeState(g.Board)
					fmt.Printf("Black moved from %v%v to %v%v", move.From.Column, move.From.Row, move.To.Column, move.To.Row)
					g.NextToMove = White
				}
			}
		}
		// after each move, check if game is over
		if isGameOver(g) {
			break
		}
		// else time for the next move (next iteration in game loop)
	}
}
func (g *Game) move(move Move, as Colour) (string, error) {
	_, p := g.Board.GetPieceAtSquare(move.From.Column, move.From.Row)
	if p.Colour != as {
		return "", errors.New("hey! not your piece")
	}
	result, moveErr := p.Move(move.To.Column, move.To.Row, g.Board, false)
	if moveErr != nil {
		return "", moveErr
	}
	if as == White {
		g.numberOfWhiteMoves++
	} else {
		g.numberOfBlackMoves++
	}
	// if no capture has been made and no pawn has been moved in the last fifty moves
	if result.Action == GoTo && p.Type != pawn {
		g.fiftyRuleCounter = g.fiftyRuleCounter + 1
	} else {
		g.fiftyRuleCounter = 0
	}
	if g.fiftyRuleCounter >= 50 {
		return "", fmt.Errorf("FiftyRuleCounter is %v", g.fiftyRuleCounter)
	}
	// keep track of positions, if threefold repetition, draw (does NOT need to be 3 times in a row)
	currentPosition := g.Board.getPosition()
	if _, ok := g.positions[currentPosition]; ok {
		g.positions[currentPosition] = g.positions[currentPosition] + 1
	} else {
		g.positions[currentPosition] = 1
	}
	if threefoldRepetitionCheck(g) {
		return "", fmt.Errorf("ThreefoldRepetitionCheck")
	}
	successMsg := fmt.Sprintf("%v %v moved from %v %v to %v %v", p.Colour, p.Type, move.From.Column,
		move.From.Row, move.To.Column, move.To.Row)
	g.History = append(g.History, Move{From: move.From, To: move.To})
	if g.NextToMove == White {
		g.NextToMove = Black
	} else {
		g.NextToMove = White
	}
	return fmt.Sprintf("%s it is now %vs turn", successMsg, g.NextToMove), nil
}
