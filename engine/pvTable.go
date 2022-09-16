package engine

import (
	chess "github.com/dylhunn/dragontoothmg"
)

var PVLength [128]int
var PVTable [128][128]chess.Move
var FollowPV bool

//var ScorePV bool

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
	var newPVLength [128]int
	var newPVTable [128][128]chess.Move
	PVLength = newPVLength
	PVTable = newPVTable
}

func FormatPV() string {
	str := ""
	for i := 0; i < PVLength[0]; i++ {
		str += PVTable[0][i].String() + " "
	}
	return str
}

func checkPV(moveList []chess.Move) chess.Move {
	if FollowPV {
		FollowPV = false
		for _, move := range moveList {
			if move == PVTable[0][Ply] {
				FollowPV = true
				return move
			}
		}
	}
	return 0
}
