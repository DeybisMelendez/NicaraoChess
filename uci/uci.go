package uci

import (
	"bufio"
	"fmt"
	"nicarao/engine"
	"os"
	"strconv"
	"strings"

	chess "github.com/dylhunn/dragontoothmg"
)

const START_POS string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func UCI() {
	var board chess.Board
	engine.ResizeTranspositionTable(512)
	engine.SetDist()
	engine.InitPieceSquareTable()
	engine.InitEvaluationMask()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			panic("Error reading input")
		}
		text = strings.TrimSpace(text)
		if text == "uci" {
			fmt.Println("id name NicaChess v0.4.0")
			fmt.Println("id author DeybisMelendez")
			fmt.Println("uciok")
		} else if text == "isready" {
			fmt.Println("readyok")
		} else if text == "ucinewgame" {
			// code to handle a new game starting
		} else if strings.HasPrefix(text, "position") {
			// code to handle setting the board position
			setPosition(text, &board)
		} else if strings.HasPrefix(text, "go") {
			//var start int64 = time.Now().UnixMilli()
			// code to handle starting the search
			fields := strings.Fields(text)
			// Set the time controls
			var wtime, btime, winc, binc, movestogo, movetime, depth int
			//var depth int
			var timeLeft int
			movestogo = 30
			for i := 0; i < len(fields); i++ {
				switch fields[i] {
				case "wtime":
					wtime, _ = strconv.Atoi(fields[i+1])
				case "btime":
					btime, _ = strconv.Atoi(fields[i+1])
				case "winc":
					winc, _ = strconv.Atoi(fields[i+1])
				case "binc":
					binc, _ = strconv.Atoi(fields[i+1])
				case "movestogo":
					movestogo, _ = strconv.Atoi(fields[i+1])
				case "movetime":
					movetime, _ = strconv.Atoi(fields[i+1])
				case "depth":
					depth, _ = strconv.Atoi(fields[i+1])
				}
			}
			if wtime != 0 || btime != 0 || movetime != 0 {
				// search with time controls
				var inc int
				if board.Wtomove {
					timeLeft = wtime
					inc = winc
				} else {
					timeLeft = btime
					inc = binc
				}
				timeLeft = (timeLeft - 50) / movestogo
				if inc > 0 {
					if timeLeft < 10000 {
						timeLeft += inc / 2
					} else {
						timeLeft += inc
					}
				}
				//go engine.SearchTime(&board, start+int64(timeLeft))
			} /* else {
				go engine.SearchDepth(&board, int8(depth))
			}*/
			go engine.Search(&board, int8(depth), int64(timeLeft))
			// code to start the search with the given time controls
			// and/or search depth
		} else if text == "quit" {
			// code to quit the program
			os.Exit(0)
		}
	}
}

func setPosition(text string, board *chess.Board) {
	fields := strings.Fields(text)
	if fields[1] == "startpos" {
		// Set up the starting position
		*board = chess.ParseFen(START_POS)
	} else if fields[1] == "fen" {
		// Set up the position specified by FEN string
		*board = chess.ParseFen(strings.Split(text, "position fen ")[1])
	}
	// Process any additional moves specified in the command
	if strings.Contains(text, "moves") {
		// Make the move on the board
		moves := strings.Split(text, "moves ")[1]
		fields := strings.Fields(moves)
		for _, m := range fields {
			move, err := chess.ParseMove(m)
			if err != nil {
				fmt.Println(err)
			}
			_ = board.Apply(move)
			engine.GamePly++
			engine.StoreRepetitionTable(board.Hash())
		}
	}
}
