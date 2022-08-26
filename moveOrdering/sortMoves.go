package moveOrdering

import (
	"math/rand"
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

var FollowPV bool

func valueMove(board *chess.Board, move chess.Move, pvMove chess.Move, bestmove chess.Move, isWhite bool, ply int) int {
	if move == bestmove {
		return 5000
	} else if move == pvMove {
		return 4000
	} else if chess.IsCapture(move, board) || move.Promote() != chess.Nothing {
		return GetMVV_LVA(move, board) + 3000
	} else if KillerMoves[0][ply] == move {
		return 2000
	} else if KillerMoves[1][ply] == move {
		return 1000
	} else {
		piece, _ := utils.GetPiece(move.From(), board)
		return GetHistoryMove(isWhite, piece, move.To()) + GetMVV_LVA(move, board)
	}
}

func SortMoves(moves []chess.Move, board *chess.Board, pvMove chess.Move, bestmove chess.Move, ply int) {
	var list = make([]int, len(moves))
	for i, m := range moves {
		list[i] = valueMove(board, m, pvMove, bestmove, board.Wtomove, ply)
	}
	Quicksort(list, moves)
}
func Quicksort(a []int, moves []chess.Move) {
	if len(a) < 2 {
		return
	}
	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]
	moves[pivot], moves[right] = moves[right], moves[pivot]
	for i := range a {
		if a[i] > a[right] {
			a[left], a[i] = a[i], a[left]
			moves[left], moves[i] = moves[i], moves[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]
	moves[left], moves[right] = moves[right], moves[left]
	Quicksort(a[:left], moves[:left])
	Quicksort(a[left+1:], moves[left+1:])

	//return a, moves
}

/*func SelectionSort(items []chess.Move) {
	var n = len(items)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if items[j] < items[minIdx] {
				minIdx = j
			}
		}
		items[i], items[minIdx] = items[minIdx], items[i]
	}
}*/
