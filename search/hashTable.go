package search

import (
	"fmt"

	chess "github.com/dylhunn/dragontoothmg"
)

const HashFlagExact, HashFlagAlpha, HashFlagBeta, NoHashEntry = 0, 1, 2, 10000

// 16Mb default hash table size
var Mb uint64 = 16
var hashEntries uint64 = 838860
var hashTable []hashEntry

type hashEntry struct {
	hash     uint64
	depth    int
	flag     int
	score    int
	bestmove chess.Move
}

func InitHasTable() {
	var newHashTable []hashEntry
	hashTable = newHashTable
	var entries int = int(hashEntries)
	for i := 0; i < entries; i++ {
		hashTable = append(hashTable, hashEntry{hash: 0, depth: 0, flag: 0, score: 0})
	}
}

func WriteHashEntry(hash uint64, score int, depth int, flag int, move chess.Move) {
	var entry hashEntry = hashTable[(hash&0x7fffffff)%hashEntries]
	entry.hash = hash
	entry.depth = depth
	entry.score = score
	entry.flag = flag
	entry.bestmove = move
}

func ReadHashEntry(hash uint64, alpha int, beta int, depth int) (int, chess.Move) {
	var entry hashEntry = hashTable[(hash&0x7fffffff)%hashEntries]
	if entry.hash == hash {
		if entry.depth >= depth {
			var score int = entry.score
			if entry.flag == HashFlagExact {
				return entry.score, entry.bestmove
			}
			if entry.flag == HashFlagAlpha && score <= alpha {
				return alpha, entry.bestmove
			}
			if entry.flag == HashFlagBeta && score >= beta {
				return beta, entry.bestmove
			}
		}
	}
	return NoHashEntry, 0
}

func SetHashTable(mb uint64) {
	//var newHashTable []hashEntry
	//hashTable = newHashTable
	if mb < 4 {
		Mb = 4
	} else if mb > 128 {
		Mb = 128
	}
	hashEntries = mb * 0x100000 / 20
	InitHasTable()
	fmt.Println("Set hash table size to", Mb, "Mb")
	fmt.Println("Hash table initialized with", hashEntries, "entries")
}
