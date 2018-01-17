package balancer

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// Balancer GRPC Balancer interface
// See: https://github.com/grpc/grpc/blob/master/doc/load-balancing.md
type Balancer interface {
	GetResolver(ctx context.Context) naming.Resolver
}
