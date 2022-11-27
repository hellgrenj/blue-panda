package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/hellgrenj/blue-panda/chess"
)

type SimpleBot struct {
	Colour    chess.Colour
	DelayInMS int
}

func NewSimpleBot(colour chess.Colour, delayInMS int) *SimpleBot {
	return &SimpleBot{Colour: colour, DelayInMS: delayInMS}
}

func (bot *SimpleBot) PickMove(g *chess.Game) (*chess.Move, error) {
	time.Sleep(time.Duration(bot.DelayInMS) * time.Millisecond) // so we can see what it is doing
	moves := g.Board.GetAllMovesFor(bot.Colour)
	bestMove, err := bot.Evaluate(g, moves)
	if err != nil {
		return nil, err
	}
	return &bestMove, nil
}

type MoveEvaluation struct {
	Move       chess.Move
	MoveResult *chess.MoveResult
	Value      int
	Attacker   *chess.Piece
}

func (bot *SimpleBot) Evaluate(game *chess.Game, moves map[chess.Move]*chess.MoveResult) (chess.Move, error) {
	var evals = make([]MoveEvaluation, 0)
	for m, r := range moves {
		found, p := game.Board.GetPieceAtSquare(m.From.Column, m.From.Row)
		if !found {
			return chess.Move{}, errors.New("Could 	not find piece at square")
		}
		attacker := p
		if r.Action == chess.Take {
			value := r.Piece.GetValue()
			// make temp move & take
			actualCurrentSquare := attacker.CurrentSquare
			attacker.CurrentSquare = chess.Square{Column: m.To.Column, Row: m.To.Row}
			r.Piece.InPlay = false
			// if piece would be lost in the next move, subtract our piece value
			if attacker.Colour == chess.White && chess.CouldAnyTakeAt(game.Board.BlackPieces, m.To.Column, m.To.Row, game.Board) {
				value = value - attacker.GetValue()
			}
			if attacker.Colour == chess.Black && chess.CouldAnyTakeAt(game.Board.WhitePieces, m.To.Column, m.To.Row, game.Board) {
				value = value - attacker.GetValue()
			}
			// undo temp move & take
			attacker.CurrentSquare = actualCurrentSquare
			r.Piece.InPlay = true

			if value == 0 {
				// 80 % chance favour the take
				if rand.Intn(100) < 80 {
					value = 1
				}
			}
			evals = append(evals, MoveEvaluation{Move: m, Attacker: attacker, MoveResult: r, Value: value})
		} else { // GoTos
			evals = append(evals, MoveEvaluation{Move: m, Attacker: attacker, MoveResult: r, Value: 0})
		}
	}

	var bestMove *chess.Move
	highestValue := 0

	for i, e := range evals {
		if e.Value > highestValue {
			highestValue = e.Value
			bestMove = &evals[i].Move
		}
	}
	if highestValue == 0 {
		// if no valuable takes, pick a random move
		bm, err := pick(moves)
		if err != nil {
			return chess.Move{}, err
		}
		bestMove = &bm
	}

	// return if legal and valid or try again
	found, attacker := game.Board.GetPieceAtSquare(bestMove.From.Column, bestMove.From.Row)
	if !found {
		return chess.Move{}, errors.New("Could 	not find piece at square")
	}
	if !attacker.MoveIsLegal(bestMove.To.Column, bestMove.To.Row, game.Board) { // check if legal
		delete(moves, *bestMove)
		newMove, newMoveErr := bot.Evaluate(game, moves)
		if newMoveErr != nil {
			return chess.Move{}, newMoveErr
		}
		return newMove, nil // return the new move we found
	} else {
		_, err := attacker.Move(bestMove.To.Column, bestMove.To.Row, game.Board, true) // check if valid
		if err != nil {
			delete(moves, *bestMove)
			newMove, newMoveErr := bot.Evaluate(game, moves)
			if newMoveErr != nil {
				return chess.Move{}, newMoveErr
			}
			return newMove, nil // return the new move we found
		} else {
			fmt.Printf("\nBot picked move %v%v to %v%v\n", bestMove.From.Column, bestMove.From.Row, bestMove.To.Column, bestMove.To.Row)
			return *bestMove, nil
		}
	}
}

func pick(m map[chess.Move]*chess.MoveResult) (chess.Move, error) {
	maxNumber := len(m)
	if maxNumber <= 0 {
		fmt.Printf("maxNumber is %v len of m is %v", maxNumber, len(m))
		return chess.Move{}, errors.New("no moves available")
	}
	k := rand.Intn(maxNumber)
	i := 0
	for x := range m {
		if i == k {
			return x, nil
		}
		i++
	}
	panic("unreachable")
}
