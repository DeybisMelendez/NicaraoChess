package utils

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

// Contains tells whether a contains x.

func Checkmate(board *chess.Board) bool {
	if board.OurKingInCheck() && len(board.GenerateLegalMoves()) == 0 {
		return true
	}
	return false
}

func GetPiece(square uint8, board *chess.Board) (int, bool) {
	squareMask := uint64(1) << square
	piece := chess.Nothing
	w := board.White
	b := board.Black
	var isWhite bool = true
	if w.Pawns&squareMask != 0 {
		piece = chess.Pawn
	} else if w.Knights&squareMask != 0 {
		piece = chess.Knight
	} else if w.Bishops&squareMask != 0 {
		piece = chess.Bishop
	} else if w.Rooks&squareMask != 0 {
		piece = chess.Rook
	} else if w.Queens&squareMask != 0 {
		piece = chess.Queen
	} else if w.Kings&squareMask != 0 {
		piece = chess.King
	} else {
		isWhite = false
		if b.Pawns&squareMask != 0 {
			piece = chess.Pawn
		} else if b.Knights&squareMask != 0 {
			piece = chess.Knight
		} else if b.Bishops&squareMask != 0 {
			piece = chess.Bishop
		} else if b.Rooks&squareMask != 0 {
			piece = chess.Rook
		} else if b.Queens&squareMask != 0 {
			piece = chess.Queen
		} else if b.Kings&squareMask != 0 {
			piece = chess.King
		}
	}
	return piece, isWhite
}

func Reverse(slc [64]int) [64]int {
	reversed := [64]int{}
	for i := 0; i < len(slc); i++ {
		// reverse the order
		reversed[i] = slc[len(slc)-1-i]
	}
	return reversed
}

func SetBits(mask uint64, square uint64) uint64 {
	return ((mask) | (1 << square))
}
