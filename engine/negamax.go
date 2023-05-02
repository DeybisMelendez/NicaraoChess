package engine

import (
	"math"
	"nicarao/utils"

	chess "github.com/dylhunn/dragontoothmg"
)

// negamax function with Alpha Beta Pruning
func negamax(board *chess.Board, depth int8, turn int16, alpha int16, beta int16) int16 {

	if isTimeToStop() {
		return 0
	}

	hash := board.Hash()
	if isRepetition(hash) {
		return 0
	}
	var bestMove chess.Move
	var hashFlag = TTUpper
	var ttValue int16 = readTTEntry(hash, alpha, beta, depth, &bestMove)
	//var isPVNode bool = beta-alpha > 1
	if ttValue != NoHashEntry {
		return ttValue
	} else {
		if bestMove == 0 && depth > 1 {
			negamax(board, depth/2, turn, alpha, beta)
			var entry = getEntry(hash)
			bestMove = entry.BestMove
		}
	}

	if depth == 0 {
		return quiesce(board, alpha, beta, turn)
	}
	var maxEval int16 = -INF
	moveList := board.GenerateLegalMoves()
	var lenMoveList int = len(moveList)
	var movesSearched uint16 = 0
	var inCheck = board.OurKingInCheck()
	//var isFirstMove bool = true
	for len(moveList) > 0 {
		var move chess.Move
		var maxScoreMove uint16
		var idx int
		/*if bestMove == 0 {
			bestMove = moveList[0]
		}*/
		for i, candidate := range moveList {
			var scoreMove uint16
			if candidate == bestMove {
				scoreMove = math.MaxUint16
			} else if chess.IsCapture(candidate, board) {
				scoreMove = 50000 + uint16(getMVV_LVA(move, board))
				/*valCap := valueCapture(candidate, board)
				if valCap >= 0 {
					scoreMove = 60000 + uint16(valCap)
				} else if valCap == 0 {
					scoreMove = 50000
				}*/
			} else if isKillerMove(candidate) {
				scoreMove = 40000
			} else {
				scoreMove = 10000 + getHistoryMove(board.Wtomove, candidate)
			}
			if scoreMove > maxScoreMove {
				maxScoreMove = scoreMove
				idx = i
			}
		}
		move = moveList[idx]
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		var moveIsCapture = chess.IsCapture(move, board)
		unmakeMoveFunc := makeMove(board, move)

		var moveGiveCheck = board.OurKingInCheck()
		var isTactical = inCheck || moveGiveCheck || moveIsCapture || move.Promote() == 0
		/*if !isTactical && !isPVNode {
			// Razoring
			if depth == 2 {
				var staticEval = Evaluate(board, turn) + 100
				if staticEval < alpha {
					unmakeMove(unmakeMoveFunc)
					break
				}
				//Futility Pruning
			} else if depth == 1 {
				var staticEval = Evaluate(board, turn) + 50
				if staticEval < alpha {
					unmakeMove(unmakeMoveFunc)
					continue
				}
			}
		}*/

		var eval int16
		/*if isFirstMove {
			eval = -negamax(board, depth-1, -turn, -beta, -alpha)
			isFirstMove = false
		} else {*/
		// LMR
		if movesSearched > 5 && depth > 4 && !isTactical { // && !moveGiveCheck { // && !inCheck {
			eval = -negamax(board, depth/3, -turn, -beta, -alpha)
		} else {
			eval = alpha + 1
		}
		if eval > alpha {
			eval = -negamax(board, depth*2/3, -turn, -beta, -alpha)
			if eval > alpha && eval < beta {
				eval = -negamax(board, depth-1, -turn, -beta, -alpha)
			}
		}
		//}

		//eval := -negamax(board, depth-1, -turn, -beta, -alpha)
		unmakeMove(unmakeMoveFunc)
		maxEval = utils.MaxInt16(maxEval, eval)
		if eval > maxEval && bestMove == 0 {
			bestMove = move
		}
		if eval > alpha {
			alpha = eval
			bestMove = move
			hashFlag = TTExact
			//writeTTEntry(hash, beta, depth, TTUpper, move)
			if !moveIsCapture {
				storeHistoryMove(move, board.Wtomove, uint16(depth))
			}
			movesSearched++
			if depth > 6 {
				depth--
			}
		}
		if alpha >= beta {
			writeTTEntry(hash, beta, depth, TTLower, move)
			if !moveIsCapture {
				storeKillerMove(move)
			}
			return beta
		}
		//movesSearched++
	}
	if lenMoveList == 0 {
		if board.OurKingInCheck() {
			//Checkmate
			return -MateValue + Ply
		} else {
			//Stalemate
			return 0
		}
	}
	if bestMove != 0 {
		writeTTEntry(hash, maxEval, depth, hashFlag, bestMove)
	}
	return maxEval
}
func makeMove(board *chess.Board, move chess.Move) func() {
	Nodes++
	Ply++
	GamePly++
	StoreRepetitionTable(board.Hash())
	return board.Apply(move)
}

func unmakeMove(f func()) {
	Ply--
	GamePly--
	f()
}

// https://www.youtube.com/watch?v=QhFtquEeffA
var repetitionTable = [256]uint64{}

func isRepetition(hash uint64) bool {
	count := 0
	for i := int16(0); i < GamePly; i++ {
		if hash == repetitionTable[i] {
			count++
		}
		if count == 2 {
			return true
		}
	}
	return false
}

func StoreRepetitionTable(hash uint64) {
	repetitionTable[GamePly] = hash
}

var killerMoves [2][64]chess.Move

func storeKillerMove(move chess.Move) {
	if Ply < 64 {
		killerMoves[1][Ply] = killerMoves[0][Ply]
		killerMoves[0][Ply] = move
	}
}

func isKillerMove(move chess.Move) bool {
	if Ply < 64 {
		return killerMoves[0][Ply] == move || killerMoves[1][Ply] == move
	}
	return false
}

var historyMoves [2][64][64]uint16

func storeHistoryMove(move chess.Move, Wtomove bool, depth uint16) {
	if Wtomove {
		historyMoves[0][move.From()][move.To()] += depth * depth * depth
	} else {
		historyMoves[1][move.From()][move.To()] += depth * depth * depth
	}
}

func getHistoryMove(isWhite bool, move chess.Move) uint16 {
	if isWhite {
		return historyMoves[0][move.From()][move.To()]
	}
	return historyMoves[1][move.From()][move.To()]
}

//Most Valuable Victim - Least Valuable Aggressor

var MVV_LVA = [7][7]int16{ //[Victims][Agressors]
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

func getMVV_LVA(move chess.Move, board *chess.Board) int16 {
	agressor, _ := chess.GetPieceType(move.From(), board)
	victim, _ := chess.GetPieceType(move.To(), board)
	return MVV_LVA[victim][agressor]
}
func ValueCapture(move chess.Move, board *chess.Board) int16 {
	agressor, _ := chess.GetPieceType(move.From(), board)
	victim, _ := chess.GetPieceType(move.To(), board)
	promote := move.Promote()
	//value := Material[agressor] - Material[victim]
	value := Material[victim] - Material[agressor]
	//value := getMVV_LVA(move, board)
	if promote != chess.Nothing {
		value += Material[promote] - Material[chess.Pawn]
	}
	return value
}
