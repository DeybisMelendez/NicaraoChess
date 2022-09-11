package engine

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

func Quiesce(board *chess.Board, alpha int, beta int, turn int) int {
	if isTimeToStop() {
		return 0
	}
	PVLength[Ply] = Ply
	standPat := Evaluate(board, turn)
	if standPat > beta {
		return beta
	}
	// Delta pruning
	/*if standPat < alpha-Delta {
		return alpha
	}*/
	alpha = utils.Max(alpha, standPat)
	moveList := board.GenerateLegalMoves()
	bestmove := checkPV(moveList)
	var score int = 0
	for {
		var val int = -1
		var idx int = -1
		var ln int = len(moveList)
		for i := 0; i < ln; i++ {
			if moveList[i] == bestmove && FollowPV {
				idx = i
				break
			}
			if chess.IsCapture(moveList[i], board) {
				var newVal int = GetMVV_LVA(moveList[i], board)
				if newVal > val {
					val = newVal
					idx = i
				}
			}
		}
		if idx == -1 {
			break
		}
		var move chess.Move = moveList[idx]
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		if SEE(board, move, move.To(), 0, 1) >= 0 {
			unmakeFunc := Make(board, move)
			score = -Quiesce(board, -beta, -alpha, -turn)
			Unmake(unmakeFunc)
			if score > alpha {
				StorePV(move)
				alpha = score
				if score >= beta {
					return beta
				}
			}
		} /* else {
			i := -1
			for i < len(moveList) {
				i++
				if moveList[i].To() == move.To() {
					moveList = append(moveList[:i], moveList[i+1:]...)
					i = -1
				}
			}
		}*/
	}
	return alpha
}
