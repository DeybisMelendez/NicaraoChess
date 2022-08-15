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

func GetPiece(square uint8, board *chess.Board) (int, bool) {
	squareMask := uint64(1) << square
	piece := chess.Nothing
	w := &board.White
	b := &board.Black
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
	} else if b.Pawns&squareMask != 0 {
		piece = chess.Pawn
		isWhite = false
	} else if b.Knights&squareMask != 0 {
		piece = chess.Knight
		isWhite = false
	} else if b.Bishops&squareMask != 0 {
		piece = chess.Bishop
		isWhite = false
	} else if b.Rooks&squareMask != 0 {
		piece = chess.Rook
		isWhite = false
	} else if b.Queens&squareMask != 0 {
		piece = chess.Queen
		isWhite = false
	} else if b.Kings&squareMask != 0 {
		piece = chess.King
		isWhite = false
	}
	return piece, isWhite
}

func Make(board *chess.Board, move chess.Move, ply *int, nodes *int) func() {
	*ply++
	*nodes++
	//eval.UpdateMaterial(board, move)
	return board.Apply(move)
}

func Unmake(f func(), ply *int) {
	*ply--
	//eval.RevertMaterial()
	f()
}
