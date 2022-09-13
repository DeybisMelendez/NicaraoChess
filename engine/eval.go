package engine

import (
	"math/bits"

	chess "github.com/dylhunn/dragontoothmg"
)

const WHITE, BLACK = 0, 1
const OPENING, ENDGAME = 0, 1

var KnightPhase int = 1
var BishopPhase int = 1
var RookPhase int = 2
var QueenPhase int = 4

var TotalPhase int = KnightPhase*4 +
	BishopPhase*4 + RookPhase*4 + QueenPhase*2 //+PawnPhase*16
var phase int = TotalPhase

func Evaluate(board *chess.Board, turn int) int {
	opening, endgame := 0, 0
	allPieces := board.White.All | board.Black.All
	allPawnCount := bits.OnesCount64(board.White.Pawns | board.Black.Pawns)
	// Phase
	phase = (bits.OnesCount64(board.White.Knights&board.Black.Knights) * KnightPhase) +
		(bits.OnesCount64(board.White.Bishops&board.Black.Bishops) * BishopPhase) +
		(bits.OnesCount64(board.White.Rooks&board.Black.Rooks) * RookPhase) +
		(bits.OnesCount64(board.White.Queens&board.Black.Queens) * QueenPhase)
	// Bishop pair
	opening += BishopPair(board.White.Bishops) - BishopPair(board.Black.Bishops)

	//Good Knight
	whiteGoodKnight := bits.OnesCount64(board.White.Knights) * GoodKnight(allPawnCount)
	blackGoodKnight := bits.OnesCount64(board.Black.Knights) * GoodKnight(allPawnCount)
	opening += whiteGoodKnight - blackGoodKnight

	//Good Rook
	whiteGoodRook := bits.OnesCount64(board.White.Rooks) * GoodRook(allPawnCount)
	blackGoodRook := bits.OnesCount64(board.White.Rooks) * GoodRook(allPawnCount)
	opening += whiteGoodRook - blackGoodRook

	opening = endgame
	//Piece Evaluation
	pieces := board.White.Pawns
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening += MaterialScore[OPENING][chess.Pawn]
		endgame += MaterialScore[ENDGAME][chess.Pawn]
		opening += PST[OPENING][chess.Pawn][ReversedBoard[square]]
		endgame += PST[ENDGAME][chess.Pawn][ReversedBoard[square]]
		val := DoublePawns(board.White.Pawns, square)
		val += IsolatedPawns(board.White.Pawns, square)
		opening -= val
		endgame -= val
		pieces &= pieces - 1
	}
	pieces = board.White.Knights
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening += MaterialScore[OPENING][chess.Knight]
		endgame += MaterialScore[ENDGAME][chess.Knight]
		opening += PST[OPENING][chess.Knight][ReversedBoard[square]]
		endgame += PST[ENDGAME][chess.Knight][ReversedBoard[square]]
		mobility := MobilityKnight(square, board.White.Knights)
		opening += mobility
		endgame += mobility
		pieces &= pieces - 1
	}
	pieces = board.White.Bishops
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening += MaterialScore[OPENING][chess.Bishop]
		endgame += MaterialScore[ENDGAME][chess.Bishop]
		opening += PST[OPENING][chess.Bishop][ReversedBoard[square]]
		endgame += PST[ENDGAME][chess.Bishop][ReversedBoard[square]]
		mobility := MobilityBishop(square, allPieces, board.White.All)
		opening += mobility
		endgame += mobility
		pieces &= pieces - 1
	}
	pieces = board.White.Rooks
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening += MaterialScore[OPENING][chess.Rook]
		endgame += MaterialScore[ENDGAME][chess.Rook]
		opening += PST[OPENING][chess.Rook][ReversedBoard[square]]
		endgame += PST[ENDGAME][chess.Rook][ReversedBoard[square]]
		mobility := MobilityRook(square, allPieces, board.White.All)
		opening += mobility
		endgame += mobility
		pieces &= pieces - 1
	}
	pieces = board.White.Queens
	for pieces != 0 {
		square := bits.TrailingZeros64(pieces)
		opening += MaterialScore[OPENING][chess.Queen]
		endgame += MaterialScore[ENDGAME][chess.Queen]
		opening += PST[OPENING][chess.Queen][ReversedBoard[square]]
		endgame += PST[ENDGAME][chess.Queen][ReversedBoard[square]]
		pieces &= pieces - 1
	}
	pieces = board.Black.Pawns
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening -= MaterialScore[OPENING][chess.Pawn]
		endgame -= MaterialScore[ENDGAME][chess.Pawn]
		opening -= PST[OPENING][chess.Pawn][square]
		endgame -= PST[ENDGAME][chess.Pawn][square]
		val := DoublePawns(board.White.Pawns, square)
		val += IsolatedPawns(board.White.Pawns, square)
		opening += val
		endgame += val
		pieces &= pieces - 1
	}
	pieces = board.Black.Knights
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening -= MaterialScore[OPENING][chess.Knight]
		endgame -= MaterialScore[ENDGAME][chess.Knight]
		opening -= PST[OPENING][chess.Knight][square]
		endgame -= PST[ENDGAME][chess.Knight][square]
		mobility := MobilityKnight(square, board.Black.All)
		opening -= mobility
		endgame -= mobility
		pieces &= pieces - 1
	}
	pieces = board.Black.Bishops
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening -= MaterialScore[OPENING][chess.Bishop]
		endgame -= MaterialScore[ENDGAME][chess.Bishop]
		opening -= PST[OPENING][chess.Bishop][square]
		endgame -= PST[ENDGAME][chess.Bishop][square]
		mobility := MobilityBishop(square, allPieces, board.Black.All)
		opening -= mobility
		endgame -= mobility
		pieces &= pieces - 1
	}
	pieces = board.Black.Rooks
	for pieces != 0 {
		square := uint8(bits.TrailingZeros64(pieces))
		opening -= MaterialScore[OPENING][chess.Rook]
		endgame -= MaterialScore[ENDGAME][chess.Rook]
		opening -= PST[OPENING][chess.Rook][square]
		endgame -= PST[ENDGAME][chess.Rook][square]
		mobility := MobilityRook(square, allPieces, board.Black.All)
		opening -= mobility
		endgame -= mobility
		pieces &= pieces - 1
	}
	pieces = board.Black.Queens
	for pieces != 0 {
		square := bits.TrailingZeros64(pieces)
		opening -= MaterialScore[OPENING][chess.Queen]
		endgame -= MaterialScore[ENDGAME][chess.Queen]
		opening -= PST[OPENING][chess.Queen][square]
		endgame -= PST[ENDGAME][chess.Queen][square]
		pieces &= pieces - 1
	}
	king := uint8(bits.TrailingZeros64(board.White.Kings))
	opening += PST[OPENING][chess.King][ReversedBoard[king]]
	endgame += PST[ENDGAME][chess.King][ReversedBoard[king]]
	opening -= BadKing(king, allPieces, board.White.All, false)
	endgame -= BadKing(king, allPieces, board.White.All, true)
	king = uint8(bits.TrailingZeros64(board.Black.Kings))
	opening -= PST[OPENING][chess.King][king]
	endgame -= PST[ENDGAME][chess.King][king]
	opening += BadKing(king, allPieces, board.Black.All, false)
	endgame += BadKing(king, allPieces, board.Black.All, true)

	if isEndgame(board) {
		if IsDraw(board) {
			return 0
		}
		return endgame * turn // Endgame
	}
	phase = ((phase * 256) + (TotalPhase / 2)) / TotalPhase
	score := ((opening * (256 - phase)) + (opening * phase)) / 256
	return score * turn
}

func isEndgame(board *chess.Board) bool {
	queens := board.White.Queens | board.Black.Queens
	if queens == 0 {
		return true
	} else {
		minorPieces := board.White.Knights | board.White.Bishops | board.White.Rooks |
			board.Black.Knights | board.Black.Bishops | board.Black.Rooks
		if bits.OnesCount64(minorPieces) < 4 && bits.OnesCount64(queens) < 2 {
			return true
		}
	}
	return false
}
