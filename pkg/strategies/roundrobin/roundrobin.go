package roundrobin

import (
	"container/ring"
)

// RoundRobin contains fields for round robin strategy
type RoundRobin struct {
	ring *ring.Ring
}

// New returns new round robing strategy
func New() *RoundRobin {
	return &RoundRobin{}
}

// Push push new endpoint
func (rr *RoundRobin) Push(endpoint string) {
	tmpRing := ring.New(1)
	tmpRing.Value = endpoint

	if rr.ring.Len() == 0 {
		rr.ring = tmpRing
		return
	}

	rr.ring.Link(tmpRing)
}

// Get get next endpoint
func (rr *RoundRobin) Get() string {
	if rr.ring.Len() == 0 {
		return ""
	}

	rr.ring = rr.ring.Next()
	return rr.ring.Value.(string)
}
