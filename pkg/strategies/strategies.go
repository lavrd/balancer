package strategies

import (
	"balancer/pkg/strategies/roundrobin"
)

const (
	RoundRobin = iota
)

type Strategy interface {
	Push(string)
	Get() string
}

func New(strategy int) Strategy {
	switch strategy {
	case RoundRobin:
		return roundrobin.New()
	default:
		return nil
	}
}
