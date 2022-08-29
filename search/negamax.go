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
	var applyFutility bool
	if !isPVNode {
		// Razoring segun stockfish
		if depth == 2 && staticEval < alpha-900 {
			var newScore int = Quiesce(board, alpha, beta, turn)
			if newScore < beta {
				return utils.Max(newScore, score)
			}
		}
		//Futility pruning
		var value int = staticEval + 100*depth*depth
		if depth < 4 && int(math.Abs(float64(alpha))) < MateScore && value <= alpha {
			applyFutility = true
		}
		if nullMove && !inCheck {
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
	}
	moveList := board.GenerateLegalMoves()
	len := len(moveList)
	if len > 1 {
		moveOrdering.SortMoves(moveList, board, PVTable[0][Ply], bestmove, Ply)
	}
	var movesSearched int = 0
	for _, move := range moveList {
		var isCapture bool = chess.IsCapture(move, board)
		var givesCheck bool = board.OurKingInCheck()
		var isPromotion bool = move.Promote() != chess.Nothing
		var isKillerMove bool = moveOrdering.IsKillerMove(move, Ply)
		from := move.From()
		piece, isWhite := utils.GetPiece(from, board)
		var isHistoryMove bool = moveOrdering.GetHistoryMove(isWhite, piece, from) > 50*depth
		var tacticalMove bool = isCapture || inCheck || givesCheck || isPromotion || isKillerMove || isHistoryMove
		if applyFutility && movesSearched > 0 && !tacticalMove {
			continue
		}
		var newDepth int = depth
		unmakeFunc := Make(board, move)
		//reductions
		if !isPVNode && !tacticalMove && newDepth > 5 {
			if beta < MateValue-Ply && alpha > -MateValue+Ply {
				newDepth -= 3
			}
		}
		if movesSearched == 0 {
			score = -Negamax(board, newDepth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if movesSearched > 4 && !tacticalMove {
				score = -Negamax(board, depth*2/3, -alpha-1, -alpha, -turn, DoNull)
			} else {
				score = alpha + 1
			}
			if score > alpha {
				score = -Negamax(board, newDepth-1, -alpha-1, -alpha, -turn, DoNull)
				if score > alpha && score < beta {
					score = -Negamax(board, newDepth-1, -beta, -alpha, -turn, DoNull)
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
			if depth > 1 && beta < MateValue-Ply && alpha > -MateValue+Ply {
				depth--
			}
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
