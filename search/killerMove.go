package search

import chess "github.com/dylhunn/dragontoothmg"

var KillerMoves [2][MaxPly]chess.Move

func StoreKillerMove(move chess.Move, board *chess.Board) {
	if Ply < MaxPly && chess.IsCapture(move, board) {
		KillerMoves[1][Ply] = KillerMoves[0][Ply]
		KillerMoves[0][Ply] = move
	}
}

func ResetKillerMoves() {
	var newKillerMoves [2][MaxPly]chess.Move
	KillerMoves = newKillerMoves
}
