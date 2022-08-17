package moveOrdering

import chess "github.com/dylhunn/dragontoothmg"

var KillerMoves [2][64]chess.Move

func StoreKillerMove(move chess.Move, board *chess.Board, ply int) {
	if ply < 64 && !chess.IsCapture(move, board) {
		KillerMoves[1][ply] = KillerMoves[0][ply]
		KillerMoves[0][ply] = move
	}
}

func IsKillerMove(move chess.Move, ply int) bool {
	return KillerMoves[0][ply] == move && KillerMoves[1][ply] == move
}

func ResetKillerMoves() {
	var newKillerMoves [2][64]chess.Move
	KillerMoves = newKillerMoves
	/*for j := 2; j < 2; j++ {
		for i := 0; i < len(KillerMoves); i++ {
			KillerMoves[j][i] = 0
		}
	}*/
}
