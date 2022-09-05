package engine

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

func Quiesce(board *chess.Board, alpha int, beta int, turn int) int {
	PVLength[Ply] = Ply
	if isTimeToStop() {
		return 0
	}
	/*if Ply >= 64 {
		return Evaluate(board, turn)
	}*/
	standPat := Evaluate(board, turn)
	if standPat > beta {
		return beta
	}
	// Delta pruning
	/*if standPat < alpha-Delta {
		return alpha
	}*/
	alpha = utils.Max(alpha, standPat)
	moveList := board.GenerateLegalMoves() //captures(board.GenerateLegalMoves(), board)
	checkPV(moveList)
	var score int = 0
	for {
		var val int = -1
		var idx int = -1
		var ln int = len(moveList)
		for i := 0; i < ln; i++ {
			//var newVal int = ValueMove(board, moveList[i], true, moveList[i].Promote() != chess.Nothing, PVTable[Ply][Ply], 0)
			// + MaterialOpening[moveList[i].Promote()]
			if chess.IsCapture(moveList[i], board) {
				if moveList[i] == PVTable[0][Ply] && FollowPV {
					idx = i
					break
				}
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
			}
			if score >= beta {
				return beta
			}
		} else {
			break
		}
	}
	return alpha
}

/*func captures(moveList []chess.Move, board *chess.Board) []chess.Move {
	var captures []chess.Move
	for _, move := range moveList {
		if chess.IsCapture(move, board) {
			captures = append(captures, move)
		}
	}
	return captures
}*/
