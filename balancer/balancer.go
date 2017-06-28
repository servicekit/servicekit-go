package balancer

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

type Balancer interface {
	GetResolver(ctx context.Context) naming.Resolver
}
