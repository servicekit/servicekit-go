package consul

import (
    "golang.org/x/net/context"
    "google.golang.org/grpc/naming"

    "github.com/servicekit/servicekit-go/balancer"
    "github.com/servicekit/servicekit-go/coordinator"
    "github.com/servicekit/servicekit-go/logger"
)

type ConsulBalancer struct {
    consul   coordinator.Coordinator
    resolver naming.Resolver

    log *logger.Logger
}

func NewConsulBalancer(consul coordinator.Coordinator, service, tag string, log *logger.Logger) (balancer.Balancer, error) {
    resolver, err := newConsulResolver(consul, service, tag, log)
    if err != nil {
        return nil, err
    }

    return &ConsulBalancer{
        consul:   consul,
        resolver: resolver,
    }, nil
}

func (b *ConsulBalancer) GetResolver(ctx context.Context) naming.Resolver {
    return b.resolver
}

func GetResolver(consul coordinator.Coordinator, service, tag string, log *logger.Logger) naming.Resolver {
    resolver, err := newConsulResolver(consul, service, tag, log)
    if err != nil {
        return nil, err
    }

    return resolver, nil
}
