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
	if nullMove && !inCheck && depth > 2 {
		//Null Move Reduction
		if Ply > 0 { // && !isEndgame(board) {
			var staticEval int = Evaluate(board, turn)
			if staticEval >= beta {
				board.Wtomove = !board.Wtomove
				nullBoard := chess.ParseFen(board.ToFen())
				board.Wtomove = !board.Wtomove
				if len(nullBoard.GenerateLegalMoves()) != 0 {
					var R int = 2
					if depth > 6 {
						R = 3
					}
					eval := -Negamax(&nullBoard, depth-R-1, -beta, -beta+1, -turn, NoNull)
					//eval := -ZWSearch(&nullBoard, depth-NullDepth-1, -beta, -turn, NoNull)
					if eval >= beta {
						depth -= R
						//return eval
					}
					if depth <= 0 {
						return Quiesce(board, alpha, beta, turn)
					}
				}
			}
		}
	}
	moveList := board.GenerateLegalMoves()
	checkPV(moveList)
	var movesSearched int = 0
	var lenMoveList int = len(moveList)
	// One Reply Extension
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
			var newVal int = ValueMove(board, moveList[i], isCapture, isPromotion, PVTable[0][Ply], hashmove)
			if newVal > val {
				val = newVal
				idx = i
			}
		}
		var isTactical bool = inCheck || isCapture || isPromotion
		var move = moveList[idx]
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		unmakeFunc := Make(board, move)
		if !isTactical && depth < 3 && depth > 0 {
			var staticEval int = Evaluate(board, turn)
			// Razoring
			if depth == 2 && staticEval+50 < alpha {
				depth--
				//Futility Pruning
			} else if depth == 1 && staticEval+50 < alpha {
				Unmake(unmakeFunc)
				continue
			}
		}
		if movesSearched == 0 {
			score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if movesSearched > 2 && !isTactical && depth > 2 {
				//score = -ZWSearch(board, depth-2, -alpha, -turn, nullMove)
				score = -Negamax(board, depth-2, -alpha-1, -alpha, -turn, DoNull)
			} else {
				score = alpha + 1
			}
			if score > alpha {
				//score = -ZWSearch(board, depth-1, -alpha, -turn, nullMove)
				score = -Negamax(board, depth-1, -alpha-1, -alpha, -turn, DoNull)
				if score > alpha && score < beta {
					score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
				}
			}
		}
		Unmake(unmakeFunc)
		if isTimeToStop() {
			return 0
		}
		movesSearched++
		if score > alpha {
			StorePV(move)
			hashFlag = HashFlagExact
			alpha = score
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

	WriteHashEntry(board.Hash(), beta, depth, hashFlag, hashmove)
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

func ResetGlobalVariables() {
	Ply = 0
	Stopped = false
	/*var newRep = [150]uint64{}
	RepetitionTable = newRep*/
}

func checkPV(moveList []chess.Move) {
	if FollowPV {
		for _, move := range moveList {
			FollowPV = false
			if move == PVTable[0][Ply] {
				FollowPV = true
				break
			}
		}
	}
}
