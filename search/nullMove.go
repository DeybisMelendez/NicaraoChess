package search

import (
	"strings"

	chess "github.com/dylhunn/dragontoothmg"
)

const DoNull = true
const NoNull = false
const NullMoveR = 2
const NullMoveFails = 10000

func NullMove(fen string, depth int, beta int) int {
	if isTimeToStop() {
		return 0
	}
	score := 0
	if Ply > 0 && depth > NullMoveR {
		if strings.Contains(fen, " w ") {
			fen = strings.ReplaceAll(fen, " w ", " b ")
		} else {
			fen = strings.ReplaceAll(fen, " b ", " w ")
		}
		nullBoard := chess.ParseFen(fen)
		if !nullBoard.OurKingInCheck() && len(nullBoard.GenerateLegalMoves()) != 0 {
			score = -Negascout(&nullBoard, depth-1-NullMoveR, -beta, -beta+1)
			if score >= beta {
				return beta
			}
		}

	}
	return NullMoveFails // Null
}
