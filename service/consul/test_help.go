package consul

import (
	"github.com/hashicorp/consul/api"

	"github.com/servicekit/servicekit-go/service"
)

type TestConsul struct {
	services         []string
	service          []*api.ServiceEntry
	getServicesError error
	getServiceError  error
	registerError    error
	deregisterError  error
}

func (c *TestConsul) GetServices() (interface{}, error) {
	if c.getServicesError != nil {
		return nil, c.getServicesError
	}

	return c.services, nil
}

func (c *TestConsul) GetService(name string) (interface{}, error) {
	if c.getServiceError != nil {
		return nil, c.getServiceError
	}

	entries := make([]*api.ServiceEntry, 0)

	for _, entry := range c.service {
		if entry.Service.Service == name {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

func (c *TestConsul) Register(info service.BackendInfo, check service.Check) error {
	return c.registerError
}

func (c *TestConsul) Deregister(serviceID string) error {
	return c.deregisterError
}

type TestAgent struct {
	registerError   error
	deregisterError error
}

func (t *TestAgent) ServiceRegister(*api.AgentServiceRegistration) error {
	return t.registerError
}

func (t *TestAgent) ServiceDeregister(string) error {
	return t.deregisterError
}

type TestCatalog struct {
	services  map[string][]string
	queryMeta *api.QueryMeta
	err       error
}

func (t *TestCatalog) Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error) {
	return t.services, t.queryMeta, t.err
}

type TestHealth struct {
	serviceEntries []*api.ServiceEntry
	queryMeta      *api.QueryMeta
	err            error
}

func (t *TestHealth) Service(service, tag string, passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	return t.serviceEntries, t.queryMeta, t.err
}

type TestClient struct {
	agent   IAgent
	catalog ICatalog
	health  IHealth
}

func (t *TestClient) Agent() IAgent {
	return t.agent
}

func (t *TestClient) Catalog() ICatalog {
	return t.catalog
}

func (t *TestClient) Health() IHealth {
	return t.health
}
