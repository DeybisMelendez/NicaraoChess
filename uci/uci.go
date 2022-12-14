package uci

import (
	"bufio"
	"fmt"
	"nicarao/engine"
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
	engine.PST = engine.PstMake()
	engine.InitHasTable()
	engine.InitEvaluationMask()
	engine.SetDist()
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
		engine.SetHashTable(uint64(mb))

	}
	if on {
		UCI()
	}
}

func uci() {
	fmt.Println("id name NicaraoChess v0.3.1")
	fmt.Println("id author Deybis Melendez")
	fmt.Println("option name Hash type spin default 16 min 4 max 16")
	fmt.Println("uciok")
}

func isReady() {
	fmt.Println("readyok")
}

func uciNewGame() {
	engine.ResetRepetitionTable()
	engine.ClearSearch()
}

func position(command string) {
	//uciNewGame()
	commands := strings.Fields(command)
	if commands[1] == "startpos" {
		board = chess.ParseFen(startpos)
	} else if commands[1] == "fen" {
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
				engine.GamePly++
				engine.RepetitionTable[engine.GamePly] = board.Hash()
			}
		}
	}
}

func goCommand(command string) {
	start := time.Now().UnixMilli()
	if !strings.Contains(command, "infinite") {
		var clock int64 = -1
		var stopTime int64 = -1
		depth := -1
		var movesToGo int64 = 30
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
		if moveTime != -1 {
			//TODO:Implementar un mejor control de tiempo
			stopTime = start + moveTime
		}
		if clock != -1 {
			timeleft := (clock - 50) / movesToGo
			if inc > 0 {
				if clock < 10000 {
					timeleft += inc * 2 / 3
				} else {
					timeleft += inc
				}
			}
			stopTime = start + timeleft
		}
		if depth == -1 {
			depth = 32
		}
		/*fmt.Println(
		"Time:", clock,
		"Inc:", inc,
		"Start:", start,
		"Stop:", stopTime,
		"Depth:", depth)*/
		go engine.Search(&board, stopTime, start, depth)
	}
}
