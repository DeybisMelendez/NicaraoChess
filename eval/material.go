package eval

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

var MaterialWeightMG = [7]int{0, 100, 320, 330, 500, 900, 10000}
var MaterialWeightEG = [7]int{0, 120, 320, 350, 550, 900}

var MaterialMG int = 0
var lastMaterial int = 0

func GetMaterial(board *chess.Board) int {
	var white int = 0
	var black int = 0
	for i := 0; i < 63; i++ {
		piece, isWhite := utils.GetPiece(uint8(i), board)
		if isWhite {
			white += MaterialWeightMG[piece]
		} else {
			black += MaterialWeightMG[piece]
		}
	}
	return white - black
}

func UpdateMaterial(board *chess.Board, move chess.Move) {
	capture, isWhite := utils.GetPiece(move.To(), board)
	if chess.IsCapture(move, board) {
		MaterialMG = lastMaterial
		if isWhite {
			MaterialMG -= MaterialWeightMG[capture]
		} else {
			MaterialMG += MaterialWeightMG[capture]
		}
	}
	promotion := move.Promote()
	if promotion != chess.Nothing {
		if isWhite {
			MaterialMG += MaterialWeightMG[promotion]
			MaterialMG -= MaterialWeightMG[chess.Pawn]
		} else {
			MaterialMG -= MaterialWeightMG[promotion]
			MaterialMG += MaterialWeightMG[chess.Pawn]
		}
	}
}

func RevertMaterial() {
	MaterialMG = lastMaterial
}

func ResetMaterial() {
	MaterialMG = 0
	lastMaterial = 0
}
