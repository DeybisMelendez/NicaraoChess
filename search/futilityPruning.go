package search

import (
	"math"

	chess "github.com/dylhunn/dragontoothmg"
)

var FutilityMargin = 90

func IsFutilityPruning(staticEval int, alpha int, board *chess.Board, inCheck bool, isCapture bool, isPromotion bool) bool {
	if !isCapture && !inCheck && !isPromotion && !board.OurKingInCheck() {
		var value int = staticEval + FutilityMargin
		if int(math.Abs(float64(alpha))) < (MateScore-500) && value <= alpha {
			return true
		}
	}
	return false
}
