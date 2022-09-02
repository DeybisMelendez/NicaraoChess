package engine

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// https://www.youtube.com/watch?v=QhFtquEeffA
var RepetitionTable = [1000]uint64{}

func IsRepetition(hash uint64) bool {
	for i := 0; i < GamePly; i++ {
		if hash == RepetitionTable[i] {
			return true
		}
	}
	return false
}

func Make(board *chess.Board, move chess.Move) func() {
	Ply++
	Nodes++
	GamePly++
	RepetitionTable[GamePly] = board.Hash()
	f := board.Apply(move)

	return f
}

func Unmake(f func()) {
	Ply--
	GamePly--
	f()
}
