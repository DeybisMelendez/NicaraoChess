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
		hashTable = append(hashTable, hashEntry{hash: 0, depth: 0, flag: 0, score: 0, bestmove: 0})
	}
}

func WriteHashEntry(hash uint64, score int, depth int, flag int, move chess.Move) {
	index := (hash & 0x7fffffff) % hashEntries
	//var entry hashEntry = hashTable[(]
	if score < -MateScore {
		score -= Ply
	} else if score > MateScore {
		score += Ply
	}
	hashTable[index].hash = hash
	hashTable[index].depth = depth
	hashTable[index].score = score
	hashTable[index].flag = flag
	hashTable[index].bestmove = move
}

func ReadHashEntry(hash uint64, alpha int, beta int, depth int, move *chess.Move) int {
	var entry hashEntry = hashTable[(hash&0x7fffffff)%hashEntries]
	if entry.hash == hash {
		if entry.depth >= depth {
			var score int = entry.score
			if score < -MateScore {
				score += Ply
			} else if score > MateScore {
				score -= Ply
			}
			if entry.flag == HashFlagExact {
				return entry.score
			}
			if entry.flag == HashFlagAlpha && score <= alpha {
				return alpha
			}
			if entry.flag == HashFlagBeta && score >= beta {
				return beta
			}
		}
		*move = entry.bestmove
	}
	return NoHashEntry
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
