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
	alpha = utils.Max(alpha, -MateScore+Ply-1)
	beta = utils.Min(beta, MateScore-Ply)
	if alpha >= beta {
		return alpha
	}
	var inCheck bool = board.OurKingInCheck()
	var staticEval int = Evaluate(board, turn)
	var applyFutility bool
	if !isPVNode {
		// Razoring segun stockfish
		if depth < 8 && staticEval < alpha-300-200*depth*depth {
			var newScore int = Quiesce(board, alpha, beta, turn)
			if newScore < beta {
				return utils.Max(newScore, score)
			}
		}
		//Futility pruning
		var value int = staticEval + 100*depth*depth
		if depth < 8 && int(math.Abs(float64(alpha))) < MateScore-500 && value <= alpha {
			applyFutility = true
		}
		//if beta < MateScore-Ply && alpha > -MateScore+Ply {
		/*if IsFutilityPruning(staticEval, depth, alpha, board) && !inCheck {
			return staticEval
		}*/
		//}

	}
	if !inCheck && !isPVNode {
		if nullMove {
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
	moveOrdering.SortMoves(moveList, board, PVTable[0][Ply], bestmove, Ply)
	var movesSearched int = 0
	for _, move := range moveList {
		var isCapture bool = chess.IsCapture(move, board)
		var givesCheck bool = board.OurKingInCheck()
		var isPromotion bool = move.Promote() != chess.Nothing
		if applyFutility && movesSearched > 0 && !isCapture && !isPromotion && !givesCheck && !inCheck {
			continue
		}
		var newDepth int = depth
		//extensions
		if depth > 9 {
			if givesCheck || inCheck || moveOrdering.IsKillerMove(move, Ply) {
				newDepth++
			} else {
				square := move.From()
				piece, isWhite := utils.GetPiece(square, board)
				if moveOrdering.GetHistoryMove(isWhite, piece, square) > 100 {
					newDepth++
				}
			}
		}
		//reductions
		if !isPVNode && bestmove != move && !inCheck && !givesCheck && !isCapture && move.Promote() == chess.Nothing && !moveOrdering.IsKillerMove(move, Ply) {
			newDepth -= 3
			if newDepth <= 0 {
				return Quiesce(board, alpha, beta, turn)
			}
		}
		unmakeFunc := Make(board, move)
		if movesSearched == 0 {
			score = -Negamax(board, newDepth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if movesSearched > FullDepthMove && isLMROk(board, inCheck, isCapture, move) && !isPVNode {
				score = -Negamax(board, pvReduction(newDepth), -alpha-1, -alpha, -turn, DoNull)
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
			if depth > 2 && depth < 7 && beta < MateScore-Ply && alpha > -MateScore+Ply {
				depth--
			}
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
	Stopped = false
	var newRep = [150]uint64{}
	RepetitionTable = newRep
}
