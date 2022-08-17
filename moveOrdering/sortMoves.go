package moveOrdering

import (
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

var FollowPV bool = true

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, color int, ply int) int {
	// PV Move : 2000
	// Killer Moves : 900-1650
	// History Move : value*100
	// MVV-LVA : 0-650
	// 0-60 Pawn - King
	if move == bestmove {
		return 3000
	}
	if ply > 0 && FollowPV {
		FollowPV = false
		if &move == &pvMove {
			FollowPV = true
			return 2000
		}
	}
	if KillerMoves[0][ply] == move {
		return getMVV_LVA(move, board) + 1000
	}
	if KillerMoves[1][ply] == move {
		return getMVV_LVA(move, board) + 900
	}
	piece, _ := utils.GetPiece(move.From(), board)
	historyMove := GetHistoryMove(color, piece, move.From())
	if historyMove != 0 {
		return historyMove * 100
	}
	return getMVV_LVA(move, board)
}

func SortMoves(moves []chess.Move, board *chess.Board, pvTable [64]chess.Move, bestmove chess.Move, ply int) []chess.Move {
	pvMove := pvTable[ply]
	color := -1
	if board.Wtomove {
		color = 1
	}
	sort.Slice(moves, func(a, b int) bool {
		valueA := valueMove(board, moves[a], pvMove, bestmove, color, ply)
		valueB := valueMove(board, moves[b], pvMove, bestmove, color, ply)
		return valueA > valueB
	})
	return moves
}
