package search

import chess "github.com/dylhunn/dragontoothmg"

const Delta = 200

var Ply int = 0
var Nodes int = 0
var StopTime int64 = -1
var Stopped bool = false

const MateScore = 4000

var MateValue = 5000

var Bestmove chess.Move

func ResetGlobalVariables() {
	Ply = 0
	Nodes = 0
	Stopped = false
	MateValue = 5000
}
