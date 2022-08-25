package moveOrdering

import (
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

var FollowPV bool

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, isWhite bool, ply int) int {
	score := GetMVV_LVA(move, board)
	piece, _ := utils.GetPiece(move.From(), board)
	historyMove := GetHistoryMove(isWhite, piece, move.To())
	if move == bestmove {
		score = 5000
	} else if ply > 0 && FollowPV {
		FollowPV = false
		if move == pvMove {
			FollowPV = true
			score = 4000
		}
	} else if chess.IsCapture(move, board) {
		score += 3000
	} else if KillerMoves[0][ply] == move {
		score = 2000
	} else if KillerMoves[1][ply] == move {
		score = 1000
	} else {
		score += historyMove
	}
	return score
}

func SortMoves(moves []chess.Move, board *chess.Board, pvMove chess.Move, bestmove chess.Move, ply int) {
	sort.Slice(moves, func(a, b int) bool {
		valueA := valueMove(board, moves[a], pvMove, bestmove, board.Wtomove, ply)
		valueB := valueMove(board, moves[b], pvMove, bestmove, board.Wtomove, ply)
		return valueA > valueB
	})
}
