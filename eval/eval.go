package eval

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

func Evaluate(board *chess.Board, color int) int {
	eval := 0
	eval += MaterialMG
	return eval * color
}

func Quiesce(board *chess.Board, color int, alpha int, beta int, ply *int, nodes *int) int {
	standPat := Evaluate(board, color)
	if standPat > beta {
		return beta
	}
	alpha = utils.Max(alpha, standPat)
	moves := filterCaptures(board.GenerateLegalMoves(), board)
	var score int = 0
	for i := 0; i < len(moves); i++ {
		unmakeFunc := utils.Make(board, moves[i], ply, nodes)
		score = -Quiesce(board, -color, -beta, -alpha, ply, nodes)
		utils.Unmake(unmakeFunc, ply)
		if score >= beta {
			return beta
		}
		alpha = utils.Max(alpha, score)
	}
	return alpha
}

func filterCaptures(moves []chess.Move, board *chess.Board) []chess.Move {
	var filteredCaptures []chess.Move
	for i := 0; i < len(moves); i++ {
		if chess.IsCapture(moves[i], board) {
			filteredCaptures = append(filteredCaptures, moves[i])
		}
	}
	return filteredCaptures
}
