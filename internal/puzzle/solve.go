package puzzle

import (
	"fmt"
	"math/rand"
	"time"

	"powquote/internal/protocol"
)

var solutionSeed = time.Now().UnixNano()

func Solve(hashData *protocol.HashData, challenge protocol.Challenge) string {
	var lasthash string

	rand.Seed(solutionSeed)
	hashData.Solution = make([]byte, 128)

	for {
		_, _ = rand.Read(hashData.Solution)
		lasthash = Hash(hashData)
		if HashMatchesChallenge(lasthash, challenge) {
			fmt.Printf("=== SEED=%v ===\n", solutionSeed)
			break
		}
	}
	return lasthash
}
