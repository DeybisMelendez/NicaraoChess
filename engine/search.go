package engine

import (
	"fmt"
	"math"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

const VAL_WINDOW = 10
const MateValue = math.MaxInt16
const INF = math.MaxInt16

var Nodes uint64
var Ply int16
var GamePly int16
var StopTime int64
var Stopped bool = false

func clearVariables() {
	transpositionTable = make([]TranspositionTable, ttSize)
	repetitionTable = [256]uint64{}
	Nodes = 0
	Ply = 0
	killerMoves = [2][64]chess.Move{}
	historyMoves = [2][64][64]uint16{}

}
func isTimeToStop() bool {
	if Stopped {
		return Stopped
	}
	if StopTime != -1 {
		if Nodes&2047 == 0 {
			if time.Now().UnixMilli() >= StopTime {
				Stopped = true
			}
		}
	}
	return Stopped
}

func Search(board *chess.Board, depth int8, timeLimit int64) {
	var start int64 = time.Now().UnixMilli()
	var alpha int16 = -INF
	var beta int16 = INF
	var turn int16 = -1
	var score int16
	var currDepth int8 = 1
	var bestmove chess.Move
	Nodes = 0
	Ply = 0
	Stopped = false
	if timeLimit > 0 {
		StopTime = start + timeLimit
	} else {
		StopTime = -1
	}
	if board.Wtomove {
		turn = 1
	}

	for { //Iterative Deepening
		if depth > 0 && currDepth > depth {
			break
		}

		score = negamax(board, currDepth, turn, alpha, beta)
		if timeLimit > 0 && isTimeToStop() {
			break
		}
		// Aspiration Window
		/*if (score <= alpha) || (score >= beta) {
			alpha = -INF
			beta = INF
			continue
		} else {
			alpha = score - VAL_WINDOW
			beta = score + VAL_WINDOW
		}*/

		// Print the best move for the current depth
		var pv string = getPV(board.ToFen(), currDepth, &bestmove)
		var time_elapsed int64 = time.Now().UnixMilli() - start

		fmt.Printf("info depth %d score cp %d nodes %d time %d pv %s\n",
			currDepth, score, Nodes, time_elapsed, pv)

		currDepth++
	}
	fmt.Println("bestmove " + bestmove.String())
	clearVariables()
}
