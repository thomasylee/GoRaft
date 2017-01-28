package global

import (
	"math/rand"
)

var TimeoutChannel chan bool = make(chan bool)

/**
 * Returns a timeout value between average - jitter and average + jitter.
 */
func GenerateTimeout(average int, jitter int) int {
	if jitter == 0 {
		return average
	}

	return average - jitter + rand.Intn(2 * jitter)
}
