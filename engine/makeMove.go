package engine

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// https://www.youtube.com/watch?v=QhFtquEeffA
var RepetitionTable = [256]uint64{}

func IsRepetition(hash uint64) bool {
	count := 0
	for i := 0; i < GamePly; i++ {
		if hash == RepetitionTable[i] {
			count++
		}
		if count == 2 {
			return true
		}
	}
	return false
}

func ResetRepetitionTable() {
	var newRep = [256]uint64{}
	RepetitionTable = newRep
	GamePly = 0
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
