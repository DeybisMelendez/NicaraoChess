package engine

import chess "github.com/dylhunn/dragontoothmg"

var KillerMoves [2][64]chess.Move

func StoreKillerMove(move chess.Move) {
	if Ply < 64 {
		KillerMoves[1][Ply] = KillerMoves[0][Ply]
		KillerMoves[0][Ply] = move
	}
}

func IsKillerMove(move chess.Move) bool {
	return KillerMoves[0][Ply] == move || KillerMoves[1][Ply] == move
}

func ResetKillerMoves() {
	var newKillerMoves [2][64]chess.Move
	KillerMoves = newKillerMoves
}
