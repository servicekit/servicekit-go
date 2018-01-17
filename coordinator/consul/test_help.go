package consul

import (
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/consul/api"
	"github.com/servicekit/servicekit-go/spec"
)

// TestConsul is a help stub for test
type TestConsul struct {
	GetServicesServices []*spec.Service
	GetServicesMeta     *api.QueryMeta
	GetServicesError    error
	RegisterError       error
	DeregisterError     error
}

// GetServices returns some service which services do we want to return
func (t *TestConsul) GetServices(ctx context.Context, name string, tag string) ([]*spec.Service, interface{}, error) {
	return t.GetServicesServices, t.GetServicesMeta, t.GetServicesError
}

// Register register a service which service do we want to register
func (t *TestConsul) Register(ctx context.Context, serv *spec.Service, ttl time.Duration) error {
	return t.RegisterError
}

// Deregister deregister a service which service do we want to deregister
func (t *TestConsul) Deregister(ctx context.Context, serviceID string) error {
	return t.DeregisterError
}
