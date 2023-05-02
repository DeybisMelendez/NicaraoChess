package engine

import (
	"math"
	"math/bits"

	chess "github.com/dylhunn/dragontoothmg"
)

var DistanceBonus [64][64]int

var KING_TROPISM [7][64][64]int16
var diagNW = [64]int{
	0, 1, 2, 3, 4, 5, 6, 7,
	1, 2, 3, 4, 5, 6, 7, 8,
	2, 3, 4, 5, 6, 7, 8, 9,
	3, 4, 5, 6, 7, 8, 9, 10,
	4, 5, 6, 7, 8, 9, 10, 11,
	5, 6, 7, 8, 9, 10, 11, 12,
	6, 7, 8, 9, 10, 11, 12, 13,
	7, 8, 9, 10, 11, 12, 13, 14,
}
var diagNE = [64]int{
	7, 6, 5, 4, 3, 2, 1, 0,
	8, 7, 6, 5, 4, 3, 2, 1,
	9, 8, 7, 6, 5, 4, 3, 2,
	10, 9, 8, 7, 6, 5, 4, 3,
	11, 10, 9, 8, 7, 6, 5, 4,
	12, 11, 10, 9, 8, 7, 6, 5,
	13, 12, 11, 10, 9, 8, 7, 6,
	14, 13, 12, 11, 10, 9, 8, 7,
}
var bonusDiagDistance = [15]int{5, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

const MOBILITY_WEIGHT = 5

const WhiteColorBoard uint64 = 0xAA55AA55AA55AA55
const BlackColorBoard uint64 = 0x55AA55AA55AA55AA
const Center uint64 = 0x1818000000
const ExtendedCenter uint64 = 0x182424180000

var PassedPawnBonus = [8]int{0, 10, 20, 40, 60, 80, 100, 200}

func doubledPawns(board *chess.Board, isWhite bool, square uint8) int16 {
	var pawns uint64
	if isWhite {
		pawns = board.White.Pawns
	} else {
		pawns = board.Black.Pawns
	}
	if bits.OnesCount64(pawns&FileMask[square]) > 1 {
		return 10
	}
	return 0
}
func isolatedPawns(board *chess.Board, isWhite bool, square uint8) int16 {
	var pawns uint64
	if isWhite {
		pawns = board.White.Pawns
	} else {
		pawns = board.Black.Pawns
	}
	if pawns&IsolatedMask[square] == 0 {
		return 10
	}
	return 0
}

/*func passedPawn(board *chess.Board, isWhite bool, square uint8) int16 {
	var pawns uint64
	if isWhite {
		pawns = board.White.Pawns
	} else {
		pawns = board.Black.Pawns
	}
}*/

func bishopPair(board *chess.Board, isWhite bool) int16 {
	var bishops uint64
	if isWhite {
		bishops = board.White.Bishops
	} else {
		bishops = board.Black.Bishops
	}
	if bits.OnesCount64(bishops) > 1 {
		return 10
	}
	return 0
}

func badBishop(board *chess.Board, isWhite bool, square uint8) int16 {
	var pawns uint64
	if isWhite {
		pawns = board.White.Pawns
	} else {
		pawns = board.Black.Pawns
	}
	squareMask := uint64(1) << square
	// White
	if squareMask&WhiteColorBoard > 0 {
		if bits.OnesCount64(pawns&WhiteColorBoard) > 3 {
			return 10
		}
	} else if squareMask&BlackColorBoard > 0 {
		if bits.OnesCount64(pawns&BlackColorBoard) > 3 {
			return 10
		}
	}
	return 0
}
func goodKnight(board *chess.Board) int16 {
	var pawns int16 = int16(uint64(bits.OnesCount64(board.White.Pawns | board.Black.Pawns)))
	return pawns * 2
}

func goodRook(board *chess.Board) int16 {
	var pawns int16 = int16(uint64(bits.OnesCount64(board.White.Pawns | board.Black.Pawns)))
	return 32 - (pawns * 2)
}

func rookToQueen(board *chess.Board, isWhite bool, square uint8) int16 {
	var queen uint64
	if isWhite {
		queen = board.Black.Queens
	} else {
		queen = board.White.Queens
	}
	if queen&FileMask[square] > 0 {
		return 10
	}
	return 0
}

// https://www.chessprogramming.org/Knight_Pattern

func mobility(pieceType chess.Piece, isWhite bool, square uint8, board *chess.Board, allPieces uint64) int16 {
	var myPieces uint64
	if isWhite {
		myPieces = board.White.All
	} else {
		myPieces = board.Black.All
	}
	switch pieceType {
	case chess.Knight:
		return int16(bits.OnesCount64(knightMasks[square]&(^myPieces)) * MOBILITY_WEIGHT)
	case chess.Bishop:
		return int16(bits.OnesCount64(chess.CalculateBishopMoveBitboard(square, allPieces)&(^myPieces)) * MOBILITY_WEIGHT)
	case chess.Rook:
		return int16(bits.OnesCount64(chess.CalculateRookMoveBitboard(square, allPieces)&(^myPieces)) * MOBILITY_WEIGHT)
	}
	return 0
}

func kingSafety(isWhite bool, square uint8, board *chess.Board, allPieces uint64) int16 {
	var myPieces uint64
	if isWhite {
		myPieces = board.White.All
	} else {
		myPieces = board.Black.All
	}
	var score int16 = int16(bits.OnesCount64(chess.CalculateBishopMoveBitboard(square, allPieces)&(^myPieces)) +
		bits.OnesCount64(chess.CalculateRookMoveBitboard(square, allPieces)&(^myPieces)))
	return score + score
}

func BadQueen(board *chess.Board, byBlack bool, square uint8) int16 {
	if board.UnderDirectAttack(byBlack, square) {
		return 30
	}
	return 0
}

var Row = [64]int{
	7, 7, 7, 7, 7, 7, 7, 7,
	6, 6, 6, 6, 6, 6, 6, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0}

var Col = [64]int{
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
	0, 1, 2, 3, 4, 5, 6, 7,
}

func SetDist() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			DistanceBonus[i][j] = 14 - int(math.Abs(float64(Col[i])-float64(Col[j]))+
				math.Abs(float64(Row[i])-float64(Row[j])))
			KING_TROPISM[chess.Knight][i][j] = int16(DistanceBonus[i][j])
			KING_TROPISM[chess.Bishop][i][j] += int16(bonusDiagDistance[int(math.Abs(float64(diagNE[i])-float64(diagNE[j])))])
			KING_TROPISM[chess.Bishop][i][j] += int16(bonusDiagDistance[int(math.Abs(float64(diagNW[i])-float64(diagNW[j])))])
			KING_TROPISM[chess.Rook][i][j] = int16(DistanceBonus[i][j] / 2)
			KING_TROPISM[chess.Queen][i][j] = int16((DistanceBonus[i][j] * 5) / 2)
		}
	}
}
