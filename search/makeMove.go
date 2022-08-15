package search

import (
	"nicarao/eval"

	chess "github.com/dylhunn/dragontoothmg"
)

func Make(board *chess.Board, move chess.Move) func() {
	Ply++
	Nodes++
	eval.UpdateMaterial(board, move)
	return board.Apply(move)
}

func Unmake(f func()) {
	Ply--
	eval.RevertMaterial()
	f()
}
