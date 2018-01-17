package requestid

import (
	"fmt"
	"net/http"

	"github.com/rs/xid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

type contextKey string

func (c contextKey) String() string {
	return "mypackage context key " + string(c)
}

// RequestIDKey is metadata key name for request ID
var RequestIDKey = "request-id"

// HandleRequestID got requestid from context
// If no requestid in context, create a new requestid
func HandleRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newRequestID()
	}

	header, ok := md[RequestIDKey]
	if !ok || len(header) == 0 {
		return newRequestID()
	}

	requestID := header[0]
	if requestID == "" {
		return newRequestID()
	}

	return requestID
}

// HandleRequestIDChain got old requestid(got from context or new) and new requestid
func HandleRequestIDChain(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return newRequestID()
	}

	header, ok := md[RequestIDKey]
	if !ok || len(header) == 0 {
		return newRequestID()
	}

	requestID := header[0]
	if requestID == "" {
		return newRequestID()
	}

	return fmt.Sprintf("%s,%s", requestID, newRequestID())
}

// newRequestID generates a requestid
func newRequestID() string {
	return xid.New().String()
}

// UpdateContextWithRequestID set a requestID to context
func UpdateContextWithRequestID(ctx context.Context, requestID string) context.Context {
	md := metadata.New(map[string]string{RequestIDKey: requestID})
	_md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md = metadata.Join(_md, md)
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
	ctx = context.WithValue(ctx, contextKey(RequestIDKey), requestID)
	return ctx
}

// GetRequestID got requestid from context
func GetRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok == false {
		return ""
	}

	header, ok := md[RequestIDKey]
	if !ok || len(header) == 0 {
		return ""
	}

	return header[0]
}

// GetRequestIDFromHTTPRequest got requestid from http request
func GetRequestIDFromHTTPRequest(ctx context.Context, r *http.Request) (context.Context, string) {
	requestID := r.Header.Get(RequestIDKey)
	if requestID == "" {
		requestID = HandleRequestID(ctx)
	}

	ctx = context.WithValue(ctx, contextKey(RequestIDKey), requestID)
	return ctx, requestID
}
