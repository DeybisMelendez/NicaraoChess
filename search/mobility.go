package search

//TODO Muy Lento para implementarlo en la evaluación!!!
import (
	"nicarao/utils"
	"strings"

	chess "github.com/dylhunn/dragontoothmg"
)

func Mobility(board *chess.Board) int {
	white, black := 0, 0
	fen := ""
	if board.Wtomove {
		fen = strings.ReplaceAll(board.ToFen(), " w ", " b ")
	} else {
		fen = strings.ReplaceAll(board.ToFen(), " b ", " w ")
	}
	rivalBoard := chess.ParseFen(fen)
	whiteMoves := board.GenerateLegalMoves()
	blackMoves := rivalBoard.GenerateLegalMoves()
	for i := 0; i < len(whiteMoves); i++ {
		piece, _ := utils.GetPiece(uint8(i), board)
		if piece != chess.Nothing && piece != chess.Queen && piece != chess.King {
			white++
		}
	}
	for i := 0; i < len(blackMoves); i++ {
		piece, _ := utils.GetPiece(uint8(i), board)
		if piece != chess.Nothing && piece != chess.Queen && piece != chess.King {
			black++
		}
	}
	return (white * 5) - (black * 5)
}
