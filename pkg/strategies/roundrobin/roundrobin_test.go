package roundrobin_test

import (
	"balancer/pkg/strategies/roundrobin"

	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	var cases = []struct {
		name      string
		endpoints []string
		expected  int
	}{
		{
			"3 endpoints",
			[]string{"1", "2", "3"},
			3,
		},
		{
			"zero endpoints",
			[]string{},
			0,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rr := roundrobin.New()

			for _, e := range c.endpoints {
				rr.Push(e)
			}

			var endpoints []string

			for _, e := range c.endpoints {
				endpoints = append(endpoints, e)
			}

			var founded int
			for _, ee := range c.endpoints {
				for _, ae := range endpoints {
					if ee == ae {
						founded++
					}
				}
			}

			require.Equal(t, c.expected, founded)
		})
	}
}
