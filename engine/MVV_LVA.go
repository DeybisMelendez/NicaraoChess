package engine

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

//Most Valuable Victim - Least Valuable Aggressor

var MVV_LVA = [7][7]int{ //[Victims][Agressors]
	//Agressors {Nothing, Pawn, Knight, Bishop, Rook, Queen, King}
	//Victim [Nothing, Pawn, Knight, Bishop, Rook, Queen, King]
	{ //Nothing
		0, 50, 40, 30, 20, 10, 0,
	},
	{ //Pawn
		0, 150, 140, 130, 120, 110, 100,
	},
	{ //Knight
		0, 250, 240, 230, 220, 210, 200,
	},
	{ //Bishop
		0, 350, 340, 330, 320, 310, 300,
	},
	{ //Rook
		0, 450, 440, 430, 420, 410, 400,
	},
	{ //Queen
		0, 550, 540, 530, 520, 510, 500,
	},
	{ //King
		0, 650, 640, 630, 620, 610, 600,
	},
}

func GetMVV_LVA(move chess.Move, board *chess.Board) int {
	agressor, _ := utils.GetPiece(move.From(), board)
	victim, _ := utils.GetPiece(move.To(), board)
	return MVV_LVA[victim][agressor]
}
