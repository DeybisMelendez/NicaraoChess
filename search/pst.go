package search

var PST = pstMake() //[phase][color][piece][square]

func pstMake() [3][2][6][64]int {
	wPawnMG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, -20, -20, 10, 10, 5,
		5, -5, -10, -10, -10, -10, -5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, 5, 10, 20, 20, 10, 5, 5,
		10, 10, 20, 30, 30, 20, 10, 10,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	bPawnMG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 20, 20, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, -10, -10, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	wPawnEG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		-20, -20, -20, -20, -20, -20, -20, -20,
		0, 0, 0, 0, 0, 0, 0, 0,
		10, 10, 10, 10, 10, 10, 10, 5,
		10, 15, 15, 15, 15, 15, 15, 10,
		20, 30, 30, 30, 30, 30, 30, 20,
		40, 60, 60, 60, 60, 60, 60, 40,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	bPawnEG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		40, 60, 60, 60, 60, 60, 60, 40,
		20, 30, 30, 30, 30, 30, 30, 20,
		10, 15, 15, 15, 15, 15, 15, 10,
		10, 10, 10, 10, 10, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
		-20, -20, -20, -20, -20, -20, -20, -20,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	wKnightOG := [64]int{
		0, -10, 0, 0, 0, 0, -10, 0,
		0, 0, 0, 15, 15, 0, 0, 0,
		0, 10, 30, 0, 0, 30, 10, 0,
		0, 0, 0, 10, 10, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	bKnightOG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 10, 10, 0, 0, 0,
		0, 10, 30, 0, 0, 30, 10, 0,
		0, 0, 0, 15, 15, 0, 0, 0,
		0, -10, 0, 0, 0, 0, -10, 0,
	}
	wKnightMG := [64]int{
		-20, -20, -20, -20, -20, -20, -20, -50,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50,
	}
	bKnightMG := [64]int{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-20, -20, -20, -20, -20, -20, -20, -50,
	}
	bBishopOG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 20, 0, 0, 0, 0, 20, 0,
		0, 0, 30, 0, 0, 30, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		10, 15, 0, 10, 10, 0, 15, 10,
		0, 0, -10, 0, 0, -10, 0, 0,
	}
	wBishopOG := [64]int{
		0, 0, -10, 0, 0, -10, 0, 0,
		10, 15, 0, 10, 10, 0, 15, 10,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 30, 0, 0, 30, 0, 0,
		0, 20, 0, 0, 0, 0, 20, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	wBishopMG := [64]int{
		0, -10, -10, -10, -10, -10, -10, 0,
		10, 10, 0, 5, 5, 0, 10, 10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 0, 15, 10, 10, 15, 0, -10,
		-10, 10, 5, 10, 10, 5, 10, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-20, -10, -10, -10, -10, -10, -10, -20,
	}
	bBishopMG := [64]int{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 10, 5, 10, 10, 5, 10, -10,
		-10, 0, 15, 10, 10, 15, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		10, 10, 0, 5, 5, 0, 10, 10,
		0, -10, -10, -10, -10, -10, -10, 0,
	}
	wRookMG := [64]int{
		0, 5, 10, 15, 15, 10, 5, 0,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		10, 25, 25, 25, 25, 25, 25, 10,
		10, 20, 20, 20, 20, 20, 20, 10,
	}
	bRookMG := [64]int{
		10, 20, 20, 20, 20, 20, 20, 10,
		10, 25, 25, 25, 25, 25, 25, 10,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 5, 10, 15, 15, 10, 5, 0,
	}
	wRookEG := [64]int{
		0, 10, 15, 25, 25, 15, 10, 0,
		30, 20, 20, 30, 30, 20, 20, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 20, 20, 20, 20, 20, 20, 30,
		20, 30, 30, 30, 30, 30, 30, 20,
	}
	bRookEG := [64]int{
		20, 30, 30, 30, 30, 30, 30, 20,
		30, 20, 20, 20, 20, 20, 20, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 25, 15, 10, 10, 15, 25, 30,
		30, 20, 20, 30, 30, 20, 20, 30,
		0, 10, 15, 25, 25, 15, 10, 0,
	}
	bQueenOG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 5, 0, 5, 0, 5, 0,
		0, 0, 0, 5, 5, 5, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	wQueenOG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 5, 5, 5, 0, 0, 0,
		0, 5, 0, 5, 0, 5, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	wQueenMG := [64]int{
		-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 15, 10, 10, 0, 0, -10,
		-10, 15, 5, 10, 10, 5, 0, -10,
		5, 0, -5, -5, -5, -5, 0, -5,
		-5, 0, -5, -5, -5, -5, 0, 5,
		-10, 0, -5, -5, -5, -5, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20,
	}
	bQueenMG := [64]int{
		-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, -5, -5, -5, -5, 0, -10,
		-5, 0, -5, -5, -5, -5, 0, 5,
		5, 0, -5, -5, -5, -5, 0, -5,
		-10, 15, 5, 10, 10, 5, 0, -10,
		-10, 0, 15, 10, 10, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20,
	}
	wKingMG := [64]int{
		20, 30, 15, 0, 0, 15, 30, 20,
		20, 20, 0, 0, 0, 0, 20, 20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
	}
	bKingMG := [64]int{
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 15, 0, 0, 15, 30, 20,
	}
	wKingEG := [64]int{
		-50, -30, -30, -30, -30, -30, -30, -50,
		-30, -30, 10, 20, 20, 10, -30, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -20, -10, 20, 20, -10, -20, -30,
		-50, -40, -30, -20, -20, -30, -40, -50,
	}
	bKingEG := [64]int{
		-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -20, -10, 20, 20, -10, -20, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -30, 10, 20, 20, 10, -30, -30,
		-50, -30, -30, -30, -30, -30, -30, -50,
	}
	ogW := [6][64]int{wPawnMG, wKnightOG, wBishopOG, wRookMG, wQueenOG, wKingMG}
	mgW := [6][64]int{wPawnMG, wKnightMG, wBishopMG, wRookMG, wQueenMG, wKingMG}
	egW := [6][64]int{wPawnEG, wKnightMG, wBishopMG, wRookEG, wQueenMG, wKingEG}

	ogB := [6][64]int{bPawnMG, bKnightOG, bBishopOG, bRookMG, bQueenOG, bKingMG}
	mgB := [6][64]int{bPawnMG, bKnightMG, bBishopMG, bRookMG, bQueenMG, bKingMG}
	egB := [6][64]int{bPawnEG, bKnightMG, bBishopMG, bRookEG, bQueenMG, bKingEG}

	og := [2][6][64]int{ogW, ogB}
	mg := [2][6][64]int{mgW, mgB}
	eg := [2][6][64]int{egW, egB}
	pst := [3][2][6][64]int{og, mg, eg}
	return pst
}
