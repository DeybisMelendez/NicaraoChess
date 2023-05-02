package engine

import (
	"math/bits"

	chess "github.com/dylhunn/dragontoothmg"
)

const (
	PawnPhase   = 0
	KnightPhase = 1
	BishopPhase = 1
	RookPhase   = 2
	QueenPhase  = 4
	OPENING     = 0
	ENDGAME     = 1
)

var TotalPhase int = KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2 // + PawnPhase*16
var Material = [7]int16{0, 1, 3, 3, 5, 9, 10}
var MaterialOpening = [7]int{0, 82, 337, 365, 477, 1025, 0}
var MaterialEndgame = [7]int{0, 94, 281, 297, 512, 936, 0}

var MaterialScore = [2][7]int{MaterialOpening, MaterialEndgame} //Opening, Endgame

//var PST = GetPieceSquareTable()

func materialEval(board *chess.Board, phase int) int16 {
	var white int = bits.OnesCount64(board.White.Pawns)*MaterialScore[phase][chess.Pawn] +
		bits.OnesCount64(board.White.Knights)*MaterialScore[phase][chess.Knight] +
		bits.OnesCount64(board.White.Bishops)*MaterialScore[phase][chess.Bishop] +
		bits.OnesCount64(board.White.Rooks)*MaterialScore[phase][chess.Rook] +
		bits.OnesCount64(board.White.Queens)*MaterialScore[phase][chess.Queen]
		//bits.OnesCount64(board.White.Kings)*MaterialScore[phase][chess.King]

	var black int = bits.OnesCount64(board.Black.Pawns)*MaterialScore[phase][chess.Pawn] +
		bits.OnesCount64(board.Black.Knights)*MaterialScore[phase][chess.Knight] +
		bits.OnesCount64(board.Black.Bishops)*MaterialScore[phase][chess.Bishop] +
		bits.OnesCount64(board.Black.Rooks)*MaterialScore[phase][chess.Rook] +
		bits.OnesCount64(board.Black.Queens)*MaterialScore[phase][chess.Queen]
		//bits.OnesCount64(board.Black.Kings)*MaterialScore[phase][chess.King]
	return int16(white - black)
}

func piecePosEval(pieces uint64, pieceType chess.Piece, isWhite bool, board *chess.Board, whiteKing uint8, blackKing uint8, allPieces uint64) (int16, int16) {
	var opening int16
	var endgame int16
	for pieces != 0 {
		var square uint8 = uint8(bits.TrailingZeros64(pieces))
		if isWhite {
			//PST WHITE
			opening += PST[OPENING][pieceType][ReversedBoard[square]]
			endgame += PST[ENDGAME][pieceType][ReversedBoard[square]]
			// King Tropism
			opening += KING_TROPISM[pieceType][square][blackKing] * Material[pieceType]
		} else {
			//PST BLACK
			opening += PST[OPENING][pieceType][square]
			endgame += PST[ENDGAME][pieceType][square]
			// King tropism
			opening += KING_TROPISM[pieceType][square][whiteKing] * Material[pieceType]
		}

		// Mobility Bishop, Knight and Rook only
		mobility := mobility(pieceType, isWhite, square, board, allPieces)
		opening += mobility
		endgame += mobility
		// Individual Evaluation Pieces
		switch pieceType {
		case chess.Pawn:
			var pawnScore int16 = -doubledPawns(board, isWhite, square) - isolatedPawns(board, isWhite, square) // + passedPawn(board, isWhite, square)
			opening += pawnScore
			endgame += pawnScore
		case chess.Knight:
			opening += goodKnight(board)
		case chess.Bishop:
			var bishopScore int16 = bishopPair(board, isWhite) - badBishop(board, isWhite, square)
			opening += bishopScore
			endgame += bishopScore
		case chess.Rook:
			opening += rookToQueen(board, isWhite, square) + goodRook(board)
		/*case chess.Queen:
		opening -= badQueen(board, isWhite, square)*/
		case chess.King:
			opening -= kingSafety(isWhite, square, board, allPieces)
		}
		pieces &= pieces - 1
	}
	return opening, endgame
}

func positionalEval(board *chess.Board) (int16, int16) {
	var whiteOpening int16
	var whiteEndgame int16
	var blackOpening int16
	var blackEndgame int16
	whiteKing := uint8(bits.TrailingZeros64(board.White.Kings))
	blackKing := uint8(bits.TrailingZeros64(board.Black.Kings))
	allPieces := board.White.All | board.Black.All
	//White Evaluations
	pieceList := []uint64{board.White.Pawns, board.White.Knights, board.White.Bishops,
		board.White.Rooks, board.White.Queens, board.White.Kings}
	for i, pieces := range pieceList {
		op, en := piecePosEval(pieces, chess.Piece(i+1), true, board, whiteKing, blackKing, allPieces)
		whiteOpening += op
		whiteEndgame += en
	}
	//Black Evaluations
	pieceList = []uint64{board.Black.Pawns, board.Black.Knights, board.Black.Bishops,
		board.Black.Rooks, board.Black.Queens, board.Black.Kings}
	for i, pieces := range pieceList {
		op, en := piecePosEval(pieces, chess.Piece(i+1), false, board, whiteKing, blackKing, allPieces)
		blackOpening += op
		blackEndgame += en
	}
	var opening int16 = whiteOpening - blackOpening
	var endgame int16 = whiteEndgame - blackEndgame
	return opening, endgame
}

func Evaluate(board *chess.Board, turn int16) int16 {
	if isDraw(board) {
		return 0
	}
	var opening, endgame = positionalEval(board)
	opening += materialEval(board, OPENING)
	endgame += materialEval(board, ENDGAME)

	phase := (bits.OnesCount64(board.White.Knights|board.Black.Knights) * KnightPhase) +
		(bits.OnesCount64(board.White.Bishops|board.Black.Bishops) * BishopPhase) +
		(bits.OnesCount64(board.White.Rooks|board.Black.Rooks) * RookPhase) +
		(bits.OnesCount64(board.White.Queens|board.Black.Queens) * QueenPhase)
	phase = (phase*256 + (TotalPhase / 2)) / TotalPhase
	score := (int(opening)*(256-phase) + int(endgame)*phase) / 256
	return int16(score) * turn
}
