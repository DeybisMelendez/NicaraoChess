package moveOrdering

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

var historyMoves [2][7][64]int

func StoreHistoryMove(move chess.Move, board *chess.Board, depth int) {
	color := -1
	if board.Wtomove {
		color = 1
	}
	piece, _ := utils.GetPiece(move.From(), board)
	if color == 1 {
		historyMoves[0][piece][move.From()] += depth * depth
	} else {
		historyMoves[1][piece][move.From()] += depth * depth
	}
}

func GetHistoryMove(color int, piece int, square uint8) int {
	if color == 1 {
		return historyMoves[0][piece][square]
	}
	return historyMoves[1][piece][square]
}

func ResetHistoryMove() {
	var newHistoryMove [2][7][64]int
	historyMoves = newHistoryMove
	/*for i := 0; i < 2; i++ {
		for j := 0; j < 7; j++ {
			for k := 0; k < 64; k++ {
				historyMoves[i][j][k] = 0
			}
		}
	}*/
}
