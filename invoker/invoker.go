package invoker

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Invoker interface {
	Invoke(ctx context.Context, conn *grpc.ClientConn, method string, request, response interface{}, opts ...grpc.CallOption) error
}

type FailfastInvoker struct {
	timeout time.Duration
}

type FailsafeInvoker struct {
	timeout time.Duration
}
