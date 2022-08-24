package search

import (
	"fmt"
	"nicarao/moveOrdering"

	chess "github.com/dylhunn/dragontoothmg"
)

const Infinity = 10000
const MateScore = 8000

var Ply int = 0
var Nodes int = 0
var StopTime int64 = -1
var Stopped bool = false
var Bestmove chess.Move

func Search(board *chess.Board, stopTime int64, depth int) {
	score := 0
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
		if depth == 0 {
			break
		}
		moveOrdering.FollowPV = true
		score = Negamax(board, currDepth, alpha, beta, turn, DoNull)
		if isTimeToStop() {
			break
		}
		Bestmove = PVTable[0][0]
		if score >= MateScore-500 {
			score = (MateScore - score + 1) / 2
			scoreType = "mate"
		} else if score <= -MateScore+500 {
			score = (MateScore + score) / 2
			scoreType = "mate"
		}
		fmt.Println("info",
			"depth", currDepth,
			"score", scoreType, score,
			"nodes", Nodes,
			"pv", FormatPV(PVTable[0]))
		/*if scoreType == "mate" {
			break
		}*/
		ResetGlobalVariables()
		depth--
		currDepth++
	}
	toPrint := "bestmove " + Bestmove.String()
	fmt.Println(toPrint)
	ClearSearch()
}

func ClearSearch() {
	InitHasTable()
	ResetPVTable()
	ResetGlobalVariables()
	moveOrdering.ResetKillerMoves()
	moveOrdering.ResetHistoryMove()
}
