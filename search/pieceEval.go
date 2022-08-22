package search

import (
	"math/bits"
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

const WhiteColorBoard uint64 = 0xAA55AA55AA55AA55
const BlackColorBoard uint64 = 0x55AA55AA55AA55AA
const Center uint64 = 0x1818000000
const ExtendedCenter uint64 = 0x182424180000

var PassedPawnBonus = [8]int{0, 10, 20, 40, 100, 200, 300, 400}

var ranks = [8]uint64{0xff, 0xff00, 0xff0000, 0xff000000, 0xff00000000, 0xff0000000000, 0xff000000000000, 0xff00000000000000}
var files = [8]uint64{0x8080808080808080, 0x4040404040404040, 0x2020202020202020, 0x1010101010101010,
	0x808080808080808, 0x404040404040404, 0x202020202020202, 0x101010101010101}

var fileMask = [64]uint64{}
var rankMask = [64]uint64{}
var isolatedMask = [64]uint64{}
var whitePassedMask = [64]uint64{}
var blackPassedMask = [64]uint64{}

var getRank = [64]int{
	7, 7, 7, 7, 7, 7, 7, 7,
	6, 6, 6, 6, 6, 6, 6, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0}

func setFileRankMask(fileNumber int, rankNumber int) uint64 {
	var mask uint64 = 0
	for rank := uint64(0); rank < 8; rank++ {
		for file := uint64(0); file < 8; file++ {
			var square uint64 = rank*8 + file
			if fileNumber != -1 {
				if file == uint64(fileNumber) {
					mask = mask | utils.SetBits(mask, square)
				}
			} else if rankNumber != -1 {
				if rank == uint64(rankNumber) {
					mask = mask | utils.SetBits(mask, square)
				}
			}
		}
	}
	return mask
}

func InitEvaluationMask() {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			var square int = rank*8 + file
			fileMask[square] = fileMask[square] | setFileRankMask(file, -1)
			rankMask[square] = rankMask[square] | setFileRankMask(-1, rank)
			isolatedMask[square] = isolatedMask[square] | setFileRankMask(file-1, -1)
			isolatedMask[square] = isolatedMask[square] | setFileRankMask(file+1, -1)
			whitePassedMask[square] = whitePassedMask[square] | setFileRankMask(file-1, -1)
			whitePassedMask[square] = whitePassedMask[square] | setFileRankMask(file, -1)
			whitePassedMask[square] = whitePassedMask[square] | setFileRankMask(file+1, -1)
			for i := 0; i < (8 - rank); i++ {
				whitePassedMask[square] = whitePassedMask[square] & rankMask[(7-i)*8+file]
			}

			blackPassedMask[square] = blackPassedMask[square] | setFileRankMask(file-1, -1)
			blackPassedMask[square] = blackPassedMask[square] | setFileRankMask(file, -1)
			blackPassedMask[square] = blackPassedMask[square] | setFileRankMask(file+1, -1)
			for i := 0; i < rank+1; i++ {
				blackPassedMask[square] = blackPassedMask[square] & rankMask[(7-i)*8+file]
			}
		}
	}
}

func DoublePawns(pawns uint64, square uint8) int {
	if bits.OnesCount64(pawns&uint64(fileMask[square])) > 1 {
		return 10
	}
	return 0
}
func IsolatedPawns(pawns uint64, square uint8) int {
	if pawns&isolatedMask[square] == 0 {
		return 10
	}
	return 0
}

func PassedPawns(pawns uint64, square uint8, isWhite bool) int {
	if isWhite {
		if pawns&whitePassedMask[square] == 0 {
			return PassedPawnBonus[getRank[63-square]]
		} else {
			return 0
		}
	}
	if pawns&blackPassedMask[square] == 0 {
		return PassedPawnBonus[getRank[square]]
	}
	return 0
}

func BadKnight(pawns int) int {
	return 16 - pawns
}

func BishopPair(bishops uint64) int {
	if bits.OnesCount64(bishops) > 1 {
		return 10 // Con 2 alfiles sumará 20
	}
	return 0
}

func BadBishop(square uint8, pawns uint64) int {
	score := 0
	squareMask := uint64(1) << square
	// White
	if squareMask&WhiteColorBoard > 0 {
		if bits.OnesCount64(pawns&WhiteColorBoard) > 3 {
			score += 20
		}
	}
	if squareMask&BlackColorBoard > 0 {
		if bits.OnesCount64(pawns&BlackColorBoard) > 3 {
			score += 20
		}
	}
	return score
}

func GoodRook(pawns int) int {
	return 32 - (pawns * 2)
}

func MobilityRook(square uint8, allPieces uint64, myPieces uint64) int {
	return bits.OnesCount64(chess.CalculateRookMoveBitboard(square, allPieces)&(^myPieces)) * 2
}

func MobilityBishop(square uint8, allPieces uint64, myPieces uint64) int {
	return bits.OnesCount64(chess.CalculateBishopMoveBitboard(square, allPieces)&(^myPieces)) * 2
}

func BadQueen(board *chess.Board, byBlack bool, square uint8) int {
	if board.UnderDirectAttack(byBlack, square) {
		return 20
	}
	return 0
}

func BadKing(square uint8, allPieces uint64, myPieces uint64, isEndgame bool) int {
	return MobilityRook(square, allPieces, myPieces) + MobilityBishop(square, allPieces, myPieces)*2
}

func CenterPawn(square uint8, isEndgame bool) int {
	if !isEndgame { //Opening
		squareMask := uint64(1) << square
		if Center&squareMask != 0 {
			return 20
		}
		if ExtendedCenter&squareMask != 0 {
			return 10
		}
	}
	return 0
}