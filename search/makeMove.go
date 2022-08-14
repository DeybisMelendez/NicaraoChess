package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

func Make(board *chess.Board, move chess.Move) func() {
	Ply++
	Nodes++
	//fmt.Println(Ply, Nodes)
	return board.Apply(move)
}

func Unmake(f func()) {
	Ply--
	f()
}
