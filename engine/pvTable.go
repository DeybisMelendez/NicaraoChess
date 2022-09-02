package engine

import (
	chess "github.com/dylhunn/dragontoothmg"
)

var PVLength [64]int
var PVTable [64][64]chess.Move
var FollowPV bool
var ScorePV bool

func StorePV(move chess.Move) {
	// Triangular PV Table
	// escribe el actual pv
	PVTable[Ply][Ply] = move
	// escribimos desde la capa mas profunda hasta la actual
	for nextPly := Ply + 1; nextPly < PVLength[Ply+1]; nextPly++ {
		PVTable[Ply][nextPly] = PVTable[Ply+1][nextPly]
	}
	// ajuste pv length
	PVLength[Ply] = PVLength[Ply+1]
}

func ResetPVTable() {
	var newPVLength [64]int
	var newPVTable [64][64]chess.Move
	PVLength = newPVLength
	PVTable = newPVTable
}

func FormatPV(moves [64]chess.Move) string {
	str := ""
	for i := 0; i < len(moves); i++ {
		if moves[i] != 0 {
			str += moves[i].String() + " "
		}
	}
	return str
}
