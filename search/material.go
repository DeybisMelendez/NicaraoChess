package search

var MaterialWeightOP = [6]int{100, 300, 300, 500, 1000, 10000}
var MaterialWeightMG = [6]int{100, 320, 330, 500, 900, 10000}
var MaterialWeightEG = [6]int{120, 320, 350, 550, 900, 10000}
var MaterialScore = [3][6]int{MaterialWeightOP, MaterialWeightMG, MaterialWeightEG}

//var MaterialMG int = 0
//var lastMaterial int = 0

/*func ValueMaterial(board *chess.Board) (int, int) {
	var whiteMG, blackMG, whiteEG, blackEG int
	for i := 0; i < 64; i++ {
		piece, isWhite := utils.GetPiece(uint8(i), board)
		if isWhite {
			whiteMG += MaterialWeightMG[piece]
			whiteEG += MaterialWeightEG[piece]
		} else {
			blackMG += MaterialWeightMG[piece]
			blackEG += MaterialWeightEG[piece]
		}
	}
	return (whiteMG - blackMG), (whiteEG - blackEG)
}*/

/*func UpdateMaterial(board *chess.Board, move chess.Move) {
	capture, isWhite := utils.GetPiece(move.To(), board)
	MaterialMG = lastMaterial
	if chess.IsCapture(move, board) {
		fmt.Println("iscapture,", MaterialWeightMG[capture], isWhite)
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
*/
