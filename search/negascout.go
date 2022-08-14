package search

import (
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

const FullDepthMove = 6

func pvReduction(depth int) int {
	return depth / 3
}

func isLMROk(board *chess.Board, move chess.Move) bool {
	var check bool = board.OurKingInCheck()
	var isCapture bool = chess.IsCapture(move, board)
	return check && isCapture
}

func Negascout(board *chess.Board, depth int, color int, alpha int, beta int) int {
	PVLength[Ply] = Ply
	var hashFlag int = HashFlagAlpha
	var moveList []chess.Move = board.GenerateLegalMoves()
	var score int = 0
	var hashScore = ReadHashEntry(board.Hash(), alpha, beta, depth)
	if hashScore != NoHashEntry && Ply > 0 {
		return hashScore
	}
	if depth == 0 || len(moveList) == 0 {
		return 0
	}
	for i := 0; i < len(moveList); i++ {
		move := moveList[i]
		GetPieces(move, board)
		unmakeFunc := Make(board, move)
		if i > FullDepthMove && isLMROk(board, move) {
			score = -Negascout(board, pvReduction(depth), -color, -alpha-1, -alpha)
			if score > alpha && score < beta {
				////https://www.chessprogramming.org/NegaScout#Guido_Schimmels
				var score2 int = -Negascout(board, depth-1, -color, -beta, -alpha)
				score = utils.Max(score, score2)
			}
		} else {
			score = -Negascout(board, depth-1, -color, -beta, -alpha)
		}
		Unmake(unmakeFunc)
		if score > alpha {
			StorePV(move)
			hashFlag = HashFlagExact
			alpha = score
		}
		if score >= beta {
			StoreKillerMove(move, board)
			WriteHashEntry(board.Hash(), beta, depth, HashFlagBeta)
			return beta
		}
	}
	WriteHashEntry(board.Hash(), beta, depth, hashFlag)
	return alpha
}
