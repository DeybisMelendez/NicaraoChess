package search

import (
	"nicarao/moveOrdering"
	"nicarao/utils"
	"time"

	chess "github.com/dylhunn/dragontoothmg"
)

func Negascout(board *chess.Board, depth int, alpha int, beta int, turn int, nullMove bool) int {
	if isTimeToStop() {
		return 0
	}
	PVLength[Ply] = Ply
	var bestmove chess.Move
	var hashScore = ReadHashEntry(board.Hash(), alpha, beta, depth, &bestmove)
	if hashScore != NoHashEntry && Ply > 0 {
		return hashScore
	}
	moveList := board.GenerateLegalMoves()
	//list := board.GenerateLegalMoves()
	//var moveList []chess.Move = moveOrdering.SortMoves(list, board, PVTable[0][Ply], bestmove, Ply)
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
	var score int = 0
	if nullMove {
		if !board.OurKingInCheck() && Ply > 0 && depth > NullMoveR {
			nullScore := NullMove(board.ToFen(), depth, beta, turn)
			if nullScore != NullMoveFails {
				return beta
			}
		}
	}
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]
		unmakeFunc := Make(board, move)
		if i > FullDepthMove && isLMROk(board, move) {
			score = -Negascout(board, pvReduction(depth), -alpha-1, -alpha, -turn, NoNull)
			if score > alpha {
				////https://www.chessprogramming.org/NegaScout#Guido_Schimmels
				var score2 int = -Negascout(board, depth-1, -beta, -alpha, -turn, DoNull)
				score = utils.Max(score, score2)
			}
		} else {
			score = -Negascout(board, depth-1, -beta, -alpha, -turn, DoNull)
		}
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(move)
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
