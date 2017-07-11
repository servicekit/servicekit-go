package consul

import (
    "golang.org/x/net/context"

    "github.com/hashicorp/consul/api"
    "github.com/servicekit/servicekit-go/spec"
)

type TestConsul struct {
    GetServicesServices []*spec.Service
    GetServicesMeta     *api.QueryMeta
    GetServicesError    error
    RegisterError       error
    DeregisterError     error
}

func (t *TestConsul) GetServices(ctx context.Context, name string, tag string) ([]*spec.Service, interface{}, error) {
    return t.GetServicesServices, t.GetServicesMeta, t.GetServicesError
}

func (t *TestConsul) Register(ctx context.Context, serv *spec.Service) error {
    return t.RegisterError
}

func (t *TestConsul) Deregister(ctx context.Context, serviceID string) error {
    return t.DeregisterError
}
