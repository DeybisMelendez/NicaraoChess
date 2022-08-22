package search

import (
	"nicarao/moveOrdering"

	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 6

func pvReduction(depth int) int {
	return depth / 3
}

func isLMROk(board *chess.Board, inCheck bool, isCapture bool, move chess.Move) bool {
	isKillerMove := moveOrdering.KillerMoves[0][Ply] == move || moveOrdering.KillerMoves[1][Ply] == move
	return !inCheck && !isCapture && !isKillerMove && move.Promote() == chess.Nothing
}
