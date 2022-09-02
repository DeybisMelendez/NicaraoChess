package engine

import (
	"fmt"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

const Infinity = 10000
const MateValue = 9900
const MateScore = 9800
const DoNull = true
const NoNull = false
const NullDepth = 2

var GamePly int = 0
var Ply int = 0
var Nodes int = 0
var StopTime int64 = -1
var Stopped bool = false

func Search(board *chess.Board, stopTime int64, depth int) {
	start := time.Now().UnixMilli()
	var bestmove chess.Move
	var lastBestmove chess.Move
	FollowPV = false
	ScorePV = false
	score := 0
	Nodes = 0
	scoreType := "cp"
	alpha := -Infinity
	beta := Infinity
	StopTime = stopTime
	currDepth := 1
	turn := -1
	if board.Wtomove {
		turn = 1
	}
	for { //Iterative Deepening
		// TODO detener en jaque mate
		lastBestmove = PVTable[0][0]
		if depth == 0 {
			break
		}
		FollowPV = true
		score = Negamax(board, currDepth, alpha, beta, turn, DoNull)
		if isTimeToStop() {
			break
		}
		if score >= MateScore {
			score = (MateValue - score + 1) / 2
			scoreType = "mate"
		} else if score <= -MateScore {
			score = -(MateValue + score) / 2
			scoreType = "mate"
		}
		fmt.Println("info",
			"depth", currDepth,
			"score", scoreType, score,
			"nodes", Nodes,
			"time", time.Now().UnixMilli()-start,
			"pv", FormatPV(PVTable[0]))
		/*if scoreType == "mate" {
			break
		}*/
		ResetGlobalVariables()
		depth--
		currDepth++
	}
	if Stopped {
		bestmove = lastBestmove
	} else {
		bestmove = PVTable[0][0]
	}
	toPrint := "bestmove " + bestmove.String()
	fmt.Println(toPrint)
	ClearSearch()
}

func ClearSearch() {
	InitHasTable()
	ResetPVTable()
	ResetGlobalVariables()
	ResetKillerMoves()
	ResetHistoryMove()
	var newRep = [1000]uint64{}
	RepetitionTable = newRep
	GamePly = 0
}
