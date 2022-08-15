package search

import "nicarao/eval"

var Ply int = 0
var Nodes int = 0

const MaxPly = 64

//var Mate int = 4000

func ResetGlobalVariables() {
	Ply = 0
	Nodes = 0
	eval.ResetMaterial()
}
