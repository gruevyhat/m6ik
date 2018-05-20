package chargen

import (
	"math/rand"
	"time"
)

var (
	seed   = rand.NewSource(time.Now().Unix())
	random = rand.New(seed)
)

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
