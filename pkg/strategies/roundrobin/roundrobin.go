package roundrobin

import (
	"container/ring"
)

type RoundRobin struct {
	ring *ring.Ring
}

func New() *RoundRobin {
	return &RoundRobin{}
}

func (rr *RoundRobin) Push(endpoint string) {
	tmpRing := ring.New(1)
	tmpRing.Value = endpoint

	if rr.ring.Len() == 0 {
		rr.ring = tmpRing
		return
	}

	rr.ring.Link(tmpRing)
}

func (rr *RoundRobin) Get() string {
	if rr.ring.Len() == 0 {
		return ""
	}

	rr.ring = rr.ring.Next()
	return rr.ring.Value.(string)
}
