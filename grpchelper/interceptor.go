package grpchelper

import (
	"reflect"
	"runtime"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/servicekit/servicekit-go/logger"
)

type requestIDKey struct{}

func UnaryServerChan(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	switch len(interceptors) {
	case 0:
		// do not want to return nil interceptor since this function was never defined to do so/for backwards compatibility
		return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			requestID := HandleRequestIDChain(ctx)
			ctx = context.WithValue(ctx, requestIDKey{}, requestID)
			return handler(ctx, req)
		}
	case 1:
		return func(ctx context.Context, req interface{}, unaryServerInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			requestID := HandleRequestIDChain(ctx)
			ctx = context.WithValue(ctx, requestIDKey{}, requestID)
			return interceptors[0](ctx, req, unaryServerInfo, handler)
		}
	default:
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			buildChain := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
				return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
					return current(currentCtx, currentReq, info, next)
				}
			}
			chain := handler
			for i := len(interceptors) - 1; i >= 0; i-- {
				chain = buildChain(interceptors[i], chain)
			}
			requestID := HandleRequestIDChain(ctx)
			ctx = UpdateContextWithRequestID(ctx, requestID)
			return chain(ctx, req)
		}
	}
}

func UnaryClientChan(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	switch len(interceptors) {
	case 0:
		// do not want to return nil interceptor since this function was never defined to do so/for backwards compatibility
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			requestID := HandleRequestIDChain(ctx)
			ctx = context.WithValue(ctx, requestIDKey{}, requestID)
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	case 1:
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			requestID := HandleRequestIDChain(ctx)
			ctx = context.WithValue(ctx, requestIDKey{}, requestID)
			return interceptors[0](ctx, method, req, reply, cc, invoker, opts...)
		}
	default:
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			buildChain := func(current grpc.UnaryClientInterceptor, next grpc.UnaryInvoker) grpc.UnaryHandler {
				return func(currentCrx context.Context, currentMethod string, currentReq, currentReply interface{}, currentClientConn *grpc.ClientConn, currentOpts ...grpc.CallOption) error {
					return current(currentCtx, currentMethod, currentReq, currentReply, currentClientConn, currentOpts)
				}
			}
			chain := invoker
			for i := len(interceptors) - 1; i >= 0; i-- {
				chain = buildChain(interceptors[i], chain)
			}
			requestID := HandleRequestIDChain(ctx)
			ctx = UpdateContextWithRequestID(ctx, requestID)
			return chain(ctx, method, req, reply, cc, opts...)
		}
	}
}

// UnaryClientInterceptor returns a new retrying unary client interceptor.
//
// The default configuration of the interceptor is to not retry *at all*. This behaviour can be
// changed through options (e.g. WithMax) on creation of the interceptor or on call (through grpc.CallOptions).
func UnaryClientInterceptor(optFuncs ...CallOption) grpc.UnaryClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultOptions, optFuncs)
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		// short circuit for simplicity, and avoiding allocations.
		if callOpts.max == 0 {
			return invoker(parentCtx, method, req, reply, cc, grpcOpts...)
		}
		var lastErr error
		for attempt := uint(0); attempt < callOpts.max; attempt++ {
			if err := waitRetryBackoff(attempt, parentCtx, callOpts); err != nil {
				return err
			}
			callCtx := perCallContext(parentCtx, callOpts, attempt)
			lastErr = invoker(callCtx, method, req, reply, cc, grpcOpts...)
			// TODO(mwitkow): Maybe dial and transport errors should be retriable?
			if lastErr == nil {
				return nil
			}
			logTrace(parentCtx, "grpc_retry attempt: %d, got err: %v", attempt, lastErr)
			if isContextError(lastErr) {
				if parentCtx.Err() != nil {
					logTrace(parentCtx, "grpc_retry attempt: %d, parent context error: %v", attempt, parentCtx.Err())
					// its the parent context deadline or cancellation.
					return lastErr
				} else {
					logTrace(parentCtx, "grpc_retry attempt: %d, context error from retry call", attempt)
					// its the callCtx deadline or cancellation, in which case try again.
					continue
				}
			}
			if !isRetriable(lastErr, callOpts) {
				return lastErr
			}
		}
		return lastErr
	}
}

type commonUnaryServerInterceptor struct {
	log *logger.Logger
}

func NewCommonUnaryServerInterceptor(log *logger.Logger) *commonUnaryServerInterceptor {
	return &commonUnaryServerInterceptor{
		log: log,
	}
}

func (i *commonUnaryServerInterceptor) RecoverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]
			i.log.Errorf("panic grpc invoke: %s, err=%v, stack:\n%s", info.FullMethod, r, string(stack))
		}
	}()

	return handler(ctx, req)
}

func (i *commonUnaryServerInterceptor) TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()

	rid := reflect.ValueOf(req).Elem().FieldByName("RequestID")

	i.log.Infof("GRPC Request: %v start. RequestID: %v", info.FullMethod, rid)

	resp, err = handler(ctx, req)

	doneTime := time.Now().Sub(startTime)

	i.log.Infof("GRPC Request: %v done. RequestID: %v, time: %v", info.FullMethod, rid, doneTime.String())

	return resp, err

}
