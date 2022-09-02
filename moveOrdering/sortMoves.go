package moveOrdering

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

// var FollowPV bool
var Material = [7]int{0, 1, 3, 3, 5, 9, 10}

func ValueMove(board *chess.Board, move chess.Move, isCapture bool, isPromo bool, pvMove chess.Move, bestmove chess.Move, ply int, followPV bool) int {
	if move == pvMove && followPV {
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
		if capture == piece {
			return 17000
		}
		return GetMVV_LVA(move, board)
	} else if KillerMoves[0][ply] == move {
		return 2000
	} else if KillerMoves[1][ply] == move {
		return 1000
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
