package search

import (
	"math"
	"nicarao/moveOrdering"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Negamax(board *chess.Board, depth int, alpha int, beta int, turn int, nullMove bool) int {
	PVLength[Ply] = Ply
	var hashFlag int = HashFlagAlpha
	var score int = 0
	var bestmove chess.Move
	var isPVNode bool = beta-alpha > 1
	if IsRepetition(board.Hash()) {
		return 0
	}
	var hashScore int = ReadHashEntry(board.Hash(), alpha, beta, depth, &bestmove)
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
	var matingValue int = MateValue - Ply
	if matingValue < beta {
		beta = matingValue
		if alpha >= matingValue {
			return matingValue
		}
	}
	matingValue = -MateValue + Ply
	if matingValue > alpha {
		alpha = matingValue
		if beta <= matingValue {
			return matingValue
		}
	}

	var inCheck bool = board.OurKingInCheck()
	var staticEval int = Evaluate(board, turn)
	if nullMove && !isPVNode && !inCheck {
		//https://www.chessprogramming.org/Null_Move_Pruning#Schemes
		//if Ply > 0 && depth > NullDepth && AllowNullMove(board) && !isEndgame(board) {
		if Ply > 0 && depth > NullDepth && staticEval >= beta && !isEndgame(board) {
			nullScore := NullMove(board, depth, alpha, beta, turn)
			if nullScore != NullMoveFails {
				return beta
			}
			if isTimeToStop() {
				return 0
			}
		}
	}
	moveList := board.GenerateLegalMoves()
	len := len(moveList)
	if len > 1 {
		moveOrdering.SortMoves(moveList, board, PVTable[0][Ply], bestmove, Ply)
	}
	var movesSearched int = 0
	sqrtDepth := math.Sqrt(float64(depth - 1))
	sqrtCount := math.Sqrt(float64(len - 1))
	R := int(math.Sqrt(sqrtDepth + sqrtCount))
	for _, move := range moveList {
		var isCapture bool = chess.IsCapture(move, board)
		var givesCheck bool = board.OurKingInCheck()
		var isPromotion bool = move.Promote() != chess.Nothing
		var isKillerMove bool = moveOrdering.IsKillerMove(move, Ply)
		var nonTactical bool = !isCapture && !inCheck && !givesCheck && !isPromotion && !isKillerMove
		if nonTactical {
			// Razoring segun stockfish
			/*if depth == 2 && staticEval < alpha+1800 {
				var newScore int = Quiesce(board, alpha, beta, turn)
				if newScore < beta {
					return utils.Max(newScore, score)
				}
			}*/
			//Futility pruning
			var value int = staticEval + 60
			if depth == 1 && int(math.Abs(float64(alpha))) < MateScore && value <= alpha {
				continue
			}
		}
		var newDepth int = depth * 2 / 3
		if !nonTactical {
			newDepth++
		}
		unmakeFunc := Make(board, move)
		if movesSearched == 0 {
			score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if movesSearched > 4 && nonTactical && depth > 2 {
				score = -Negamax(board, R, -alpha-1, -alpha, -turn, DoNull)
			} else {
				score = alpha + 1
			}
			if score > alpha {
				score = -Negamax(board, newDepth, -alpha-1, -alpha, -turn, DoNull)
				if score > alpha && score < beta {
					score = -Negamax(board, newDepth, -beta, -alpha, -turn, DoNull)
				}
			}
		}
		if isTimeToStop() {
			return 0
		}
		Unmake(unmakeFunc)
		movesSearched++
		if score > alpha {
			StorePV(move)
			moveOrdering.StoreHistoryMove(move, board, depth)
			bestmove = move
			hashFlag = HashFlagExact
			alpha = score
			// Reduce other moves
			/*if depth > 2 {
				depth--
			}*/
			if score >= beta {
				moveOrdering.StoreKillerMove(move, board, Ply)
				WriteHashEntry(board.Hash(), beta, depth, HashFlagBeta, move)
				return beta
			}
		}
	}

	if len == 0 {
		if board.OurKingInCheck() {
			//Checkmate
			return -MateValue + Ply
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
	Stopped = false
	var newRep = [150]uint64{}
	RepetitionTable = newRep
}
