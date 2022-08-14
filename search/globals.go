package search

import (
	chess "github.com/dylhunn/dragontoothmg"
)

const MaxPly = 64

var Ply int = 0
var Nodes int = 0

//var Mate int = 4000

func ResetGlobalVariables() {
	Ply = 0
	Nodes = 0
}

func GetPieces(move chess.Move, board *chess.Board) {
	println(board.White.All)
}
