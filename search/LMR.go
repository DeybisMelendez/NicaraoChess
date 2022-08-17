package search

// Late Move Reduction
import (
	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 6

func pvReduction(depth int) int {
	//return depth - 2
	if depth > 2 {
		return int(float32(depth) / 3.0)
	}
	return depth - 1
}

func isLMROk(board *chess.Board, move chess.Move) bool {
	var check bool = board.OurKingInCheck()
	var isCapture bool = chess.IsCapture(move, board)
	return check && isCapture
}
