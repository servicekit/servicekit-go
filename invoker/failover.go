package invoker

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// MaxTries represents the max count for retry
	MaxTries = 8
	// MaxTimeout represents the timeout seconds
	MaxTimeout = 8 * time.Second
	// MaxDelay represents the delay seconds
	MaxDelay = 8 * time.Minute
)

// FailoverInvoker can be retry when invoke failed
type FailoverInvoker struct {
	tries   int
	timeout time.Duration
	delay   Delay
}

// NewFailoverInvoker returns an Invoker
func NewFailoverInvoker(tries int, timeout time.Duration, delay Delay) Invoker {
	if tries > MaxTries {
		tries = MaxTries
	}

	if timeout > MaxTimeout {
		timeout = MaxTimeout
	}

	return &FailoverInvoker{
		tries:   tries,
		timeout: timeout,
		delay:   delay,
	}
}

// Invoke method invoke grpc.Invoke that can be retry when invoke failed
func (f *FailoverInvoker) Invoke(ctx context.Context, conn *grpc.ClientConn, method string, request, response interface{}, opts ...grpc.CallOption) error {
	var err error

	for i := 0; i < f.tries; i++ {
		err = grpc.Invoke(ctx, method, request, response, conn, opts...)
		if err != nil {
			delay := f.delay.GetDelay()
			time.Sleep(delay)
			fmt.Println(i, delay)
			continue
		}

		return nil
	}

	return err
}
