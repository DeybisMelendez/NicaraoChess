package search

// Late Move Reduction
import (
	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 3

func pvReduction(depth int) int {
	return depth - 2
}

func isLMROk(board *chess.Board, move chess.Move) bool {
	return !board.OurKingInCheck() && !chess.IsCapture(move, board) && move.Promote() == chess.Nothing
}
