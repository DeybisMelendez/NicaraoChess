package engine

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

func SEE(board *chess.Board, move chess.Move, square uint8, value int, turn int) int {
	piece, _ := utils.GetPiece(square, board)
	value += Material[piece] * turn
	unmake := Make(board, move)
	moveList := board.GenerateLegalMoves()
	for _, next := range moveList {
		if next.To() == square {
			if chess.IsCapture(next, board) {
				bestvalue := SEE(board, next, square, value, -turn)
				value = utils.Max(bestvalue, value)
			} /* else {
				break
			}*/
		}
	}
	Unmake(unmake)
	return value
}
