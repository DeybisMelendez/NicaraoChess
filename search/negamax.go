package search

import (
	"nicarao/moveOrdering"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Negamax(board *chess.Board, depth int, alpha int, beta int, turn int, nullMove bool) int {
	if isTimeToStop() {
		return 0
	}
	PVLength[Ply] = Ply
	moveList := board.GenerateLegalMoves()
	var bestmove chess.Move
	//var isPVNode bool = beta-alpha > 1
	var hashScore = ReadHashEntry(board.Hash(), alpha, beta, depth, &bestmove)
	if hashScore != NoHashEntry && Ply > 0 {
		return hashScore
	}
	moveOrdering.SortMoves(moveList, board, PVTable[0][Ply], bestmove, Ply)

	if depth == 0 || len(moveList) == 0 {
		return Quiesce(board, alpha, beta, turn) //Evaluate(board,turn) //
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
	} else {
		//Evaluation pruning
		staticEval := Evaluate(board, turn)
		if nullMove {
			if Ply > 0 && depth > NullMoveR {
				nullScore := NullMove(board.ToFen(), depth, beta, turn)
				if nullScore != NullMoveFails {
					return beta
				}
			}
		}
		CheckFutilityPruning(staticEval, depth, alpha)
	}
	bSearchPv := true
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]
		unmakeFunc := Make(board, move)
		//Futility Pruning
		if IsFutilityPruning(board, move, i) {
			Unmake(unmakeFunc)
			continue
		}
		if i >= FullDepthMove && depth > 2 {
			score = -Negamax(board, pvReduction(depth), -alpha-1, -alpha, -turn, NoNull)
		} else {
			score = alpha + 1
		}
		if score > alpha {
			if bSearchPv {
				score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
			} else {
				score = -Negamax(board, depth-1, -alpha-1, -alpha, -turn, NoNull)
				if score > alpha {
					score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
				}
			}
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
				bSearchPv = false
				return beta
			}
		}
	}

	if len(moveList) == 0 {
		if board.OurKingInCheck() {
			//Checkmate
			return -MateScore + Ply
		} else {
			//Stalemate
			return 0
		}
	}
	if IsThreeFoldRepetition(board.Hash()) {
		return 0
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
