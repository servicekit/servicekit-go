package grpchelper

import (
    "reflect"
    "runtime"
    "time"

    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"

    "github.com/servicekit/servicekit-go/logger"
)

type requestIDKey struct{}

func UnaryServerChan(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
    switch len(interceptors) {
    case 0:
        // do not want to return nil interceptor since this function was never defined to do so/for backwards compatibility
        return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
            requestID := HandleRequestID(ctx)
            ctx = context.WithValue(ctx, requestIDKey{}, requestID)
            return handler(ctx, req)
        }
    case 1:
        return func(ctx context.Context, req interface{}, unaryServerInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
            requestID := HandleRequestID(ctx)
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
            md := metadata.New(map[string]string{DefaultXRequestIDKey: requestID})
            _md, ok := metadata.FromContext(ctx)
            if ok {
                md = metadata.Join(_md, md)
            }
            ctx = metadata.NewIncomingContext(ctx, md)
            return chain(ctx, req)
        }
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
