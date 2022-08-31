package search

import (
	"nicarao/moveOrdering"
	"strings"
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
	if nullMove && !inCheck {
		if Ply > 0 && depth > NullDepth && !isEndgame(board) {
			var staticEval int = Evaluate(board, turn)
			if staticEval >= beta {
				var fen = board.ToFen()
				if strings.Contains(fen, " w ") {
					fen = strings.ReplaceAll(fen, " w ", " b ")
				} else {
					fen = strings.ReplaceAll(fen, " b ", " w ")
				}
				nullBoard := chess.ParseFen(fen)
				if !nullBoard.OurKingInCheck() && len(nullBoard.GenerateLegalMoves()) != 0 {
					eval := -Negamax(&nullBoard, depth-NullDepth-1, -beta, -beta+1, -turn, NoNull)
					if eval >= beta {
						return eval
					}
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
		var isPromotion bool = move.Promote() != chess.Nothing
		unmakeFunc := Make(board, move)
		var givesCheck bool = board.OurKingInCheck()
		var isTactical bool = inCheck || givesCheck || isCapture || isPromotion
		if !isTactical && depth < 3 && depth > 0 {
			var staticEval int = Evaluate(board, turn)
			// Razoring
			if depth == 2 && staticEval+50 < alpha {
				depth--
			} else if depth == 1 && staticEval+50 < alpha {
				Unmake(unmakeFunc)
				continue
			}
		}
		if movesSearched == 0 {
			score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
		} else {
			if movesSearched > 3 && !isTactical && depth > 2 {
				score = -Negamax(board, depth-2, -alpha-1, -alpha, -turn, DoNull)
			} else {
				score = alpha + 1
			}
			if score > alpha {
				score = -Negamax(board, depth-1, -alpha-1, -alpha, -turn, DoNull)
				if score > alpha && score < beta {
					score = -Negamax(board, depth-1, -beta, -alpha, -turn, DoNull)
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
			bestmove = move
			hashFlag = HashFlagExact
			alpha = score
			if score >= beta {
				moveOrdering.StoreKillerMove(move, board, Ply)
				moveOrdering.StoreHistoryMove(move, board, depth)
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
