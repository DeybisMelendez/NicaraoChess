package moveOrdering

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

var historyMoves [2][7][64]int

func StoreHistoryMove(move chess.Move, board *chess.Board, depth int) {
	if !chess.IsCapture(move, board) {
		piece, _ := utils.GetPiece(move.To(), board)
		if board.Wtomove {
			historyMoves[0][piece][move.To()] += depth * depth
		} else {
			historyMoves[1][piece][move.To()] += depth * depth
		}
	}
}

func GetHistoryMove(isWhite bool, piece int, square uint8) int {
	if isWhite {
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
