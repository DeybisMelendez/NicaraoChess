package search

const HashFlagExact, HashFlagAlpha, HashFlagBeta, NoHashEntry = 0, 1, 2, 10000

// 16Mb default hash table size
var hashEntries uint64 = 838860
var hashTable []hashEntry

type hashEntry struct {
	hash  uint64
	depth int
	flag  int
	score int
	//bestmove
}

func InitHasTable() {
	var entries int = int(hashEntries)
	for i := 0; i < entries; i++ {
		hashTable = append(hashTable, hashEntry{hash: 0, depth: 0, flag: 0, score: 0})
	}
}

func WriteHashEntry(hash uint64, score int, depth int, flag int) {
	var entry hashEntry = hashTable[(hash&0x7fffffff)%hashEntries]
	entry.hash = hash
	entry.depth = depth
	entry.score = score
	entry.flag = flag
}

func ReadHashEntry(hash uint64, alpha int, beta int, depth int) int {
	var entry hashEntry = hashTable[(hash&0x7fffffff)%hashEntries]
	if entry.hash == hash {
		if entry.depth >= depth {
			var score int = entry.score
			if entry.flag == HashFlagExact {
				return score
			}
			if entry.flag == HashFlagAlpha && score <= alpha {
				return alpha
			}
			if entry.flag == HashFlagBeta && score >= beta {
				return beta
			}
		}
	}
	return NoHashEntry
}
