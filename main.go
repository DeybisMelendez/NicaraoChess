package main

import (
	"fmt"
	"nicarao/search"
	"nicarao/uci"

	chess "github.com/dylhunn/dragontoothmg"
)

func main() {
	uci.Init()
	//Ddebug()
	uci.UCI()
}

func Debug() {

	var board = chess.ParseFen("8/2pp1p1P/4k2P/8/8/p3K3/p1PP1P2/8 w - - 0 1")
	/*fmt.Println("peon pasado negro", search.PassedPawns(board.White.Pawns, 8, false))
	fmt.Println("peon doblado negro", search.DoublePawns(board.Black.Pawns, 8))
	fmt.Println("peon aislado negro", search.IsolatedPawns(board.Black.Pawns, 53))
	fmt.Println("peon pasado blanco", search.PassedPawns(board.Black.Pawns, 55, true))
	fmt.Println("peon doblado blanco", search.DoublePawns(board.White.Pawns, 55))
	fmt.Println("peon aislado blanco", search.IsolatedPawns(board.White.Pawns, 13))*/
	fmt.Println("peon pasado negro", search.PassedPawns(board.White.Pawns, 53, false))

	fmt.Println("peon doblado negro", search.DoublePawns(board.Black.Pawns, 51))
	fmt.Println("peon aislado negro", search.IsolatedPawns(board.Black.Pawns, 50))

	fmt.Println("peon pasado blanco", search.PassedPawns(board.Black.Pawns, 13, true))

	fmt.Println("peon doblado blanco", search.DoublePawns(board.White.Pawns, 11))
	fmt.Println("peon aislado blanco", search.IsolatedPawns(board.White.Pawns, 10))
}
