package engine

import (
	"math/bits"

	chess "github.com/dylhunn/dragontoothmg"
)

func IsDraw(board *chess.Board) bool {
	if board.White.Pawns|board.Black.Pawns == 0 {
		if board.White.Rooks|board.Black.Rooks == 0 {
			white := board.White.Knights | board.White.Bishops
			black := board.Black.Knights | board.Black.Bishops
			if bits.OnesCount64(white|black) < 2 {
				return true
			}
			if bits.OnesCount64(white) == 1 && bits.OnesCount64(black) == 1 {
				return true
			}
			if (bits.OnesCount64(board.White.Knights) == 2 && bits.OnesCount64(black) < 2) ||
				(bits.OnesCount64(board.Black.Knights) == 2 && bits.OnesCount64(white) < 2) {
				return true
			}
		}
	}
	if bits.OnesCount64(board.White.All|board.Black.All) == 2 {
		return true
	}
	return false
}
