package coordinator

import (
	"golang.org/x/net/context"

	"github.com/servicekit/servicekit-go/service"
)

type Coordinator interface {
	GetServices(ctx context.Context, name string, tag string) ([]*service.Service, interface{}, error)
	Register(ctx context.Context, serv *service.Service) error
	Deregister(ctx context.Context, serviceID string) error
}
