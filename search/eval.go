package search

import (
	"math/bits"
	"nicarao/moveOrdering"
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

const WHITE, BLACK = 0, 1
const OPENING, ENDGAME = 0, 1

var KnightPhase int = 1
var BishopPhase int = 1
var RookPhase int = 2
var QueenPhase int = 4

const Delta = 180

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
	whiteBishopPair := BishopPair(board.White.Bishops)
	blackBishopPair := BishopPair(board.Black.Bishops)
	opening += whiteBishopPair - blackBishopPair

	//Good Knight
	whiteGoodKnight := bits.OnesCount64(board.White.Knights) * GoodKnight(allPawnCount)
	blackGoodKnight := bits.OnesCount64(board.Black.Knights) * GoodKnight(allPawnCount)
	opening += whiteGoodKnight - blackGoodKnight

	//Good Rook
	whiteGoodRook := bits.OnesCount64(board.White.Rooks) * GoodRook(allPawnCount)
	blackGoodRook := bits.OnesCount64(board.White.Rooks) * GoodRook(allPawnCount)
	opening += whiteGoodRook - blackGoodRook

	opening = endgame

	//Pawn structure
	opening += bits.OnesCount64(board.White.Pawns&Center) * 20
	opening -= bits.OnesCount64(board.Black.Pawns&Center) * 20
	opening += bits.OnesCount64(board.White.Pawns&ExtendedCenter) * 10
	opening -= bits.OnesCount64(board.Black.Pawns&ExtendedCenter) * 10
	for square := uint8(bits.TrailingZeros64(allPieces)); square < 64-uint8(bits.LeadingZeros64(allPieces)); square++ {
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
				opening += val * color
				endgame += val * color
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
	}
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
	var filteredCaptures []chess.Move
	for _, move := range moves {
		if chess.IsCapture(move, board) {
			filteredCaptures = append(filteredCaptures, move)
		}
	}
	sort.Slice(filteredCaptures, func(a, b int) bool {
		valueA := moveOrdering.GetMVV_LVA(filteredCaptures[a], board)
		valueB := moveOrdering.GetMVV_LVA(filteredCaptures[b], board)
		return valueA > valueB
	})
	return filteredCaptures
}

func getMinorPieces(board *chess.Board) int {
	bin := board.White.Knights | board.White.Bishops | board.White.Rooks |
		board.Black.Knights | board.Black.Bishops | board.Black.Rooks
	return bits.OnesCount64(bin)
}

func isEndgame(board *chess.Board) bool {
	minorPieces := getMinorPieces(board)
	queens := bits.OnesCount64(board.White.Queens | board.Black.Queens)
	if (minorPieces < 4 && queens < 2) || queens == 0 {
		return true
	}
	return false
}
