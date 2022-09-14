package engine

import (
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Negamax(board *chess.Board, depth int, alpha int, beta int, turn int, nullMove bool) int {
	PVLength[Ply] = Ply
	var hashFlag int = HashFlagAlpha
	var hashmove chess.Move
	var isPVNode bool = beta-alpha > 1
	var score int = ReadHashEntry(board.Hash(), alpha, beta, depth, &hashmove)
	if score != NoHashEntry && !isPVNode && Ply > 0 {
		return score
	}
	if isTimeToStop() {
		return 0
	}
	if IsRepetition(board.Hash()) {
		return 0
	}
	if depth == 0 || Ply >= 64 {
		return Quiesce(board, alpha, beta, turn)
	}
	// Mate Distance pruning
	if alpha < -MateValue {
		alpha = -MateValue
	}
	if beta > MateValue-1 {
		beta = MateValue - 1
	}
	if alpha >= beta {
		return alpha
	}
	var inCheck bool = board.OurKingInCheck()
	var moveList []chess.Move = board.GenerateLegalMoves()
	var bestmove chess.Move = checkPV(moveList)
	var movesSearched int = 0
	var lenMoveList int = len(moveList)
	// Check Extension
	/*if inCheck {
		depth++
	}*/
	var scoreMoveList []int
	var isCapture bool
	var isPromotion bool
	for _, v := range moveList {
		isCapture = chess.IsCapture(v, board)
		isPromotion = v.Promote() != chess.Nothing
		value := ValueMove(board, v, isCapture, isPromotion, bestmove, hashmove)
		scoreMoveList = append(scoreMoveList, value)
	}
	for len(moveList) > 0 {
		var val int = -1
		var idx int = 0
		//var ln int = len(moveList)
		var isCapture bool
		var isPromotion bool
		for i, newVal := range scoreMoveList {
			//var newVal int = scoreMoveList[i]
			if newVal > val {
				val = newVal
				idx = i
			}
		}
		var move = moveList[idx]
		isCapture = chess.IsCapture(move, board)
		isPromotion = move.Promote() != chess.Nothing
		var isTactical bool = inCheck || isCapture || isPromotion
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		scoreMoveList = append(scoreMoveList[:idx], scoreMoveList[idx+1:]...)
		unmakeFunc := Make(board, move)

		if !isTactical && !isPVNode {
			// Razoring
			if depth == 2 {
				var staticEval = Evaluate(board, turn) + 50
				if staticEval < alpha {
					Unmake(unmakeFunc)
					break
				}
				//Futility Pruning
			} else if depth == 1 {
				if !board.OurKingInCheck() {
					var staticEval = Evaluate(board, turn) + 50
					if staticEval < alpha {
						Unmake(unmakeFunc)
						continue
					}
				}
			}
		}
		//var newDepth = depth
		/*if newDepth > 3 {
			if isPVNode && score == NoHashEntry {
				newDepth -= 2
			}
		}*/
		/*if movesSearched == 0 {
			score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
		} else {*/
		if movesSearched > 2 { // && !isTactical && !IsKillerMove(move) {
			score = -Negamax(board, depth*2/3, -alpha-1, -alpha, -turn, DoNull)
		} else {
			score = alpha + 1
		}
		if score > alpha {
			score = -Negamax(board, depth-1, -alpha-1, -alpha, -turn, DoNull)
			if score > alpha && score < beta {
				score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
			}
		}
		//}
		Unmake(unmakeFunc)
		movesSearched++
		if score > alpha {
			StorePV(move)
			hashFlag = HashFlagExact
			alpha = score
			bestmove = move
			if depth > 6 && depth < 12 {
				depth--
			}
			if score >= beta {
				if !isCapture {
					StoreKillerMove(move)
					StoreHistoryMove(move, board.Wtomove, depth)
				}
				WriteHashEntry(board.Hash(), beta, depth, HashFlagBeta, move)
				return beta
			}
		}

	}
	if lenMoveList == 0 {
		if inCheck {
			//Checkmate
			return -MateValue + Ply
		} else {
			//Stalemate
			return 0
		}
	}
	WriteHashEntry(board.Hash(), alpha, depth, hashFlag, bestmove)
	return alpha
}

func isTimeToStop() bool {
	if Stopped {
		return true
	}
	if StopTime != -1 && Nodes&8191 == 0 {
		if time.Now().UnixMilli() >= StopTime {
			Stopped = true
			return true
		}
	}
	return false
}
