package consul

import (
	"github.com/hashicorp/consul/api"

	"github.com/servicekit/servicekit-go/service"
)

type IAgent interface {
	ServiceRegister(*api.AgentServiceRegistration) error
	ServiceDeregister(serviceID string) error
}

type ICatalog interface {
	Services(q *api.QueryOptions) (map[string][]string, *api.QueryMeta, error)
}

type IHealth interface {
	Service(service, tag string, passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error)
}

type IClient interface {
	Agent() IAgent
	Catalog() ICatalog
	Health() IHealth
}

type Consul struct {
	c IClient
}

func NewConsul(c IClient) service.ServiceCoordinate {
	// create a reusable client
	// c, err := api.NewClient(&api.Config{Address: cfg.Addr, Scheme: cfg.Scheme, Token: cfg.Token})
	// if err != nil {
	//	return nil, err
	// }

	return &Consul{
		c: c,
	}
}

func (c *Consul) GetServices() (interface{}, error) {
	services, _, err := c.c.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	_services := make([]string, len(services))

	count := 0
	for k, _ := range services {
		_services[count] = k
		count += 1
	}

	return _services, nil
}

func (c *Consul) GetService(name string) (interface{}, error) {
	serviceEntries, _, err := c.c.Health().Service(name, "", true, nil)

	if err != nil {
		// log.Warn(err)
		return nil, err
	}

	return serviceEntries, nil

}

func (c *Consul) Register(info service.BackendInfo, check service.Check) error {
	service := &api.AgentServiceRegistration{
		ID:      info.ID,
		Name:    info.Service,
		Address: info.Address,
		Port:    info.Port,
		Tags:    info.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:          check.GetHTTPURL(),
			Interval:      check.GetInterval().String(),
			Timeout:       check.GetTimeout().String(),
			TLSSkipVerify: check.GetTLSSkipVerify(),
		},
	}

	if err := c.c.Agent().ServiceRegister(service); err != nil {
		return err
	}

	return nil
}

func (c *Consul) Deregister(serviceID string) error {
	return c.c.Agent().ServiceDeregister(serviceID)
}
