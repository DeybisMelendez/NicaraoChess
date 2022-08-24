package uci

import (
	"bufio"
	"fmt"
	"nicarao/search"
	"os"
	"strconv"
	"strings"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

const startpos string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var board chess.Board
var on bool = true
var inputReader *bufio.Reader = bufio.NewReader(os.Stdin)

func Init() {
	search.PST = search.PstMake()
	search.InitHasTable()
	search.InitEvaluationMask()
}
func UCI() {
	command, _ := inputReader.ReadString('\n')
	command = strings.ReplaceAll(command, "\n", "")
	if command == "uci" {
		uci()
	} else if command == "isready" {
		isReady()
	} else if command == "quit" {
		on = false
	} else if command == "ucinewgame" {
		uciNewGame()
	} else if strings.Contains(command, "position") {
		position(command)
	} else if strings.Contains(command, "go") {
		goCommand(command)
	} else if strings.Contains(command, "setoption name Hash value") {
		slice := strings.Fields(command)
		mb, _ := strconv.Atoi(slice[len(slice)-1])
		search.SetHashTable(uint64(mb))

	}
	if on {
		UCI()
	}
}

func uci() {
	fmt.Println("id name NicaraoChess")
	fmt.Println("id author Deybis Melendez")
	fmt.Println("option name Hash type spin default 16 min 4 max 128")
	fmt.Println("uciok")
}

func isReady() {
	fmt.Println("readyok")
}

func uciNewGame() {
	//search.SetHashTable(search.Mb)
	search.ClearSearch()
	board = chess.ParseFen(startpos)
}

func position(command string) {
	commands := strings.Fields(command)
	if commands[1] == "startpos" {
		uciNewGame()
	} else if commands[1] == "fen" {
		search.ClearSearch()
		board = chess.ParseFen(strings.Split(command, "position fen ")[1])
	}
	if strings.Contains(command, "moves ") {
		split := strings.Split(command, "moves ")[1]
		if len(split) > 0 {
			moves := strings.Fields(split)
			for i := 0; i < len(moves); i++ {
				move, err := chess.ParseMove(moves[i])
				if err != nil {
					fmt.Println(err)
				}
				_ = board.Apply(move)
			}
		}
	}
}

func goCommand(command string) {
	if !strings.Contains(command, "infinite") {
		var clock int64 = -1
		var stopTime int64 = -1
		depth := -1
		var movesToGo int64 = 50
		var moveTime int64 = -1
		var inc int64 = 0
		goCommand := strings.Fields(command)
		len := len(goCommand)
		if len > 1 {
			if goCommand[1] == "wtime" && board.Wtomove {
				x, _ := strconv.Atoi(goCommand[2])
				clock = int64(x)
			}
		}
		if len > 3 {
			if goCommand[3] == "btime" && !board.Wtomove {
				x, _ := strconv.Atoi(goCommand[4])
				clock = int64(x)
			}
		}
		if len > 5 {
			if goCommand[5] == "winc" && board.Wtomove {
				x, _ := strconv.Atoi(goCommand[6])
				inc = int64(x)
			}
		}
		if len > 7 {
			if goCommand[7] == "binc" && !board.Wtomove {
				x, _ := strconv.Atoi(goCommand[8])
				inc = int64(x)
			}
		}
		if len > 9 {
			if goCommand[9] == "movestogo" {
				x, _ := strconv.Atoi(goCommand[10])
				movesToGo = int64(x)
			}
		}
		if goCommand[1] == "movetime" {
			x, _ := strconv.Atoi(goCommand[2])
			moveTime = int64(x)
		}
		if goCommand[1] == "depth" {
			depth, _ = strconv.Atoi(goCommand[2])
		}
		start := time.Now().UnixMilli()
		if moveTime != -1 {
			//TODO:Implementar un mejor control de tiempo
			stopTime = start + moveTime
		}
		if clock != -1 {
			timeleft := (clock-50)/movesToGo + inc
			stopTime = start + timeleft
		}
		if depth == -1 {
			depth = 64
		}
		fmt.Println(
			"Time:", clock,
			"Inc:", inc,
			"Start:", start,
			"Stop:", stopTime,
			"Depth:", depth)
		go search.Search(&board, stopTime, depth)
	}
}
