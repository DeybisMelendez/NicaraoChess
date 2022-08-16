package search

import (
	"nicarao/moveOrdering"
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

func Evaluate(board *chess.Board, color int) int {
	eval := 0
	//GetMaterial(board)
	eval += ValueMaterial(board)
	return eval * color
}

func Quiesce(board *chess.Board, color int, alpha int, beta int) int {
	standPat := Evaluate(board, color)
	if standPat > beta {
		return beta
	}
	alpha = utils.Max(alpha, standPat)
	moves := filterCaptures(board.GenerateLegalMoves(), board, color)
	var score int = 0
	for i := 0; i < len(moves); i++ {
		unmakeFunc := Make(board, moves[i])
		score = -Quiesce(board, -color, -beta, -alpha)
		Unmake(unmakeFunc)
		if isTimeToStop() {
			return 0
		}
		if score >= beta {
			return beta
		}
		if score > alpha {
			StorePV(moves[i])
			alpha = score
		}
		//alpha = utils.Max(alpha, score)
	}
	return alpha
}

func filterCaptures(moves []chess.Move, board *chess.Board, color int) []chess.Move {
	var filteredCaptures []chess.Move
	for i := 0; i < len(moves); i++ {
		if chess.IsCapture(moves[i], board) {
			filteredCaptures = append(filteredCaptures, moves[i])
		}
	}
	filteredCaptures = moveOrdering.SortMoves(filteredCaptures, board, PVTable[0], Bestmove, color, Ply)
	return filteredCaptures
}
