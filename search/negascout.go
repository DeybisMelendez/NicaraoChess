package search

import (
	"math"
	"nicarao/moveOrdering"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Negascout(board *chess.Board, depth int, alpha int, beta int, turn int, nullMove bool) int {
	if isTimeToStop() {
		return 0
	}
	PVLength[Ply] = Ply
	var bestmove chess.Move
	//var isPVNode bool = beta-alpha > 1
	var hashScore = ReadHashEntry(board.Hash(), alpha, beta, depth, &bestmove)
	if hashScore != NoHashEntry && Ply > 0 {
		return hashScore
	}
	moveList := board.GenerateLegalMoves()
	moveOrdering.SortMoves(moveList, board, PVTable[0][Ply], bestmove, Ply)

	if depth == 0 || len(moveList) == 0 {
		return Quiesce(board, alpha, beta, turn) //eval.Evaluate(board) //
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
	FutilityPruning = false
	var score int = 0
	//check extension
	var incheck bool = board.OurKingInCheck()
	if incheck {
		depth++
	}
	if !incheck {
		//Evaluation pruning
		staticEval := Evaluate(board, turn)
		if depth < 3 && int(math.Abs(float64(beta-1))) > -MateScore+100 {
			evalMargin := MaterialWeightOP[0] * depth
			if staticEval-evalMargin >= beta {
				return staticEval - evalMargin
			}
		}
		if nullMove {
			if Ply > 0 && depth > NullMoveR {
				nullScore := NullMove(board.ToFen(), depth, beta, turn)
				if nullScore != NullMoveFails {
					return beta
				}
			}
			//TODO razoring
			//https://github.com/maksimKorzh/wukongJS/blob/main/wukong.js#L1576
			/*score = staticEval + MaterialWeightOP[0]
			if score < beta && depth == 1 {
				var newScore int = Quiesce(board, alpha, beta, turn)
				if newScore < beta {
					if newScore > score {
						return newScore
					} else {
						return score
					}
				}
			}*/
		}
		CheckFutilityPruning(staticEval, depth, alpha)
	}
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]
		unmakeFunc := Make(board, move)
		//Futility Pruning
		if IsFutilityPruning(board, move, i) {
			Unmake(unmakeFunc)
			continue
		}
		if i > FullDepthMove && isLMROk(board, move) {
			score = -Negascout(board, pvReduction(depth), -beta, -alpha, -turn, DoNull)
		} else {
			score = -Negascout(board, depth-1, -beta, -alpha, -turn, DoNull)
		}
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(move)
			moveOrdering.StoreHistoryMove(move, board, depth)
			bestmove = move
			hashFlag = HashFlagExact
			alpha = score
			if score >= beta {
				moveOrdering.StoreKillerMove(move, board, Ply)
				WriteHashEntry(board.Hash(), beta, depth, HashFlagBeta, move)
				return beta
			}
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

func ResetGlobalVariables() {
	Ply = 0
	Nodes = 0
	Stopped = false
}
