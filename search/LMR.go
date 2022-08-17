package search

// Late Move Reduction
import (
	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 6

func pvReduction(depth int) int {
	//return depth - 2
	if depth > 2 {
		return depth - 2
		//return int(float32(depth) / 3.0)
	}
	return depth - 1
}

func isLMROk(board *chess.Board, move chess.Move) bool {
	//var isNotKillerMove bool = !moveOrdering.IsKillerMove(move, Ply)
	var notCheck bool = !board.OurKingInCheck()
	var isNotCapture bool = !chess.IsCapture(move, board)
	var isNotPromotion bool = move.Promote() == chess.Nothing
	return notCheck && isNotCapture && isNotPromotion
}
