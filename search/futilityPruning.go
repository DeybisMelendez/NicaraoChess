package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

const FutilityMargin int = 100

func IsFutilityPruning(staticEval int, alpha int, board *chess.Board, inCheck bool, isCapture bool, isPromotion bool) bool {
	if !isCapture && !inCheck && !board.OurKingInCheck() && !isPromotion {
		var value int = staticEval + FutilityMargin
		if value < MateScore && value <= alpha {
			return true
		}
	}
	return false
}
