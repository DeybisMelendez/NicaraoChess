package engine

var PST = [2][7][64]int{} //[phase][piece][square]

var ReversedBoard = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

func PstMake() [2][7][64]int {
	// El tablero se miraría como si llevara las piezas negras
	PawnOG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		53, 50, 47, 49, 53, 43, 61, 55,
		8, 11, 19, 29, 32, 23, 7, 12,
		4, 5, 10, 24, 25, 4, 7, 8,
		0, -1, 1, 21, 19, -2, -1, 0,
		5, -6, -8, -3, 6, -11, -5, 4,
		5, 9, 9, -30, -11, 10, 10, 6,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	KnightOG := [64]int{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, -2, 5, 0, -20, -40,
		-28, 5, 6, 13, 16, 10, -11, -38,
		-31, 0, 17, 9, 10, 13, 14, -30,
		-19, -11, 15, 9, 16, 15, 1, -37,
		-30, 4, 3, 6, 14, 14, -6, -30,
		-40, -20, -11, 4, 5, 5, -20, -29,
		-50, -29, -30, -33, -41, -41, -39, -50,
	}
	BishopOG := [64]int{
		-20, -10, -10, -10, -10, -14, -10, -20,
		-10, 0, 0, 9, 0, 0, 1, -10,
		-10, -2, -4, 10, 10, 4, 0, 1,
		-10, 4, 7, 21, -1, -6, 11, -9,
		1, 0, 9, 10, 10, 9, -1, -12,
		-1, 9, 9, 11, 12, 9, -1, -10,
		-10, 16, 6, 11, 10, 9, 4, -10,
		-20, -10, -20, -13, -1, -9, -10, -20,
	}
	RookOG := [64]int{
		-1, 20, 20, 20, 19, 21, 20, 10,
		8, 29, 25, 25, 25, 36, 16, 12,
		6, 2, 11, -2, 2, 10, 11, 5,
		-5, -1, 11, -1, 0, -1, -1, -6,
		-5, 0, -3, 0, 0, -6, 0, -1,
		-6, -2, -2, 0, 2, 0, 1, 4,
		-5, -1, -3, 6, 2, 2, 11, -14,
		-4, 1, 11, 20, 21, 9, 1, 4,
	}
	QueenOG := [64]int{
		-23, -10, -10, -5, -5, -10, -9, -16,
		-10, -1, 0, -1, -9, 0, -6, -8,
		1, 0, 3, 3, -1, -1, 9, -10,
		-4, -11, 5, -6, 0, 5, 2, -4,
		1, 3, 5, 6, 5, 9, -1, -5,
		-17, 16, 3, -6, 6, 10, 3, -10,
		-10, 1, 5, 0, 0, 11, 2, -21,
		-25, -10, -10, -2, -6, 0, -10, -20,
	}
	KingOG := [64]int{
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -42, -41, -48, -50, -40, -40, -30,
		-19, -29, -41, -51, -50, -40, -40, -30,
		-30, -29, -49, -50, -50, -40, -40, -30,
		-21, -30, -29, -36, -40, -28, -30, -20,
		-12, -20, -29, -20, -19, -20, -20, -11,
		20, 9, -21, -17, -18, -14, 20, 21,
		20, 30, -11, -10, 0, -10, 29, 20,
	}
	PawnEG := [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		22, 60, 49, 60, 71, 56, 71, 30,
		10, 30, 30, 28, 38, 36, 39, 16,
		5, 10, 9, 10, 10, 12, 10, 10,
		0, 0, 6, 3, 5, 6, 2, 0,
		0, 1, 1, 0, 11, 0, 0, 0,
		0, 0, 0, 0, 4, -1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	KnightEG := [64]int{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, -7, 0, 0, -20, -40,
		-31, 11, 9, 15, 15, 10, 0, -30,
		-29, 7, 17, 19, 11, 15, 12, -30,
		-20, -5, 15, 16, 21, 14, 1, -26,
		-30, 4, 10, 13, 15, 11, 11, -30,
		-40, -20, -5, 4, 5, 5, -20, -40,
		-50, -39, -30, -31, -32, -41, -39, -50,
	}
	BishopEG := [64]int{
		-20, -10, -10, -10, -10, -19, -10, -20,
		-10, 0, 0, -10, 0, 0, 0, -10,
		-10, -1, 5, 10, 10, 2, 0, -4,
		-10, 15, 5, 5, 10, 0, 16, -10,
		-10, 0, 20, 10, 10, 17, 0, -12,
		-12, 10, 10, 10, 10, 10, 9, -10,
		-10, 5, 1, 7, 7, 1, 5, -10,
		-20, -10, -20, -10, -4, -22, -10, -20,
	}
	RookEG := [64]int{
		-1, 20, 20, 20, 19, 21, 20, 10,
		8, 29, 25, 25, 25, 36, 16, 12,
		6, 2, 11, -2, 2, 10, 11, 5,
		-5, -1, 11, -1, 0, -1, -1, -6,
		-5, 0, -3, 0, 0, -6, 0, -1,
		-6, -2, -2, 0, 2, 0, 1, 4,
		-5, -1, -3, 6, 2, 2, 11, -14,
		-4, 1, 11, 20, 21, 9, 1, 4,
	}
	QueenEG := [64]int{
		-21, -10, -10, -5, -5, -10, -9, -17,
		-10, 1, 0, -1, -2, 0, -1, -13,
		-9, 0, 5, 4, 3, 0, 6, -10,
		-5, -11, 5, 5, -6, 5, 5, -5,
		3, 0, 5, 5, 5, 8, 0, -5,
		-14, 6, 6, 10, 5, 5, -1, -10,
		-9, 0, 5, 0, 0, 11, 2, -11,
		-22, -10, -10, -3, -5, -8, -10, -20,
	}
	KingEG := [64]int{
		-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -21, -21, 9, 0, -10, -20, -30,
		-19, 0, 19, 23, 30, 20, -10, -30,
		-30, -6, 19, 40, 40, 30, 1, -30,
		-29, -10, 19, 47, 40, 33, -10, -30,
		-33, -10, 9, 37, 34, 19, -8, -30,
		-31, -40, -11, -10, 5, -2, -30, -29,
		-50, -30, -30, -30, -24, -30, -30, -51,
	}
	nothing := [64]int{}
	opening := [7][64]int{nothing, PawnOG, KnightOG, BishopOG, RookOG, QueenOG, KingOG}

	endgame := [7][64]int{nothing, PawnEG, KnightEG, BishopEG, RookEG, QueenEG, KingEG}

	/*whiteMiddle := [7][64]int{nothing, utils.Reverse(PawnMG), utils.Reverse(KnightMG), utils.Reverse(BishopMG),
		utils.Reverse(RookMG), utils.Reverse(QueenMG), utils.Reverse(KingMG)}

	whiteEnd := [7][64]int{nothing, utils.Reverse(PawnEG), utils.Reverse(KnightMG), utils.Reverse(BishopMG),
		utils.Reverse(RookEG), utils.Reverse(QueenMG), utils.Reverse(KingEG)}

	white := [2][7][64]int{whiteMiddle, whiteEnd}
	black := [2][7][64]int{blackMiddle, blackEnd}
	*/
	pst := [2][7][64]int{opening, endgame}
	return pst
}