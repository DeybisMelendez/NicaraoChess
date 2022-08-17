package search

// Late Move Reduction
import (
	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 6

func pvReduction(depth int) int {
	return int(float32(depth) / 3.0)
}

func isLMROk(board *chess.Board, move chess.Move) bool {
	var check bool = board.OurKingInCheck()
	var isCapture bool = chess.IsCapture(move, board)
	return check && isCapture
}
