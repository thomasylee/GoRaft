package global

import (
	"math/rand"
)

// TimeoutChannel is the channel used for kicking off leader election when
// a node has not received a leader heartbeat in a while.
var TimeoutChannel chan bool = make(chan bool)

// GenerateTimeout returns a timeout value between average - jitter and
// average + jitter.
func GenerateTimeout(average uint32, jitter uint32) uint32 {
	if jitter == 0 {
		return average
	}

	return average - jitter + uint32(rand.Intn(2*int(jitter)))
}
