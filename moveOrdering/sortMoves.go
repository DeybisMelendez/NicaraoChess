package moveOrdering

import (
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

// var FollowPV bool
var Material = [7]int{0, 1, 3, 3, 5, 9, 10}

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, isWhite bool, ply int) int {
	if move == pvMove {
		return 6000
	} else if move == bestmove {
		return 5000
	} else if chess.IsCapture(move, board) || move.Promote() != chess.Nothing {
		piece, _ := utils.GetPiece(move.From(), board)
		capture, _ := utils.GetPiece(move.To(), board)
		promo := move.Promote()
		if Material[capture] > Material[piece] || promo != chess.Nothing {
			return 4500 + Material[capture] - Material[piece] + Material[promo]
		}
		if capture == piece {
			return 4000
		}
		return GetMVV_LVA(move, board)
	} else if KillerMoves[0][ply] == move {
		return 3000
	} else if KillerMoves[1][ply] == move {
		return 2000
	} else {
		return GetHistoryMove(isWhite, move) + 1000
	}
}

func SortMoves(moves []chess.Move, board *chess.Board, pvMove chess.Move, bestmove chess.Move, ply int) {
	sort.Slice(moves, func(a, b int) bool {
		valueA := valueMove(board, moves[a], pvMove, bestmove, board.Wtomove, ply)
		valueB := valueMove(board, moves[b], pvMove, bestmove, board.Wtomove, ply)
		return valueA > valueB
	})
}
