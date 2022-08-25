package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// https://www.youtube.com/watch?v=QhFtquEeffA
var RepetitionTable = [150]uint64{}

func IsRepetition(hash uint64) bool {
	for i := 0; i < Ply; i++ {
		if hash == RepetitionTable[i] {
			return true
		}
	}
	return false
}

func Make(board *chess.Board, move chess.Move) func() {
	RepetitionTable[Ply] = board.Hash()
	Ply++
	Nodes++
	f := board.Apply(move)

	return f
}

func Unmake(f func()) {
	Ply--
	f()
}
