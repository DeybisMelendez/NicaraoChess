package search

import (
	"math"

	chess "github.com/dylhunn/dragontoothmg"
)

var FutilityMargin = [4]int{0, 100, 300, 500}
var FutilityPruning bool

func CheckFutilityPruning(staticEval int, depth int, alpha int) {
	//FutilityPruning = false
	if depth < 4 &&
		int(math.Abs(float64(alpha))) < MateScore &&
		staticEval+FutilityMargin[depth] <= alpha {
		FutilityPruning = true
	}
}

func IsFutilityPruning(board *chess.Board, i int, inCheck bool, isCapture bool) bool {
	if i > 0 && FutilityPruning && !isCapture && !inCheck && !board.OurKingInCheck() {
		return true
	}
	return false
}
