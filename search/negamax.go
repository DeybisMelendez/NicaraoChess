package search

import (
	"math"
	"nicarao/moveOrdering"
	"nicarao/utils"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Negamax(board *chess.Board, depth int, alpha int, beta int, turn int, nullMove bool) int {
	PVLength[Ply] = Ply
	var hashFlag int = HashFlagAlpha
	var score int = 0
	var bestmove chess.Move
	var isPVNode bool = beta-alpha > 1
	if IsThreeFoldRepetition(board.Hash()) {
		return 0
	}
	var hashScore = ReadHashEntry(board.Hash(), alpha, beta, depth, &bestmove)
	if hashScore != NoHashEntry && Ply > 0 && !isPVNode {
		return hashScore
	}
	if isTimeToStop() {
		return 0
	}
	if depth == 0 {
		return Quiesce(board, alpha, beta, turn) //Evaluate(board,turn) //
	}
	// Mate Distance pruning
	/*if alpha < -MateScore {
		alpha = -MateScore
	}
	if beta > MateScore-1 {
		beta = MateScore - 1
	}
	if alpha >= beta {
		return alpha
	}*/
	//check extension
	var inCheck bool = board.OurKingInCheck()
	if inCheck {
		depth++
	}
	if !inCheck && !isPVNode {
		var staticEval int = Evaluate(board, turn)
		//evaluation pruning
		if depth < 3 && int(math.Abs(float64(beta-1))) > -MateScore+100 {
			evalMargin := 100 * depth
			if staticEval-evalMargin >= beta {
				return staticEval - evalMargin
			}
		}
		if nullMove {
			//https://www.chessprogramming.org/Null_Move_Pruning#Schemes
			if Ply > 0 && depth > NullMoveR && AllowNullMove(board) && !isEndgame(board) {
				//if Ply > 0 && depth > NullMoveR && staticEval > beta && !isEndgame(board) {
				nullScore := NullMove(board.ToFen(), depth, beta, turn)
				if nullScore != NullMoveFails {
					return beta
				}
				if isTimeToStop() {
					return 0
				}
			}
			//Razoring
			score = staticEval + 100
			if score < beta && depth == 1 {
				var newScore int = Quiesce(board, alpha, beta, turn)
				if newScore < beta {
					return utils.Max(newScore, score)
				}
			}
		}
	}
	moveList := board.GenerateLegalMoves()
	moveOrdering.SortMoves(moveList, board, PVTable[0][Ply], bestmove, Ply)
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]
		var isCapture bool = chess.IsCapture(move, board)
		unmakeFunc := Make(board, move)
		/*if depth == 1 && !isPVNode {
			var staticEval int = Evaluate(board, turn)
			var isPromotion bool = move.Promote() != chess.Nothing
			if IsFutilityPruning(staticEval, alpha, board, inCheck, isCapture, isPromotion) {
				Unmake(unmakeFunc)
				continue
			}
		}*/
		if i < FullDepthMove {
			score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if isLMROk(board, inCheck, isCapture, move) && !isPVNode {
				score = -Negamax(board, pvReduction(depth), -alpha-1, -alpha, -turn, NoNull)
			} else {
				score = alpha + 1
			}
			if score > alpha {
				score = -Negamax(board, depth-1, -alpha-1, -alpha, -turn, NoNull)
				if score > alpha && score < beta {
					score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
				}
			}
		}
		if isTimeToStop() {
			return 0
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

	if len(moveList) == 0 {
		if board.OurKingInCheck() {
			//Checkmate
			return -MateScore + Ply
		} else {
			//Stalemate
			return 0
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
