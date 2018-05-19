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
	weightSum := 0.0
	for _, w := range weights {
		weightSum += w
	}
	cumWeights := []float64{}
	for i, w := range weights {
		if i == 0 {
			cumWeights = append(cumWeights, w)
		} else {
			cumWeights = append(cumWeights, cumWeights[i-1]+w)
		}
	}
	r := random.Float64() * weightSum
	var out string
	for i, w := range weights {
		if r < w/cumWeights[i] {
			out = choices[i]
			break
		}
	}
	return out
}
