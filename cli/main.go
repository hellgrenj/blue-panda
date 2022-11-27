package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hellgrenj/blue-panda/chess"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	Menu()
}
func Menu() {

	fmt.Println("What do you want to play?")
	fmt.Println("1. Human vs Human")
	fmt.Println("2. Human vs Computer")
	fmt.Println("3. Computer vs Computer")
	fmt.Println("4. 100 games of Computer vs Computer")
	reader := bufio.NewReader(os.Stdin)
	gameType, _ := reader.ReadString('\n')
	switch gameType {
	case "1\n":
		whitePlayer := &Player{Colour: chess.White}
		blackPlayer := &Player{Colour: chess.Black}
		startGame(whitePlayer, blackPlayer)
	case "2\n":
		whitePlayer := &Player{Colour: chess.White}
		blackPlayer := NewSimpleBot(chess.Black, 1500)
		startGame(whitePlayer, blackPlayer)
	case "3\n":
		whitePlayer := NewSimpleBot(chess.White, 200)
		blackPlayer := NewSimpleBot(chess.Black, 200)
		startGame(whitePlayer, blackPlayer)
	case "4\n":
		whitePlayer := NewSimpleBot(chess.White, 0)
		blackPlayer := NewSimpleBot(chess.Black, 0)
		results := make(map[chess.Result]int)
		for i := 0; i < 100; i++ {
			result := startGame(whitePlayer, blackPlayer)
			results[result]++
		}
		fmt.Println("\nResults:")
		for k, v := range results {
			if k.Draw {
				fmt.Printf("Draw (%v): %d times\n", k.Reason, v)
			} else {
				fmt.Printf("%v Wins (%v): %d times\n", k.Winner, k.Reason, v)
			}
		}
	default:
		fmt.Println("Invalid option. You can enter 1, 2, 3 or 4. please try again")
		Menu()
	}
}
func startGame(whitePlayer chess.Player, blackPlayer chess.Player) chess.Result {
	game := chess.NewGame(whitePlayer, blackPlayer, &CLIPrinter{})
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGSEGV)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("\ngame aborted by signal %v, printing history of moves:\n", sig)
		printHistory(game)
		os.Exit(0)
	}()
	result := game.Start()
	if result.Draw {
		fmt.Printf("Draw, reason: %v\n", result.Reason)
	} else {
		fmt.Printf("%v wins! (%v)\n\n", result.Winner, result.Reason)
	}
	printHistory(game)
	return result
}
func printHistory(g *chess.Game) {
	f, err := os.Create("history")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	for _, move := range g.History {
		f.WriteString(fmt.Sprintf("%v", move))
		f.WriteString("\n")
	}
}

type CLIPrinter struct {
}

func (v *CLIPrinter) VisualizeState(b *chess.Board) {
	// clear terminal screen
	fmt.Print("\033[H\033[2J")
	Print(b)
}

type Player struct {
	Colour chess.Colour
}

func (p *Player) PickMove(g *chess.Game) (*chess.Move, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n%v to move: ", g.NextToMove)
	move, _ := reader.ReadString('\n')

	// input validation
	match, _ := regexp.MatchString("[A-Ha-h][1-8] [A-Ha-h][1-8]", move)
	if !match {
		return nil, errors.New("invalid input, please enter a move like this: E2 E4 or e2 e4 (single space between squares)")
	}

	// split move into from and to
	moveParts := strings.Split(move, " ")
	from := moveParts[0]
	to := strings.Trim(moveParts[1], "\n")
	fromColumn := strings.ToUpper(string(from[0]))
	fromRow, fromSquareErr := strconv.Atoi(from[1:])

	if fromSquareErr != nil {
		return nil, fmt.Errorf("invalid from square")
	}
	toColumn := strings.ToUpper(string(to[0]))
	toRow, toSquareErr := strconv.Atoi(to[1:])
	if toSquareErr != nil {
		fmt.Println(toSquareErr)
		return nil, fmt.Errorf("invalid to square")
	}

	return &chess.Move{
		From: chess.Square{Column: fromColumn, Row: fromRow},
		To:   chess.Square{Column: toColumn, Row: toRow},
	}, nil
}
func Print(b *chess.Board) {

	for i, square := range b.Squares {
		if i%8 == 0 {
			fmt.Printf("\n")
		}
		occupied, piece := b.GetPieceAtSquare(square.Column, square.Row)
		if occupied {
			fmt.Printf(" %v ", piece.GetAbbreveation())
		} else {
			fmt.Printf(" .. ")
		}
	}
	fmt.Printf("\n")
	fmt.Printf("White's captured pieces:")
	for _, piece := range b.BlackPieces {
		if !piece.InPlay {
			fmt.Printf(" %v ", piece.GetAbbreveation())
		}
	}
	fmt.Println("")
	fmt.Printf("Black's captured pieces:")
	for _, piece := range b.WhitePieces {
		if !piece.InPlay {
			fmt.Printf(" %v ", piece.GetAbbreveation())
		}
	}
	fmt.Println("")

}
