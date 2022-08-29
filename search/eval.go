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

const Delta = 200

var TotalPhase int = KnightPhase*4 +
	BishopPhase*4 + RookPhase*4 + QueenPhase*2 //+PawnPhase*16
var phase int = TotalPhase

func Evaluate(board *chess.Board, turn int) int {
	opening, endgame := 0, 0
	whitePawnCount := bits.OnesCount64(board.White.Pawns)
	blackPawnCount := bits.OnesCount64(board.White.Pawns)
	allPieces := board.White.All | board.Black.All
	allPawnCount := whitePawnCount + blackPawnCount
	opening += bits.OnesCount64(board.White.Pawns&Center) * 20
	opening -= bits.OnesCount64(board.Black.Pawns&Center) * 20
	opening += bits.OnesCount64(board.White.Pawns&ExtendedCenter) * 10
	opening -= bits.OnesCount64(board.Black.Pawns&ExtendedCenter) * 10
	all := bits.OnesCount64(allPieces)
	var count int
	for square := uint8(0); square < 64; square++ {
		if count >= all {
			break
		}
		if (uint64(1)<<square)&allPieces != 0 {
			count++
			piece, isWhite := utils.GetPiece(square, board)
			//Material & PST Evaluation
			if isWhite {
				opening += MaterialScore[OPENING][piece]
				endgame += MaterialScore[ENDGAME][piece]
				opening += PST[OPENING][piece][ReversedBoard[square]]
				endgame += PST[ENDGAME][piece][ReversedBoard[square]]
				switch piece {
				case chess.Pawn:
					val := PassedPawns(board.White.Pawns, square, true)
					val -= DoublePawns(board.White.Pawns, square)
					val -= IsolatedPawns(board.White.Pawns, square)
					//opening += CenterPawn(square)
					opening += val
					endgame += val
				case chess.Knight:
					phase -= KnightPhase
					val := GoodKnight(allPawnCount)
					val += MobilityKnight(square, allPieces)
					opening += val
					endgame += val
				case chess.Bishop:
					phase -= BishopPhase
					val := BishopPair(board.White.Bishops)
					val += MobilityBishop(square, allPieces, board.White.All)
					val -= BadBishop(square, board.White.Pawns)
					opening += val
					endgame += val

				case chess.Rook:
					phase -= RookPhase
					val := GoodRook(allPawnCount)
					val += MobilityRook(square, allPieces, board.White.All)
					opening += val
					endgame += val

				case chess.Queen:
					phase -= QueenPhase
					val := BadQueen(board, true, square)
					opening -= val
					endgame -= val

				case chess.King:
					opening -= BadKing(square, allPieces, board.White.All, false, board.Black.Queens)
					endgame -= BadKing(square, allPieces, board.Black.All, true, board.Black.Queens)
					opening -= AttackedKing(board.OurKingInCheck(), board.Wtomove, true)
					endgame -= AttackedKing(board.OurKingInCheck(), board.Wtomove, false)
				}
			} else {
				opening -= MaterialScore[OPENING][piece]
				endgame -= MaterialScore[ENDGAME][piece]
				opening -= PST[OPENING][piece][square]
				endgame -= PST[ENDGAME][piece][square]
				//Piece Evaluation
				switch piece {
				case chess.Pawn:
					val := PassedPawns(board.Black.Pawns, square, false)
					val -= IsolatedPawns(board.Black.Pawns, square)
					val -= DoublePawns(board.Black.Pawns, square)
					//opening -= CenterPawn(square)
					opening -= val
					endgame -= val
				case chess.Knight:
					phase -= KnightPhase
					val := GoodKnight(allPawnCount)
					val += MobilityKnight(square, allPieces)
					opening -= val
					endgame -= val
				case chess.Bishop:
					phase -= BishopPhase
					val := BishopPair(board.Black.Bishops)
					val += MobilityBishop(square, allPieces, board.Black.All)
					val -= BadBishop(square, board.Black.Pawns)
					opening -= val
					endgame -= val
				case chess.Rook:
					phase -= RookPhase
					val := GoodRook(allPawnCount)
					val += MobilityRook(square, allPieces, board.Black.All)
					opening -= val
					endgame -= val
				case chess.Queen:
					phase -= QueenPhase
					val := BadQueen(board, false, square)
					opening += val
					endgame += val
				case chess.King:
					opening += BadKing(square, allPieces, board.Black.All, false, board.White.Queens)
					endgame += BadKing(square, allPieces, board.Black.All, true, board.White.Queens)
					opening += AttackedKing(board.OurKingInCheck(), board.Wtomove, true)
					endgame += AttackedKing(board.OurKingInCheck(), board.Wtomove, false)
				}
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
