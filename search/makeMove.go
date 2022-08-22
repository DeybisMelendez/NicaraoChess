package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// https://www.youtube.com/watch?v=QhFtquEeffA
var repetitionTable = []uint64{}

func IsThreeFoldRepetition(hash uint64) bool {
	repetitionCount := 0
	for i := 0; i < len(repetitionTable); i++ {
		if hash == repetitionTable[i] {
			repetitionCount++
		}
		if repetitionCount >= 3 {
			return true
		}
	}
	return false
}

func Make(board *chess.Board, move chess.Move) func() {
	Ply++
	Nodes++
	f := board.Apply(move)
	repetitionTable = append(repetitionTable, board.Hash())
	return f
}

func Unmake(f func()) {
	Ply--
	// pop last element
	if len(repetitionTable) > 0 {
		repetitionTable = repetitionTable[:len(repetitionTable)-1]
	}
	f()
}
