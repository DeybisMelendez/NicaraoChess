package search

import (
	"strings"

	chess "github.com/dylhunn/dragontoothmg"
)

const DoNull = true
const NoNull = false
const NullDepth = 3
const NullDivisor = 6
const NullMoveFails = 10000

// R = null_depth + depth / null_divisor
func NullMove(board *chess.Board, depth int, alpha int, beta int, turn int) int {
	score := 0
	var fen = board.ToFen()
	if strings.Contains(fen, " w ") {
		fen = strings.ReplaceAll(fen, " w ", " b ")
	} else {
		fen = strings.ReplaceAll(fen, " b ", " w ")
	}
	nullBoard := chess.ParseFen(fen)
	if !nullBoard.OurKingInCheck() && len(nullBoard.GenerateLegalMoves()) != 0 {
		var R = NullDepth + depth/NullDivisor
		if depth-R-1 > 0 {
			score = -Negamax(&nullBoard, depth-1-NullDepth, -beta, -beta+1, -turn, NoNull)
		} else {
			score = -Quiesce(&nullBoard, -beta, -alpha, -turn)
		}
		if score >= beta {
			return beta
		}
	}
	return NullMoveFails // Null
}

func AllowNullMove(board *chess.Board) bool {
	if board.Wtomove {
		if (board.White.Knights&board.White.Bishops)&(board.White.Rooks&board.White.Queens) > 0 {
			return true
		}
	} else {
		if (board.Black.Knights&board.Black.Bishops)&(board.Black.Rooks&board.Black.Queens) > 0 {
			return true
		}
	}
	return false
}
