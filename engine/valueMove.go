package engine

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

// var FollowPV bool

func ValueMove(board *chess.Board, move chess.Move, isCapture bool, isPromo bool, pvMove chess.Move, bestmove chess.Move) int {
	if move == pvMove && FollowPV {
		return 20000
	} else if move == bestmove {
		return 19000
	} else if isCapture || isPromo {
		piece, _ := utils.GetPiece(move.From(), board)
		capture, _ := utils.GetPiece(move.To(), board)
		promo := move.Promote()
		if Material[capture] > Material[piece] {
			return 18000 + Material[capture] - Material[piece] + Material[promo]
		}
		if isPromo {
			return 17000 + Material[promo]
		}
		/*if capture == piece {
			return 3000
		}*/
		return GetMVV_LVA(move, board)
	} else if KillerMoves[0][Ply] == move {
		return 2000 + GetHistoryMove(board.Wtomove, move)
	} else if KillerMoves[1][Ply] == move {
		return 1000 + GetHistoryMove(board.Wtomove, move)
	} else {
		return GetHistoryMove(board.Wtomove, move)
	}
}

/*func SortMoves(moves []chess.Move, board *chess.Board, pvMove chess.Move, bestmove chess.Move, ply int) {
	var n = len(moves)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if valueMove(board, moves[j], pvMove, bestmove, ply) > valueMove(board, moves[minIdx], pvMove, bestmove, ply) {
				minIdx = j
			}
		}
		moves[i], moves[minIdx] = moves[minIdx], moves[i]
	}
}*/
