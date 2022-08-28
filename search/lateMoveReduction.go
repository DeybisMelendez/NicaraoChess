package search

import (
	"nicarao/moveOrdering"

	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 6

func pvReduction(depth int) int {
	return depth * 2 / 3
}
func isLMROk(board *chess.Board, inCheck bool, isCapture bool, move chess.Move) bool {
	return !inCheck && !isCapture && move.Promote() == chess.Nothing && !moveOrdering.IsKillerMove(move, Ply)
}
