package balancer_test

import (
	"balancer/pkg/balancer"
	"balancer/pkg/strategies"
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEndpoints(t *testing.T) {
	var cases = []struct {
		name     string
		endpoint string
		expected bool
	}{
		{
			"success",
			"http://endpoint.com/",
			false,
		},
		{
			"invalid",
			"smth",
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			b, err := balancer.New(&balancer.Opts{})
			require.NoError(t, err)

			err = b.AddEndpoint(c.endpoint)
			require.Equal(t, c.expected, balancer.IsErrInvalidEndpoint(err))
		})
	}
}

const (
	TCPServer1Port = 1001
	TCPServer2Port = 1002
	BalancerPort   = 2001
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go TCPServer(TCPServer1Port, ctx, t)
	go TCPServer(TCPServer2Port, ctx, t)

	b, err := balancer.New(&balancer.Opts{
		Port:     BalancerPort,
		Strategy: strategies.RoundRobin,
	})
	require.NoError(t, err)

	err = b.AddEndpoint(fmt.Sprintf("localhost:%d", TCPServer1Port))
	require.NoError(t, err)
	err = b.AddEndpoint(fmt.Sprintf("localhost:%d", TCPServer2Port))
	require.NoError(t, err)

	go func() {
		err := b.Run(ctx)
		require.NoError(t, err)
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", BalancerPort))
	require.NoError(t, err)

	_, err = conn.Write([]byte{1, 2, 3})
	require.NoError(t, err)

	var buff = make([]byte, 3)
	_, err = conn.Read(buff)
	require.NoError(t, err)

	require.Equal(t, []byte{1, 2, 3}, buff)
}

func TCPServer(port int, ctx context.Context, t *testing.T) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	require.NoError(t, err)
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			require.NoError(t, err)
			defer conn.Close()

			var buff = make([]byte, 3)

			_, err = conn.Read(buff)
			require.NoError(t, err)

			_, err = conn.Write(buff)
			require.NoError(t, err)
		}
	}
}
