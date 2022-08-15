package moveOrdering

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

var historyMoves [2][6][64]int

func StoreHistoryMove(move chess.Move, board *chess.Board, depth int, color int) {
	piece, _ := utils.GetPiece(move.From(), board)
	if color == 1 {
		historyMoves[0][piece-1][move.From()] += depth * depth
	} else {
		historyMoves[1][piece-1][move.From()] += depth * depth
	}
}

func GetHistoryMove(color int, piece int, square uint8) int {
	if color == 1 {
		return historyMoves[0][piece-1][square]
	}
	return historyMoves[1][piece-1][square]
}
