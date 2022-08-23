package moveOrdering

import (
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

var FollowPV bool

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, color int, ply int) int {
	score := 0
	piece, _ := utils.GetPiece(move.From(), board)
	historyMove := GetHistoryMove(color, piece, move.To())
	if move == bestmove {
		score = 6000000
	} else if ply > 0 && FollowPV {
		FollowPV = false
		if move == pvMove {
			FollowPV = true
			score = 5000000
		}
	} else if KillerMoves[0][ply] == move {
		score = 4000000
	} else if KillerMoves[1][ply] == move {
		score = 3000000
	} else if historyMove != 0 {
		score = 2000000 + historyMove
	} else {
		score = GetMVV_LVA(move, board)
	}
	return score
}

func SortMoves(moves []chess.Move, board *chess.Board, pvMove chess.Move, bestmove chess.Move, ply int) {
	color := -1
	if board.Wtomove {
		color = 1
	}
	sort.Slice(moves, func(a, b int) bool {
		valueA := valueMove(board, moves[a], pvMove, bestmove, color, ply)
		valueB := valueMove(board, moves[b], pvMove, bestmove, color, ply)
		return valueA > valueB
	})
}
