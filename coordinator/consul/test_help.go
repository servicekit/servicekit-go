package consul

import (
	"golang.org/x/net/context"

	"github.com/hashicorp/consul/api"
	"github.com/servicekit/servicekit-go/service"
)

type TestConsul struct {
	GetServicesServices []*service.Service
	GetServicesMeta     *api.QueryMeta
	GetServicesError    error
	RegisterError       error
	DeregisterError     error
}

func (t *TestConsul) GetServices(ctx context.Context, name string, tag string) ([]*service.Service, interface{}, error) {
	return t.GetServicesServices, t.GetServicesMeta, t.GetServicesError
}

func (t *TestConsul) Register(ctx context.Context, serv *service.Service) error {
	return t.RegisterError
}

func (t *TestConsul) Deregister(ctx context.Context, serviceID string) error {
	return t.DeregisterError
}
