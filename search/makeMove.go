package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// https://www.youtube.com/watch?v=QhFtquEeffA
var RepetitionTable = [150]uint64{}

func IsThreeFoldRepetition(hash uint64) bool {
	var count int = 0
	for i := 0; i < Ply; i++ {
		if hash == RepetitionTable[i] {
			count++
		}
		if count > 1 {
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
