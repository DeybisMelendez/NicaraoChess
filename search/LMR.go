package search

// Late Move Reduction
import (
	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 3

func pvReduction(depth int) int {
	return depth / 3
}

func isLMROk(board *chess.Board, move chess.Move) bool {
	var notCheck bool = !board.OurKingInCheck()
	var isNotCapture bool = !chess.IsCapture(move, board)
	var isNotPromotion bool = move.Promote() == chess.Nothing
	return notCheck && isNotCapture && isNotPromotion
}
