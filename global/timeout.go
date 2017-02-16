package global

import (
	"math/rand"
)

var TimeoutChannel chan bool = make(chan bool)

/**
 * Returns a timeout value between average - jitter and average + jitter.
 */
func GenerateTimeout(average uint32, jitter uint32) uint32 {
	if jitter == 0 {
		return average
	}

	return average - jitter + uint32(rand.Intn(2 * int(jitter)))
}
