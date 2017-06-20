package consul

import (
	"context"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/servicekit/servicekit-go/service"
)

type ConsulServices struct {
	consul   service.ServiceCoordinate
	services map[string]service.Service
	targets  []string
}

func NewConsulServices(consul *Consul, targets []string) service.Services {
	return &ConsulServices{
		consul:   consul,
		services: make(map[string]service.Service),
		targets:  targets,
	}
}

func (c *ConsulServices) parseEntry(entries []*api.ServiceEntry) []service.Backend {
	backends := make([]service.Backend, len(entries))

	for idx, entry := range entries {
		backends[idx] = &ConsulBackend{
			ID:                entry.Service.ID,
			Service:           entry.Service.Service,
			Tags:              entry.Service.Tags,
			Port:              entry.Service.Port,
			Address:           entry.Service.Address,
			EnableTagOverride: entry.Service.EnableTagOverride,
			CreateIndex:       entry.Service.CreateIndex,
			ModifyIndex:       entry.Service.ModifyIndex,
		}
	}

	return backends
}

func (c *ConsulServices) updateBackend(services []string) error {
	for _, s := range services {
		serviceEntries, err := c.consul.GetService(s)
		if err != nil {
			return err
		}

		backends := c.parseEntry(serviceEntries.([]*api.ServiceEntry))

		if _, ok := c.services[s]; ok == false {
			c.services[s] = NewConsulService(s)
		}
		c.services[s].Put(backends)

	}

	return nil
}

func (c *ConsulServices) update() error {
	services, err := c.consul.GetServices()
	if err != nil {
		return err
	}

	err = c.updateBackend(services.([]string))
	if err != nil {
		return err
	}

	return nil

}

func (c *ConsulServices) GetServicesInfo() map[string]*service.ServiceInfo {
	serviceInfos := make(map[string]*service.ServiceInfo)

	for k, v := range c.services {
		backends := v.GetBackends()
		backendInfos := make([]*service.BackendInfo, len(backends))
		for idx, b := range backends {
			backendInfos[idx] = b.GetBackendInfo()
		}

		serviceInfos[k] = &service.ServiceInfo{
			Name:         k,
			BackendInfos: backendInfos,
		}
	}

	return serviceInfos
}

func (c *ConsulServices) GetBackends(name string) []service.Backend {
	if service, ok := c.services[name]; ok == false {
		return nil
	} else {
		return service.GetBackends()
	}
}

func (c *ConsulServices) GetBackend(name string) service.Backend {
	if service, ok := c.services[name]; ok == false {
		return nil
	} else {
		return service.GetBackend("roundrobin")
	}
}

func (c *ConsulServices) Register(info service.BackendInfo, check service.Check) error {
	return c.consul.Register(info, check)
}

func (c *ConsulServices) Deregister(serviceID string) error {
	return c.consul.Deregister(serviceID)
}

func (c *ConsulServices) start(ctx context.Context, interval time.Duration) {
	// log.Infof("services started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.update()
		}

		time.Sleep(interval)
	}

	// log.Infof("services stopped")
}

type ConsulService struct {
	name           string
	backends       []service.Backend
	blackHouse     []service.Backend
	backupBackends []service.Backend
	idx            uint8
	backendsStatus map[string]bool
	sync.RWMutex
}

func NewConsulService(name string) service.Service {
	return &ConsulService{
		name:           name,
		backends:       make([]service.Backend, 0),
		blackHouse:     make([]service.Backend, 0),
		backupBackends: make([]service.Backend, 0),
		backendsStatus: make(map[string]bool),
	}

}

func (s *ConsulService) GetName() string {
	return s.name
}

func (s *ConsulService) getByRoundRobin() service.Backend {
	if len(s.backends) == 0 {
		return nil
	}

	s.RLock()
	backend := s.backends[int(s.idx)%len(s.backends)]
	s.RUnlock()
	s.idx += 1

	return backend
}

func (s *ConsulService) GetBackends() []service.Backend {
	return s.backends
}

func (s *ConsulService) GetBackend(policy string, opt ...interface{}) service.Backend {
	backend := s.getByRoundRobin()
	return backend
}

func (s *ConsulService) Put(backends []service.Backend) error {
	s.Lock()
	defer s.Unlock()

	s.backends = backends
	s.blackHouse = make([]service.Backend, 0)

	return nil
}
