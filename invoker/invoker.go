package invoker

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Invoker has a method Invoke used to invoke grpc.Invoke
type Invoker interface {
	Invoke(ctx context.Context, conn *grpc.ClientConn, method string, request, response interface{}, opts ...grpc.CallOption) error
}

// FailfastInvoker will be returned immediately when invoking timeout
type FailfastInvoker struct {
	timeout time.Duration
}

// FailsafeInvoker will be returned immediately when invoking timeout or failed
type FailsafeInvoker struct {
	timeout time.Duration
}
