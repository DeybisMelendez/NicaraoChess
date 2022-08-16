package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// posiblemente tenga que completar el pvmove si no funciona correctamente, falta followPV
var PVLength [MaxPly]int
var PVTable [MaxPly][MaxPly]chess.Move

func StorePV(move chess.Move) {
	// Triangular PV Table
	// escribe el actual pv
	if Ply < 63 {
		PVTable[Ply][Ply] = move
		// escribimos desde la capa mas profunda hasta la actual
		for nextPly := Ply + 1; nextPly < PVLength[Ply+1]; nextPly++ {
			PVTable[Ply][nextPly] = PVTable[Ply+1][nextPly]
		}
		// ajuste pv length
		PVLength[Ply] = PVLength[Ply+1]
	}

}

func ResetPVTable() {
	var newPVLength [MaxPly]int
	var newPVTable [MaxPly][MaxPly]chess.Move
	PVLength = newPVLength
	PVTable = newPVTable
	/*for i := 0; i < len(PVLength); i++ {
		PVLength[i] = 0
	}
	for j := 0; j < len(PVTable); j++ {
		for i := 0; i < len(PVTable[j]); i++ {
			PVTable[j][i] = 0
		}
	}*/
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
