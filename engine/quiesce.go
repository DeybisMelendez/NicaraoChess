package engine

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

// DELTA_PRUNING for Quiesce
const DELTA_PRUNING = 20

func quiesce(board *chess.Board, alpha int16, beta int16, turn int16) int16 {
	if isTimeToStop() {
		return 0
	}
	var standPat int16 = Evaluate(board, turn)
	// Standing Pat
	if standPat >= beta {
		return beta
	}
	alpha = utils.MaxInt16(alpha, standPat)
	// Delta pruning
	/*if standPat+DELTA_PRUNING <= alpha {
		return alpha
	}*/
	moveList := board.GenerateLegalMoves()
	if len(moveList) == 0 {
		return standPat
	}
	for {
		var maxScoreMove int16
		var hasCapture bool
		var idx int
		for i, candidate := range moveList {
			var scoreMove int16
			if chess.IsCapture(candidate, board) {
				hasCapture = true
				scoreMove = getMVV_LVA(candidate, board)
				if scoreMove > maxScoreMove {
					maxScoreMove = scoreMove
					idx = i
				}
			}
		}
		if !hasCapture {
			break
		}
		var move chess.Move = moveList[idx]
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		unmakeMoveFunc := makeMove(board, move)
		eval := -quiesce(board, -beta, -alpha, -turn)
		unmakeMove(unmakeMoveFunc)
		alpha = utils.MaxInt16(alpha, eval)
		if eval >= beta {
			return beta
		}

		/*if eval >= beta-DELTA_PRUNING {
			return beta
		}*/
	}
	return alpha
}
