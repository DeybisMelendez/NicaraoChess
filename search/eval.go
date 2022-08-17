package search

import (
	"nicarao/moveOrdering"
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

func Evaluate(board *chess.Board) int {
	moves := board.GenerateLegalMoves()
	if len(moves) == 0 {
		if board.OurKingInCheck() {
			//Checkmate
			return -MateScore + Ply
		} else {
			//Stalemate
			return 0
		}
	}
	eval := 0
	//GetMaterial(board)
	eval += ValueMaterial(board)
	eval += PSTEval(board)
	if board.Wtomove {
		return eval
	}
	return -eval //eval * color
}

func Quiesce(board *chess.Board, alpha int, beta int) int {
	if isTimeToStop() {
		return 0
	}
	standPat := Evaluate(board)
	if standPat > beta {
		return beta
	}
	alpha = utils.Max(alpha, standPat)
	moves := filterCaptures(board.GenerateLegalMoves(), board)
	var score int = 0
	for i := 0; i < len(moves); i++ {
		if isTimeToStop() {
			return 0
		}
		unmakeFunc := Make(board, moves[i])
		score = -Quiesce(board, -beta, -alpha)
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(moves[i])
			alpha = score
		}
		if score >= beta {
			return beta
		}

		//alpha = utils.Max(alpha, score)
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
	filteredCaptures = moveOrdering.SortMoves(filteredCaptures, board, PVTable[0], 0, Ply)
	return filteredCaptures
}
