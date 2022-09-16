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
				if newVal > val { // && SEE(board, moveList[i], moveList[i].To(), 0, 1) >= 0 {
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
		//if SEE(board, move, move.To(), 0, 1) >= 0 {
		unmakeFunc := Make(board, move)
		score = -Quiesce(board, -beta, -alpha, -turn)
		Unmake(unmakeFunc)
		if alpha >= beta {
			return beta
		}
		if score > alpha {
			StorePV(move)
			alpha = score
		}
		//}
	}
	return alpha
}
