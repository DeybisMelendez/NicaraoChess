package search

import (
	"math"

	chess "github.com/dylhunn/dragontoothmg"
)

var FutilityMargin = [4]int{0, 90, 320, 500}

func IsFutilityPruning(staticEval int, depth int, alpha int, board *chess.Board, inCheck bool, isCapture bool, isPromotion bool) bool {
	if !isCapture && !inCheck && !isPromotion && !board.OurKingInCheck() {
		var value int = staticEval + FutilityMargin[depth]
		if int(math.Abs(float64(alpha))) < (MateScore-500) && value <= alpha {
			return true
		}
	}
	return false
}
