package grpchelper

import (
    "fmt"

    "github.com/rs/xid"
    "golang.org/x/net/context"
    "google.golang.org/grpc/metadata"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "x-request-id"

func HandleRequestID(ctx context.Context) string {
    md, ok := metadata.FromContext(ctx)
    if !ok {
        return newRequestID()
    }

    header, ok := md[DefaultXRequestIDKey]
    if !ok || len(header) == 0 {
        return newRequestID()
    }

    requestID := header[0]
    if requestID == "" {
        return newRequestID()
    }

    return requestID
}

func HandleRequestIDChain(ctx context.Context) string {
    md, ok := metadata.FromContext(ctx)
    if !ok {
        return newRequestID()
    }

    header, ok := md[DefaultXRequestIDKey]
    if !ok || len(header) == 0 {
        return newRequestID()
    }

    requestID := header[0]
    if requestID == "" {
        return newRequestID()
    }

    return fmt.Sprintf("%s,%s", requestID, newRequestID())
}

func newRequestID() string {
    return xid.New().String()
}

func UpdateContextWithRequestID(ctx context.Context, requestID string) context.Context {
    md := metadata.New(map[string]string{DefaultXRequestIDKey: requestID})
    _md, ok := metadata.FromContext(ctx)
    if ok {
        md = metadata.Join(_md, md)
    }

    ctx = metadata.NewContext(ctx, md)

    return ctx
}

func GetRequestID(ctx context.Context) string {
    md, ok := metadata.FromContext(ctx)
    if ok == false {
        return ""
    }

    header, ok := md[DefaultXRequestIDKey]
    if !ok || len(header) == 0 {
        return ""
    }

    return header[0]
}
