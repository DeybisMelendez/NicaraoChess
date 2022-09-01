package moveOrdering

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

// var FollowPV bool
var Material = [7]int{0, 1, 3, 3, 5, 9, 10}

func ValueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, ply int) int {
	if move == pvMove {
		return 6000
	} else if move == bestmove {
		return 5000
	} else if chess.IsCapture(move, board) || move.Promote() != chess.Nothing {
		piece, _ := utils.GetPiece(move.From(), board)
		capture, _ := utils.GetPiece(move.To(), board)
		promo := move.Promote()
		if Material[capture] > Material[piece] || promo != chess.Nothing {
			return 4500 + Material[capture] - Material[piece] + Material[promo]
		}
		if capture == piece {
			return 4000
		}
		return GetMVV_LVA(move, board)
	} else if KillerMoves[0][ply] == move {
		return 3000
	} else if KillerMoves[1][ply] == move {
		return 2000
	} else {
		return GetHistoryMove(board.Wtomove, move) + 1000
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
