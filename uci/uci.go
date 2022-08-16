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

func UCI() {
	command, _ := inputReader.ReadString('\n')
	//command := ""
	//fmt.Scan(&command)
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
	board = chess.ParseFen(startpos)
}

func position(command string) {
	commands := strings.Fields(command)
	if commands[1] == "startpos" {
		uciNewGame()
	} else if commands[1] == "fen" {
		board = chess.ParseFen(strings.Split(command, "position fen ")[1])
	}
	if strings.Contains(command, "moves ") {
		split := strings.Split(command, "moves ")[1]
		if len(split) > 0 {
			moves := strings.Fields(split)
			for i := 0; i < len(moves); i++ {
				move, _ := chess.ParseMove(moves[i])
				_ = board.Apply(move)
			}
		}
	}
}

func goCommand(command string) {
	if !strings.Contains(command, "infinite") {
		clock := -1
		var stopTime int64 = -1
		depth := -1
		movesToGo := 30
		var moveTime int64 = -1
		inc := 0
		goCommand := strings.Fields(command)
		len := len(goCommand)
		if len > 1 {
			if goCommand[1] == "wtime" && board.Wtomove {
				clock, _ = strconv.Atoi(goCommand[2])
			}
		}
		if len > 3 {
			if goCommand[3] == "btime" && !board.Wtomove {
				clock, _ = strconv.Atoi(goCommand[4])
			}
		}
		if len > 5 {
			if goCommand[5] == "winc" && board.Wtomove {
				inc, _ = strconv.Atoi(goCommand[6])
			}
		}
		if len > 7 {
			if goCommand[7] == "binc" && !board.Wtomove {
				inc, _ = strconv.Atoi(goCommand[8])
			}
		}
		if len > 9 {
			if goCommand[9] == "movestogo" {
				movesToGo, _ = strconv.Atoi(goCommand[10])
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
		fmt.Println(moveTime, movesToGo)
		if moveTime != -1 {
			/*timeTotal := clock - 50
			movetime := int(float32(timeTotal)/float32(movesToGo)) + inc
			if inc > 0 && timeTotal < 5*inc {
				movetime = int(float32(75*inc) / 100)
			}*/
			stopTime = start + moveTime - 50 //ajuste
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
