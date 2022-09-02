package engine

/*
func ZWSearch(board *chess.Board, depth int, beta int, turn int, nullMove bool) int {
	var score int = 0
	var alpha int = beta - 1
	var hashFlag int = HashFlagAlpha
	var bestmove chess.Move
	var isPVNode bool = beta-alpha > 1
	ReadHashEntry(board.Hash(), alpha, beta, depth, &bestmove)
	if hashScore != NoHashEntry && Ply > 0 && !isPVNode {
		return hashScore
	}
	if isTimeToStop() {
		return 0
	}
	if depth == 0 {
		return Quiesce(board, alpha, beta, turn) //Evaluate(board,turn) //
	}
	var inCheck bool = board.OurKingInCheck()
	if nullMove && !inCheck {
		if Ply > 0 && depth > NullDepth && !isEndgame(board) {
			var staticEval int = Evaluate(board, turn)
			if staticEval >= beta {
				board.Wtomove = !board.Wtomove
				nullBoard := chess.ParseFen(board.ToFen())
				board.Wtomove = !board.Wtomove
				if !nullBoard.OurKingInCheck() && len(nullBoard.GenerateLegalMoves()) != 0 {
					eval := -ZWSearch(&nullBoard, depth-NullDepth-1, alpha, -turn, NoNull)
					if eval >= beta {
						return eval
					}
				}
			}
		}
	}
	moveList := board.GenerateLegalMoves()
	var lenMoveList int = len(moveList)
	for len(moveList) > 0 {
		var val int = -1
		var idx int = 0
		var ln int = len(moveList)
		for i := 0; i < ln; i++ {
			var newVal int = moveOrdering.ValueMove(board, moveList[i], 0, 0, Ply, false)
			if newVal > val {
				val = newVal
				idx = i
			}
		}
		var move = moveList[idx]
		moveList = append(moveList[:idx], moveList[idx+1:]...)
		unmakeFunc := Make(board, move)
		score = -ZWSearch(board, depth-1, 1-beta, -turn, nullMove)
		Unmake(unmakeFunc)
		if isTimeToStop() {
			return 0
		}
		if score >= beta {
			//WriteHashEntry(board.Hash(), beta, depth, HashFlagBeta, move)
			return beta
		}
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
	//WriteHashEntry(board.Hash(), beta, depth, hashFlag, bestmove)
	return beta - 1
}
*/
