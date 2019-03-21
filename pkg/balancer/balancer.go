package balancer

import (
	"balancer/pkg/strategies"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

const (
	// TCP network
	TCP = "tcp"
)

var (
	// ErrInvalidEndpointPrefix contains err invalid endpoint prefix
	ErrInvalidEndpointPrefix = "invelid endpoint: "
	// ErrInvalidEndpoint means that you pass invalid endpoint to balancer
	ErrInvalidEndpoint = func(endpoint string) error {
		return fmt.Errorf("%s%s", ErrInvalidEndpointPrefix, endpoint)
	}
)

// Balancer contains balancer fields
type Balancer struct {
	Strategy strategies.Strategy
	Port     int
}

// Opts contains options for balancer
type Opts struct {
	Port     int
	Strategy int
}

// New returns balancer
func New(opts *Opts) (*Balancer, error) {
	strategy := strategies.New(opts.Strategy)

	return &Balancer{
		Strategy: strategy,
		Port:     opts.Port,
	}, nil
}

// Run run listen requests to balancer
func (b *Balancer) Run(ctx context.Context) error {
	listener, err := net.Listen(TCP, fmt.Sprintf(":%d", b.Port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				return err
			}

			go b.proxy(conn)
		}
	}
}

func (b *Balancer) proxy(in net.Conn) {
	defer in.Close()

	out, err := net.Dial(TCP, b.Strategy.Get())
	if err != nil {
		return
	}
	defer out.Close()

	var (
		done = make(chan struct{})
	)

	go func() {
		io.Copy(out, in)
		done <- struct{}{}
	}()

	go func() {
		io.Copy(in, out)
		done <- struct{}{}
	}()

	<-done
	<-done
}

// AddEndpoint add endpoint to balancer
func (b *Balancer) AddEndpoint(endpoint string) error {
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return ErrInvalidEndpoint(endpoint)
	}

	b.Strategy.Push(endpoint)

	return nil
}

// IsErrInvalidEndpoint returns true if err is ErrInvalidEndpoint
func IsErrInvalidEndpoint(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), ErrInvalidEndpointPrefix)
}
