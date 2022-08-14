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
	PVTable[Ply][Ply] = move
	// escribimos desde la capa mas profunda hasta la actual
	for nextPly := Ply + 1; nextPly < PVLength[Ply+1]; nextPly++ {
		PVTable[Ply][nextPly] = move
	}
	// ajuste pv length
	PVLength[Ply] = PVLength[Ply+1]
}

/*func containsMove(a [MaxPly]dragontoothmg.Move, x dragontoothmg.Move) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
*/
