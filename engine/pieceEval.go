package engine

import (
	"math/bits"

	chess "github.com/dylhunn/dragontoothmg"
)

var Material = [7]int{0, 1, 3, 3, 5, 9, 10}
var MaterialOpening = [7]int{0, 90, 320, 330, 500, 900, 10000}
var MaterialEndgame = [7]int{0, 100, 320, 350, 550, 900, 10000}
var MaterialScore = [2][7]int{MaterialOpening, MaterialEndgame} //Opening, Endgame

const WhiteColorBoard uint64 = 0xAA55AA55AA55AA55
const BlackColorBoard uint64 = 0x55AA55AA55AA55AA
const Center uint64 = 0x1818000000
const ExtendedCenter uint64 = 0x182424180000

var PassedPawnBonus = [8]int{0, 10, 20, 40, 60, 80, 100, 200}

func DoublePawns(pawns uint64, square uint8) int {
	if bits.OnesCount64(pawns&FileMask[square]) > 1 {
		return 10
	}
	return 0
}
func IsolatedPawns(pawns uint64, square uint8) int {
	if pawns&IsolatedMask[square] == 0 {
		return 10
	}
	return 0
}

func GoodKnight(pawns int) int {
	return pawns * 2
}

func BishopPair(bishops uint64) int {
	if bits.OnesCount64(bishops) > 1 {
		return 20
	}
	return 0
}

func BadBishop(square uint8, pawns uint64) int {
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

func GoodRook(pawns int) int {
	return 32 - (pawns * 2)
}

// https://www.chessprogramming.org/Knight_Pattern

func MobilityRook(square uint8, allPieces uint64, myPieces uint64) int {
	return bits.OnesCount64(chess.CalculateRookMoveBitboard(square, allPieces)&(^myPieces)) * 2
}

func MobilityBishop(square uint8, allPieces uint64, myPieces uint64) int {
	return bits.OnesCount64(chess.CalculateBishopMoveBitboard(square, allPieces)&(^myPieces)) * 2
}

func MobilityKnight(square uint8, myPieces uint64) int {
	return bits.OnesCount64(knightMasks[square]&(^myPieces)) * 2
}

func BadQueen(board *chess.Board, byBlack bool, square uint8) int {
	if board.UnderDirectAttack(byBlack, square) {
		return 30
	}
	return 0
}

func BadKing(square uint8, allPieces uint64, myPieces uint64, isEndgame bool) int {
	score := MobilityRook(square, allPieces, myPieces) + MobilityBishop(square, allPieces, myPieces)
	if isEndgame {
		return 0
	}
	return score
}
