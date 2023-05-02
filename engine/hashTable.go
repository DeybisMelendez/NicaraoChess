package engine

import (
	"fmt"
	"math"
	"unsafe"

	chess "github.com/dylhunn/dragontoothmg"
)

const (
	TTExact int8 = iota
	TTLower
	TTUpper
	NoHashEntry = math.MinInt16
)

var ttSize uint64 = 8 * 1024 * 1024 / 20

type TranspositionTable struct {
	Hash     uint64
	Depth    int8
	Flag     int8
	Value    int16
	BestMove chess.Move
}

var transpositionTable []TranspositionTable

func ResizeTranspositionTable(sizeMB uint64) {
	// Convertir MB a bytes
	var sizeBytes uint64 = sizeMB * 1024 * 1024
	// Calcular el nuevo tamaño de la tabla de transposición
	ttSize = sizeBytes / uint64(unsafe.Sizeof(TranspositionTable{})+unsafe.Alignof(TranspositionTable{}))
	transpositionTable = make([]TranspositionTable, ttSize)
	fmt.Println("Set hash table size to", sizeMB, "Mb")
	fmt.Println("Hash table initialized with", ttSize, "entries")
}

func writeTTEntry(hash uint64, value int16, depth int8, flag int8, bestmove chess.Move) {
	ttEntry := &transpositionTable[hash%ttSize]
	if depth > ttEntry.Depth {
		ttEntry.Depth = depth
		ttEntry.Flag = flag
		ttEntry.Value = value
		ttEntry.BestMove = bestmove
		ttEntry.Hash = hash
	}
}

func readTTEntry(hash uint64, alpha int16, beta int16, depth int8, move *chess.Move) int16 {
	var ttEntry TranspositionTable = getEntry(hash)
	if ttEntry.Hash == hash {
		if ttEntry.Depth >= depth {
			if ttEntry.Flag == TTExact {
				return ttEntry.Value
			}
			if ttEntry.Flag == TTUpper && ttEntry.Value <= alpha {
				return alpha
			}
			if ttEntry.Flag == TTLower && ttEntry.Value >= beta {
				return beta
			}
		}
		*move = ttEntry.BestMove
	}
	return NoHashEntry
}

func getPV(fen string, depth int8, bestmove *chess.Move) string {
	var pv string
	var board chess.Board = chess.ParseFen(fen)

	for i := int8(0); i < depth; i++ {
		hash := board.Hash()
		entry := getEntry(hash)
		if entry.Hash == hash && entry.BestMove != 0 {
			if i == 0 {
				*bestmove = entry.BestMove
			}
			pv += entry.BestMove.String() + " "
			board.Apply(entry.BestMove)
		} else {
			break
		}
	}
	return pv
}

func getEntry(hash uint64) TranspositionTable {
	return transpositionTable[hash%ttSize]
}
