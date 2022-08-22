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

var MaterialOpening = [7]int{0, 100, 320, 330, 500, 900, 10000}
var MaterialEndgame = [7]int{0, 120, 320, 350, 550, 900, 10000}
var MaterialScore = [2][7]int{MaterialOpening, MaterialEndgame} //Opening, Endgame

var KnightPhase int = 1
var BishopPhase int = 1
var RookPhase int = 2
var QueenPhase int = 4

const Delta = 200

func Evaluate(board *chess.Board, turn int) int {
	moves := board.GenerateLegalMoves()
	if len(moves) == 0 {
		if board.OurKingInCheck() {
			//Checkmate
			return -MateScore + Ply
		} else {
			//Stalemate
			return 0
		}
	}
	if IsThreeFoldRepetition(board.Hash()) {
		return 0
	}
	opening, endgame := 0, 0

	whitePawnCount := int(bits.OnesCount64(board.White.Pawns))
	blackPawnCount := int(bits.OnesCount64(board.White.Pawns))
	allPieces := board.White.All | board.Black.All
	allPawnCount := whitePawnCount + blackPawnCount

	var TotalPhase int = KnightPhase*4 +
		BishopPhase*4 + RookPhase*4 + QueenPhase*2 //+PawnPhase*16
	var phase int = TotalPhase

	for square := uint8(0); square < 64; square++ {
		piece, isWhite := utils.GetPiece(square, board)
		if piece != chess.Nothing {
			//Material & PST Evaluation
			if isWhite {
				opening += MaterialScore[OPENING][piece]
				endgame += MaterialScore[ENDGAME][piece]

				opening += PST[WHITE][OPENING][piece][square]
				endgame += PST[WHITE][ENDGAME][piece][square]
			} else {
				opening -= MaterialScore[OPENING][piece]
				endgame -= MaterialScore[ENDGAME][piece]
				opening -= PST[BLACK][OPENING][piece][square]
				endgame -= PST[BLACK][ENDGAME][piece][square]
			}
			//Piece Evaluation
			switch piece {
			case chess.Pawn:
				if isWhite {
					opening += CenterPawn(square, false)
					/*opening -= DoublePawns(board.White.Pawns, square)
					opening -= IsolatedPawns(board.White.Pawns, square)
					opening += PassedPawns(board.White.Pawns, square, true)

					endgame -= DoublePawns(board.White.Pawns, square)
					endgame -= IsolatedPawns(board.White.Pawns, square)
					endgame += PassedPawns(board.White.Pawns, square, true)*/
				} else {
					opening -= CenterPawn(square, false)
					/*opening += DoublePawns(board.Black.Pawns, square)
					opening += IsolatedPawns(board.Black.Pawns, square)
					opening -= PassedPawns(board.Black.Pawns, square, false)

					endgame += DoublePawns(board.Black.Pawns, square)
					endgame += IsolatedPawns(board.Black.Pawns, square)
					endgame -= PassedPawns(board.Black.Pawns, square, false)*/
				}
			case chess.Knight:
				phase -= KnightPhase
				val := BadKnight(allPawnCount)
				if isWhite {
					opening += val
					endgame += val
				} else {
					opening -= val
					endgame -= val
				}
			case chess.Bishop:
				phase -= BishopPhase
				if isWhite {
					val := BishopPair(board.White.Bishops)
					val += MobilityBishop(square, allPieces, board.White.All)
					val -= BadBishop(square, board.White.Pawns)
					opening += val
					endgame += val
				} else {
					val := BishopPair(board.Black.Bishops)
					val += MobilityBishop(square, allPieces, board.Black.All)
					val -= BadBishop(square, board.Black.Pawns)
					opening -= val
					endgame -= val
				}
			case chess.Rook:
				phase -= RookPhase
				if isWhite {
					val := GoodRook(allPawnCount)
					val += MobilityRook(square, allPieces, board.White.All)
					opening += val
					endgame += val
				} else {
					val := GoodRook(allPawnCount)
					val += MobilityRook(square, allPieces, board.Black.All)
					opening -= val
					endgame -= val
				}
			case chess.Queen:
				phase -= QueenPhase
				if isWhite {
					val := BadQueen(board, true, square)
					opening -= val
					endgame -= val
				} else {
					val := BadQueen(board, false, square)
					opening += val
					endgame += val
				}
			case chess.King:
				if isWhite {
					opening -= BadKing(square, allPieces, board.White.All, false)
				} else {
					opening += BadKing(square, allPieces, board.Black.All, false)
				}
			}
		}
	}
	phase = ((phase * 256) + (TotalPhase / 2)) / TotalPhase
	score := ((opening * (256 - phase)) + (opening * phase)) / 256
	if isEndgame(board) {
		return endgame * turn // Endgame
	}
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
	if standPat < alpha-Delta {
		return alpha
	}
	alpha = utils.Max(alpha, standPat)
	moves := filterCaptures(board.GenerateLegalMoves(), board)
	var score int = 0
	for i := 0; i < len(moves); i++ {
		if isTimeToStop() {
			return 0
		}
		unmakeFunc := Make(board, moves[i])
		score = -Quiesce(board, -beta, -alpha, -turn)
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(moves[i])
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
	for i := 0; i < len(moves); i++ {
		if chess.IsCapture(moves[i], board) {
			filteredCaptures = append(filteredCaptures, moves[i])
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
