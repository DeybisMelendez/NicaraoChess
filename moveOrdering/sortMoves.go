package moveOrdering

import (
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

var FollowPV bool

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, isWhite bool, ply int) int {
	if move == bestmove {
		return 5000
	} else if move == pvMove {
		return 4000
	} else if chess.IsCapture(move, board) || move.Promote() != chess.Nothing {
		return GetMVV_LVA(move, board) + 3000
	} else if KillerMoves[0][ply] == move {
		return 2000
	} else if KillerMoves[1][ply] == move {
		return 1000
	} else {
		piece, _ := utils.GetPiece(move.From(), board)
		return GetHistoryMove(isWhite, piece, move.To()) + GetMVV_LVA(move, board)
	}
}

func SortMoves(moves []chess.Move, board *chess.Board, pvMove chess.Move, bestmove chess.Move, ply int) {
	sort.Slice(moves, func(a, b int) bool {
		valueA := valueMove(board, moves[a], pvMove, bestmove, board.Wtomove, ply)
		valueB := valueMove(board, moves[b], pvMove, bestmove, board.Wtomove, ply)
		return valueA > valueB
	})
}
