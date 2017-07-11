package coordinator

import (
    "golang.org/x/net/context"

    "github.com/servicekit/servicekit-go/spec"
)

type Coordinator interface {
    GetServices(ctx context.Context, name string, tag string) ([]*spec.Service, interface{}, error)
    Register(ctx context.Context, serv *spec.Service) error
    Deregister(ctx context.Context, serviceID string) error
}
