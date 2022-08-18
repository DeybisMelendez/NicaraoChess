package search

import (
	"nicarao/moveOrdering"
	"nicarao/utils"
	"sort"

	chess "github.com/dylhunn/dragontoothmg"
)

var MaterialWeightOP = [6]int{100, 300, 300, 500, 1000, 10000}
var MaterialWeightMG = [6]int{100, 320, 330, 500, 900, 10000}
var MaterialWeightEG = [6]int{120, 320, 350, 550, 900, 10000}
var MaterialScore = [3][6]int{MaterialWeightOP, MaterialWeightMG, MaterialWeightEG}

const Delta = 200
const MiddleGamePhaseScore = 7300 //16xP + 4xB + 4xN + 4xR + 2xQ - R - 2xP
// const EndGamePhaseScore = 1800    // 2Q
var phase = 1 // 0 opening, 1 middlegame, 2 endgame

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
	piecesCount := getPiecesCount(board)
	queens := utils.NumOfSetBits(board.White.Queens) + utils.NumOfSetBits(board.Black.Queens)
	if piecesCount > 28 {
		phase = 0 //Opening
	} else if (piecesCount < 16 && queens == 0) || queens < 2 {
		phase = 2 // Endgame
	}
	score := 0
	for i := uint8(0); i < 64; i++ {
		piece, isWhite := utils.GetPiece(i, board)
		if piece != chess.Nothing {
			color := 1
			material, pst := 0, 0
			if isWhite {
				color = 0
			}
			material = MaterialScore[phase][piece-1]
			pst = PST[phase][color][piece-1][i]
			if isWhite {
				score += material + pst
			} else {
				score -= (material + pst)
			}
		}
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

func getPiecesCount(board *chess.Board) uint64 {
	return utils.NumOfSetBits(board.White.All) + utils.NumOfSetBits(board.Black.All)
}
