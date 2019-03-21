package strategies

import (
	"balancer/pkg/strategies/roundrobin"
)

const (
	// RoundRobin strategy
	RoundRobin = iota
)

// Strategy interface
type Strategy interface {
	Push(string)
	Get() string
}

// New returns new strategy
func New(strategy int) Strategy {
	switch strategy {
	case RoundRobin:
		return roundrobin.New()
	default:
		return nil
	}
}
