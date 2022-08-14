package main

import (
	"nicarao/search"

	chess "github.com/dylhunn/dragontoothmg"
)

const startpos string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
const infinity int = 5000

var depth, color = 10, 1

func main() {
	board := chess.ParseFen(startpos)
	search.InitHasTable()
	search.Negascout(&board, depth, color, -infinity, infinity)
}
