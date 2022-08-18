package search

import (
	"math"

	chess "github.com/dylhunn/dragontoothmg"
)

var FutilityMargin = [4]int{0, 100, 300, 500}
var FutilityPruning bool

func CheckFutilityPruning(staticEval int, depth int, alpha int) {
	if depth < 4 {
		if int(math.Abs(float64(alpha))) < MateScore && staticEval+FutilityMargin[depth] <= alpha {
			FutilityPruning = true
		}
	}
}

func IsFutilityPruning(board *chess.Board, move chess.Move, i int) bool {
	if i > 0 && FutilityPruning && !chess.IsCapture(move, board) &&
		!board.OurKingInCheck() && move.Promote() == chess.Nothing {
		return true
	}
	return false
}
