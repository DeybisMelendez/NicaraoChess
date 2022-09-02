package moveOrdering

import (
	chess "github.com/dylhunn/dragontoothmg"
)

var historyMoves [2][64][64]int

func StoreHistoryMove(move chess.Move, Wtomove bool, depth int) {
	//piece, _ := utils.GetPiece(move.To(), board)
	if Wtomove {
		historyMoves[0][move.From()][move.To()] += depth * depth
	} else {
		historyMoves[1][move.From()][move.To()] += depth * depth
	}
}

func GetHistoryMove(isWhite bool, move chess.Move) int {
	if isWhite {
		return historyMoves[0][move.From()][move.To()]
	}
	return historyMoves[1][move.From()][move.To()]
}

func ResetHistoryMove() {
	var newHistoryMove [2][64][64]int
	historyMoves = newHistoryMove
	/*for i := 0; i < 2; i++ {
		for j := 0; j < 7; j++ {
			for k := 0; k < 64; k++ {
				historyMoves[i][j][k] = 0
			}
		}
	}*/
}
