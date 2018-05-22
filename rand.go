package chargen

import (
	"math/rand"
	"strconv"
	"time"
)

var (
	seed   = rand.NewSource(time.Now().UnixNano())
	random = rand.New(seed)
)

func setSeed(hash string) *rand.Rand {
	src, _ := strconv.Atoi(hash)
	newSeed := rand.NewSource(int64(src))
	return rand.New(newSeed)
}

func sampleWithoutReplacement(choices []string, n int) []string {
	samples := []string{}
	idxs := rand.Perm(len(choices))
	for i := 0; i < n; i++ {
		samples = append(samples, choices[idxs[i]])
	}
	return samples
}

func randomChoice(choices []string) string {
	r := random.Intn(len(choices))
	return choices[r]
}

func randomInt(min, max int) int {
	return random.Intn(max-min) + min
}

func weightedRandomChoice(choices []string, weights []float64) string {
	sum := 0.0
	for _, w := range weights {
		sum += w
	}
	r := random.Float64() * sum
	total := 0.0
	for i, w := range weights {
		total += w
		if r < total {
			return choices[i]
		}
	}
	return choices[0]
}
