package search

import (
	"fmt"
	"nicarao/moveOrdering"
	"nicarao/utils"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Search(board *chess.Board, stopTime int64, depth int) {
	score := 0
	InitHasTable()
	ResetPVTable()
	ResetGlobalVariables()
	moveOrdering.ResetKillerMoves()
	moveOrdering.ResetHistoryMove()
	//eval.ResetMaterial()
	//eval.GetMaterial(board)
	const infinity = 10000
	alpha := -infinity
	beta := infinity
	StopTime = stopTime
	currDepth := 1
	/*color := -1
	//var bestmove chess.Move
	if board.Wtomove {
		color = 1
	}*/
	for {
		// TODO detener en jaque mate
		if depth == 0 {
			break
		}
		//Mate = MateScore
		moveOrdering.FollowPV = true
		//eval.GetMaterial(board)
		score = Negascout(board, currDepth, alpha, beta)
		if isTimeToStop() {
			break
		}
		Bestmove = PVTable[0][0]
		fmt.Println("info",
			"depth", currDepth,
			"score cp", score,
			"nodes", Nodes,
			"pv", FormatPV(PVTable[0]))
		ResetGlobalVariables()
		//eval.ResetMaterial()
		//eval.GetMaterial(board)
		depth--
		currDepth++
	}
	//Bestmove = PVTable[0][0]
	toPrint := "bestmove " + Bestmove.String()
	fmt.Println(toPrint)
}

func Negascout(board *chess.Board, depth int, alpha int, beta int) int {
	if isTimeToStop() {
		return 0
	}
	PVLength[Ply] = Ply
	var hashScore, bestmove = ReadHashEntry(board.Hash(), alpha, beta, depth)
	if hashScore != NoHashEntry && Ply > 0 {
		//StorePV(bestmove)
		return hashScore
	}
	list := board.GenerateLegalMoves()
	var moveList []chess.Move = moveOrdering.SortMoves(list, board, PVTable[0], bestmove, Ply)
	//var moveList []chess.Move = board.GenerateLegalMoves()
	if depth == 0 || len(moveList) == 0 {
		return Quiesce(board, alpha, beta) //eval.Evaluate(board) //
	}
	// Mate Distance pruning
	if alpha < -MateScore {
		alpha = -MateScore
	}
	if beta > MateScore-1 {
		beta = MateScore - 1
	}
	if alpha >= beta {
		return alpha
	}
	var hashFlag int = HashFlagAlpha
	var score int = 0
	/*if nullMove {
		nullScore := NullMove(*board, board.OurKingInCheck(), depth, alpha, beta)
		if nullScore != NullMoveFails {
			return beta
		}
	}*/
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]
		unmakeFunc := Make(board, move)
		if i > FullDepthMove && isLMROk(board, move) {
			score = -Negascout(board, pvReduction(depth), -alpha-1, -alpha)
			if score > alpha {
				////https://www.chessprogramming.org/NegaScout#Guido_Schimmels
				var score2 int = -Negascout(board, depth-1, -beta, -alpha)
				score = utils.Max(score, score2)
			}
		} else {
			score = -Negascout(board, depth-1, -beta, -alpha)
		}
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(move)
			bestmove = move
			moveOrdering.StoreHistoryMove(move, board, depth)
			hashFlag = HashFlagExact
			alpha = score
		}
		if score >= beta {
			moveOrdering.StoreKillerMove(move, board, Ply)
			WriteHashEntry(board.Hash(), beta, depth, HashFlagBeta, move)
			return beta
		}

	}
	WriteHashEntry(board.Hash(), beta, depth, hashFlag, bestmove)
	return alpha
}

func isTimeToStop() bool {
	if Stopped {
		return true
	}
	if StopTime != -1 {
		if time.Now().UnixMilli() >= StopTime {
			Stopped = true
			return true
		}
	}
	return false
}
