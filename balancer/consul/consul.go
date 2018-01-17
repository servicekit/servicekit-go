package consul

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"

	"github.com/servicekit/servicekit-go/balancer"
	"github.com/servicekit/servicekit-go/coordinator"
	"github.com/servicekit/servicekit-go/logger"
)

// Balancer save meta info in Consul
type Balancer struct {
	consul   coordinator.Coordinator
	resolver naming.Resolver

	log *logger.Logger
}

// NewBalancer initializes and returns a new Balancer.
func NewBalancer(consul coordinator.Coordinator, service, tag string, log *logger.Logger) (balancer.Balancer, error) {
	resolver, err := newResolver(consul, service, tag, log)
	if err != nil {
		return nil, err
	}

	return &Balancer{
		consul:   consul,
		resolver: resolver,
	}, nil
}

// GetResolver returns a Resolver
func (b *Balancer) GetResolver(ctx context.Context) naming.Resolver {
	return b.resolver
}

// GetResolver returns a Resolver
func GetResolver(consul coordinator.Coordinator, service, tag string, log *logger.Logger) (naming.Resolver, error) {
	resolver, err := newResolver(consul, service, tag, log)
	if err != nil {
		return nil, err
	}

	return resolver, nil
}
