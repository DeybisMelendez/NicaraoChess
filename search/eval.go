package search

import (
	"math/bits"
	"nicarao/moveOrdering"
	"nicarao/utils"

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
	incheck := board.OurKingInCheck()
	king := uint8(bits.TrailingZeros64(board.White.Kings))
	opening += PST[OPENING][chess.King][ReversedBoard[king]]
	endgame += PST[ENDGAME][chess.King][ReversedBoard[king]]
	opening -= BadKing(king, allPieces, board.White.All, false)
	opening -= AttackedKing(incheck, board.Wtomove, true)
	endgame -= BadKing(king, allPieces, board.White.All, true)
	endgame -= AttackedKing(incheck, board.Wtomove, false)
	king = uint8(bits.TrailingZeros64(board.Black.Kings))
	opening -= PST[OPENING][chess.King][king]
	endgame -= PST[ENDGAME][chess.King][king]
	opening += BadKing(king, allPieces, board.Black.All, false)
	opening += AttackedKing(incheck, board.Wtomove, true)
	endgame += BadKing(king, allPieces, board.Black.All, true)
	endgame += AttackedKing(incheck, board.Wtomove, false)

	//Pawn structure
	/*opening += bits.OnesCount64(board.White.Pawns&Center) * 20
	opening -= bits.OnesCount64(board.Black.Pawns&Center) * 20
	opening += bits.OnesCount64(board.White.Pawns&ExtendedCenter) * 10
	opening -= bits.OnesCount64(board.Black.Pawns&ExtendedCenter) * 10*/
	/*for square := uint8(bits.TrailingZeros64(allPieces)); square < 64-uint8(bits.LeadingZeros64(allPieces)); square++ {
		if (uint64(1)<<square)&allPieces != 0 {
			piece, isWhite := utils.GetPiece(square, board)
			//Material & PST Evaluation
			var b *chess.Bitboards
			var color int = 1
			if isWhite {
				b = &board.White
				opening += MaterialScore[OPENING][piece]
				endgame += MaterialScore[ENDGAME][piece]
				opening += PST[OPENING][piece][ReversedBoard[square]]
				endgame += PST[ENDGAME][piece][ReversedBoard[square]]
			} else {
				color = -1
				b = &board.Black
				opening -= MaterialScore[OPENING][piece]
				endgame -= MaterialScore[ENDGAME][piece]
				opening -= PST[OPENING][piece][square]
				endgame -= PST[ENDGAME][piece][square]
			}
			switch piece {
			case chess.Pawn:
				val := DoublePawns(b.Pawns, square)
				val += IsolatedPawns(b.Pawns, square)
				opening -= val * color
				endgame -= val * color
			case chess.Knight:
				mobility := MobilityKnight(square, allPieces)
				opening += mobility * color
				endgame += mobility * color
			case chess.Bishop:
				mobility := MobilityBishop(square, allPieces, b.All)
				opening += mobility * color
				endgame += mobility * color
			case chess.Rook:
				mobility := MobilityRook(square, allPieces, b.All)
				opening += mobility * color
				endgame += mobility * color
			}
		}
	}*/
	if isEndgame(board) {
		if IsDraw(board) {
			return 0
		}
		return endgame * turn // Endgame
	}
	//opening += Mobility(board)
	phase = ((phase * 256) + (TotalPhase / 2)) / TotalPhase
	score := ((opening * (256 - phase)) + (opening * phase)) / 256
	return score * turn
}

func Quiesce(board *chess.Board, alpha int, beta int, turn int) int {
	if isTimeToStop() {
		return 0
	}
	standPat := Evaluate(board, turn)
	if standPat > beta {
		return beta
	}
	// Delta pruning
	/*if standPat < alpha-Delta {
		return alpha
	}*/
	alpha = utils.Max(alpha, standPat)
	moves := filterCaptures(board.GenerateLegalMoves(), board)
	var score int = 0
	for _, move := range moves {
		if isTimeToStop() {
			return 0
		}
		unmakeFunc := Make(board, move)
		score = -Quiesce(board, -beta, -alpha, -turn)
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(move)
			alpha = score
		}
		if score >= beta {
			return beta
		}
	}
	return alpha
}

func filterCaptures(moves []chess.Move, board *chess.Board) []chess.Move {
	var captures []chess.Move
	for _, move := range moves {
		if chess.IsCapture(move, board) {
			captures = append(captures, move)
		}
	}
	var n = len(captures)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if moveOrdering.GetMVV_LVA(captures[j], board) > moveOrdering.GetMVV_LVA(captures[minIdx], board) {
				minIdx = j
			}
		}
		captures[i], captures[minIdx] = captures[minIdx], captures[i]
	}
	/*sort.Slice(filteredCaptures, func(a, b int) bool {
		valueA := moveOrdering.GetMVV_LVA(filteredCaptures[a], board)
		valueB := moveOrdering.GetMVV_LVA(filteredCaptures[b], board)
		return valueA > valueB
	})*/
	return captures
}

/*func getMinorPieces(board *chess.Board) int {
	bin := board.White.Knights | board.White.Bishops | board.White.Rooks |
		board.Black.Knights | board.Black.Bishops | board.Black.Rooks
	return bits.OnesCount64(bin)
}*/

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
