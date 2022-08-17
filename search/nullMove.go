package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

const DoNull = true
const NoNull = false
const NullMoveR = 2
const NullMoveFails = 10000

func NullMove(board chess.Board, inCheck bool, depth int, alpha int, beta int) int {
	if isTimeToStop() {
		return 0
	}
	score := 0
	if Ply > 0 && depth >= NullMoveR+1 && !inCheck {
		//color := 1
		// Flip the player to move
		//board.hash ^= whiteToMoveZobristC
		board.Wtomove = !board.Wtomove

		// Restore the halfmove clock
		board.Halfmoveclock = uint8(0)
		/*fmt.Println("before:", fen)
		if strings.Contains(fen, " w ") {
			fen = strings.ReplaceAll(fen, " w ", " b ")
			color = -1
		} else {
			fen = strings.ReplaceAll(fen, " b ", " w ")
		}
		fmt.Println("after:", fen)
		board := chess.ParseFen(fen)*/
		score = -Negascout(&board, depth-1-NullMoveR, -beta, -beta+1)
		if score >= beta {
			return beta
		}
	}
	return NullMoveFails // Null
}
