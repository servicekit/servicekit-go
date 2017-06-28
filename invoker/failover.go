package invoker

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	MaxTries   = 8
	MaxTimeout = 8 * time.Second
	MaxDelay   = 8 * time.Minute
)

type FailoverInvoker struct {
	tries   int
	timeout time.Duration
	delay   Delay
}

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
