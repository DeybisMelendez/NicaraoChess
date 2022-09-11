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
	if depth == 0 {
		return Quiesce(board, alpha, beta, turn) //Evaluate(board,turn) //
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
	/*R := 2
	if nullMove && !inCheck && depth > R {
		if !isEndgame(board) {
			//var staticEval = Evaluate(board, turn)
			//if staticEval >= beta {
			//Null Move Reduction
			board.Wtomove = !board.Wtomove
			halfmoveclock := board.Halfmoveclock
			board.Halfmoveclock = 0
			nullBoard := chess.ParseFen(board.ToFen())
			board.Wtomove = !board.Wtomove
			board.Halfmoveclock = halfmoveclock
			if len(nullBoard.GenerateLegalMoves()) != 0 {
				//if depth-R-1 > 0 {
				score = -Negamax(&nullBoard, depth-R-1, -beta, -beta+1, -turn, NoNull)
				//} else {
				//	score = Quiesce(board, -beta, -beta+1, -turn)
				//}
				//eval := -ZWSearch(&nullBoard, depth-NullDepth-1, -beta, -turn, NoNull)
				if score >= beta {
					return beta
				}
			}
			//}
		}
	}*/
	var moveList []chess.Move = board.GenerateLegalMoves()
	var bestmove chess.Move = checkPV(moveList)
	if hashmove == 0 && bestmove == 0 && isPVNode && nullMove && depth >= 6 && Ply > 1 {
		Negamax(board, depth-2, alpha, beta, turn, DoNull)
		/*if score > alpha {
			if score < beta {

			}
		}*/
	}
	var movesSearched int = 0
	var lenMoveList int = len(moveList)
	/*var staticEval = Evaluate(board, turn)
	// Razoring
	if depth == 2 && staticEval+50 < alpha {
		depth--
	}*/
	// One Reply & Check Extension
	if lenMoveList == 1 || inCheck {
		depth++
	}
	for len(moveList) > 0 {
		var val int = -1
		var idx int = 0
		var ln int = len(moveList)
		var isCapture bool
		var isPromotion bool
		for i := 0; i < ln; i++ {
			isCapture = chess.IsCapture(moveList[i], board)
			isPromotion = moveList[i].Promote() != chess.Nothing
			var newVal int = ValueMove(board, moveList[i], isCapture, isPromotion, bestmove, hashmove)
			if newVal > val {
				val = newVal
				idx = i
			}
		}
		var isTactical bool = inCheck || isCapture || isPromotion
		var move = moveList[idx]
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		unmakeFunc := Make(board, move)
		//Futility Pruning
		if !isTactical && depth == 1 && !isPVNode {
			if !board.OurKingInCheck() {
				var staticEval int = Evaluate(board, turn)
				if staticEval+50 < alpha {
					Unmake(unmakeFunc)
					continue
				}

			}
		}
		var newDepth = depth
		if newDepth > 3 {
			if isPVNode && score == NoHashEntry {
				newDepth -= 2
			}
		}
		if movesSearched == 0 {
			score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if movesSearched > 2 && !isTactical && !IsKillerMove(move) { // && GetHistoryMove(board.Wtomove, move) > depth*depth*10 {
				//score = -ZWSearch(board, depth-2, -alpha, -turn, nullMove)
				score = -Negamax(board, newDepth*2/3, -alpha-1, -alpha, -turn, DoNull)
			} else {
				score = alpha + 1
			}
			if score > alpha {
				//score = -ZWSearch(board, depth-1, -alpha, -turn, nullMove)
				score = -Negamax(board, newDepth-1, -alpha-1, -alpha, -turn, DoNull)
				if score > alpha && score < beta {
					score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
				}
			}
		}
		Unmake(unmakeFunc)
		movesSearched++
		if score > alpha {
			StorePV(move)
			hashFlag = HashFlagExact
			alpha = score
			bestmove = move
			if depth > 4 && depth < 9 {
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
	if StopTime != -1 && Nodes&1023 == 0 {
		if time.Now().UnixMilli() >= StopTime {
			Stopped = true
			return true
		}
	}
	return false
}
