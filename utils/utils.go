package utils

import (
	"fmt"

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

/*func Checkmate(board *chess.Board) bool {
	if board.OurKingInCheck() && len(board.GenerateLegalMoves()) == 0 {
		return true
	}
	return false
}*/

func GetPiece(square uint8, board *chess.Board) (int, bool) {
	squareMask := uint64(1) << square
	if board.White.All&squareMask != 0 {
		if board.White.Pawns&squareMask != 0 {
			return chess.Pawn, true
		} else if board.White.Knights&squareMask != 0 {
			return chess.Knight, true
		} else if board.White.Bishops&squareMask != 0 {
			return chess.Bishop, true
		} else if board.White.Rooks&squareMask != 0 {
			return chess.Rook, true
		} else if board.White.Queens&squareMask != 0 {
			return chess.Queen, true
		} else if board.White.Kings&squareMask != 0 {
			return chess.King, true
		}
	} else {
		if board.Black.Pawns&squareMask != 0 {
			return chess.Pawn, false
		} else if board.Black.Knights&squareMask != 0 {
			return chess.Knight, false
		} else if board.Black.Bishops&squareMask != 0 {
			return chess.Bishop, false
		} else if board.Black.Rooks&squareMask != 0 {
			return chess.Rook, false
		} else if board.Black.Queens&squareMask != 0 {
			return chess.Queen, false
		} else if board.Black.Kings&squareMask != 0 {
			return chess.King, false
		}
	}
	return chess.Nothing, false
}

/*func Reverse(slc [64]int) [64]int {
	reversed := [64]int{}
	for i := 0; i < len(slc); i++ {
		// reverse the order
		reversed[i] = slc[len(slc)-1-i]
	}
	return reversed
}*/

func SetBits(mask uint64, square uint64) uint64 {
	return ((mask) | (1 << square))
}

func PrintBits(bits uint64) {
	for y := uint64(0); y < 8; y++ {
		line := ""
		for x := uint64(0); x < 8; x++ {
			if bits&(1<<((7-y)*8+x)) > 0 {
				line += "1"
			} else {
				line += "0"
			}
		}
		fmt.Println(line)
	}
}
