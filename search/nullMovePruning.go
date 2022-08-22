package search

import (
	"strings"

	chess "github.com/dylhunn/dragontoothmg"
)

const DoNull = true
const NoNull = false
const NullMoveR = 2
const NullMoveFails = 10000

func NullMove(fen string, depth int, beta int, turn int) int {
	score := 0
	if strings.Contains(fen, " w ") {
		fen = strings.ReplaceAll(fen, " w ", " b ")
	} else {
		fen = strings.ReplaceAll(fen, " b ", " w ")
	}
	nullBoard := chess.ParseFen(fen)
	if !nullBoard.OurKingInCheck() && len(nullBoard.GenerateLegalMoves()) != 0 {
		score = -Negamax(&nullBoard, depth-1-NullMoveR, -beta, -beta+1, -turn, NoNull)
		if score >= beta {
			return beta
		}
	}
	return NullMoveFails // Null
}
