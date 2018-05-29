package m6ik

import (
	"hash/fnv"
	"math/rand"
	"strconv"
	"time"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func setSeed(seed string) *rand.Rand {
	h := int64(hash(seed))
	newSeed := rand.NewSource(h)
	return rand.New(newSeed)
}

var (
	seed   = strconv.FormatInt(time.Now().UTC().UnixNano(), 16)
	random = setSeed(seed)
)

func sampleWithoutReplacement(choices []string, n int) []string {
	samples := []string{}
	idxs := random.Perm(len(choices))
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
	r := random.Float64()*sum - 1.0
	total := 0.0
	for i, w := range weights {
		total += w
		if r <= total {
			return choices[i]
		}
	}
	return choices[0]
}
