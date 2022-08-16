package moveOrdering

import (
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, color int, ply int) int {
	// PV Move : 2000
	// Killer Moves : 900-1650
	// History Move : value*10
	// MVV-LVA : 0-650
	if move == bestmove {
		return 3000
	}
	if move == pvMove {
		return 2000
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
		return historyMove * 10
	}
	return getMVV_LVA(move, board)
}

func SortMoves(moves []chess.Move, board *chess.Board, pvTable [64]chess.Move, bestmove chess.Move, color int, ply int) []chess.Move {
	pvMove := pvTable[ply]
	sort.Slice(moves, func(a, b int) bool {
		valueA := valueMove(board, moves[a], pvMove, bestmove, color, ply)
		valueB := valueMove(board, moves[b], pvMove, bestmove, color, ply)
		return valueA > valueB
	})
	return moves
}
