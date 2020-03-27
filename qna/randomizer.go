package qna

import "math/rand"

type Randomizer struct {
	rand       *rand.Rand
	lastRandom int
}

func (r *Randomizer) Reset(seed int64) {
	r.rand = rand.New(rand.NewSource(seed))
	r.lastRandom = -1
}

func (r *Randomizer) Random(inclusiveMin, exclusiveMax int) int {
	if r.rand == nil {
		r.Reset(1)
	}
	// say max was 2, and min 0.
	// we should get either 0 or 1.
	// length is 2-0= 2, n = [0..2), ret= 2+0
	length := exclusiveMax - inclusiveMin
	n := r.rand.Intn(length) + inclusiveMin // [0,length)
	if n == r.lastRandom {
		n = ((n + 1) % length) + inclusiveMin // also, [0,length)
	}
	r.lastRandom = n
	return n
}
