package engine

import (
	chess "github.com/dylhunn/dragontoothmg"
)

// var FollowPV bool

func ValueMove(board *chess.Board, move chess.Move, isCapture bool, isPromo bool, pvMove chess.Move, bestmove chess.Move) int {
	if move == pvMove && FollowPV {
		return 200000
	} else if move == bestmove {
		return 190000
	} else if isPromo {
		return 180000 + Material[move.Promote()]
		/*else if isCapture || isPromo {
		piece, _ := utils.GetPiece(move.From(), board)
		capture, _ := utils.GetPiece(move.To(), board)
		promo := move.Promote()
		if Material[capture] > Material[piece] {
			return 4000 + Material[capture] + Material[promo]
		}
		if isPromo {
			return 5000 + Material[promo]
		}
		return GetMVV_LVA(move, board)*/
	} else if isCapture {
		if SEE(board, move, move.To(), 0, 1) > 0 {
			return 3000 + GetMVV_LVA(move, board)
		} else {
			return GetMVV_LVA(move, board)
		}
	} else if KillerMoves[0][Ply] == move {
		return 2000 // + GetHistoryMove(board.Wtomove, move)
	} else if KillerMoves[1][Ply] == move {
		return 1000 // + GetHistoryMove(board.Wtomove, move)
	} else {
		return GetHistoryMove(board.Wtomove, move)
	}
}
