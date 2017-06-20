package consul

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/servicekit/servicekit-go/service"
)

func TestNewConsulServices(t *testing.T) {
	NewConsulServices(nil, make([]string, 0))
}

func getEntries() []*api.ServiceEntry {
	entries := make([]*api.ServiceEntry, 0)
	entries = append(entries, &api.ServiceEntry{
		Service: &api.AgentService{
			ID:      "id1",
			Service: "service1",
			Tags:    make([]string, 0),
			Port:    80,
			Address: "address1",
		},
	})

	entries = append(entries, &api.ServiceEntry{
		Service: &api.AgentService{
			ID:      "id2",
			Service: "service2",
			Tags:    make([]string, 0),
			Port:    81,
			Address: "address2",
		},
	})

	entries = append(entries, &api.ServiceEntry{
		Service: &api.AgentService{
			ID:      "id3",
			Service: "service1",
			Tags:    make([]string, 0),
			Port:    81,
			Address: "address3",
		},
	})

	return entries
}

func TestConsulServicesParseEntry(t *testing.T) {
	c := NewConsulServices(nil, make([]string, 0))

	res := c.(*ConsulServices).parseEntry(getEntries())

	if len(res) != 3 || res[0].GetID() != "id1" || res[1].GetID() != "id2" {
		t.Fatalf("parseEntry error. %v", res[0])
	}

}

func TestConsulServicesUpdateBackendSucceed(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services: services,
		service:  getEntries(),
	}

	res, _ := consul.GetServices()

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err := c.(*ConsulServices).updateBackend(res.([]string))
	if err != nil {
		t.Fatalf(err.Error())
	}

	backends := c.GetBackends("service1")
	if len(backends) != 2 {
		t.Fatalf("backends length should be 2. got: ", len(backends))
	}
	if backends[0].GetID() != "id1" || backends[1].GetID() != "id3" {
		t.Fatalf("backends[0] should be id1 backends[1] should be id3")
	}

	backends = c.GetBackends("service2")
	if len(backends) != 1 {
		t.Fatalf("backends length should be 1. got: ", len(backends))
	}
	if backends[0].GetID() != "id2" {
		t.Fatalf("backends[1] should be id2")
	}
}

func TestConsulServicesUpdateBackendHasError(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services:        services,
		service:         getEntries(),
		getServiceError: fmt.Errorf("get service error"),
	}

	res, _ := consul.GetServices()

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err := c.(*ConsulServices).updateBackend(res.([]string))
	if err.Error() != "get service error" {
		t.Fatalf("updateBackend should be returns error")
	}
}

func TestConsulServicesUpdateSucceed(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services: services,
		service:  getEntries(),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	c.(*ConsulServices).update()

	backends := c.GetBackends("service1")
	if len(backends) != 2 {
		t.Fatalf("backends length should be 2. got: ", len(backends))
	}
	if backends[0].GetID() != "id1" || backends[1].GetID() != "id3" {
		t.Fatalf("backends[0] should be id1 backends[1] should be id3")
	}

	backends = c.GetBackends("service2")
	if len(backends) != 1 {
		t.Fatalf("backends length should be 1. got: ", len(backends))
	}
	if backends[0].GetID() != "id2" {
		t.Fatalf("backends[1] should be id2")
	}

}

func TestConsulServicesUpdateHasGetServicesError(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services:         services,
		service:          getEntries(),
		getServicesError: fmt.Errorf("get services error"),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err := c.(*ConsulServices).update()
	if err.Error() != "get services error" {
		t.Fatalf("update should be return error")
	}
}

func TestConsulServicesUpdateHasUpdateBackendError(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services:        services,
		service:         getEntries(),
		getServiceError: fmt.Errorf("get service error"),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err := c.(*ConsulServices).update()
	if err.Error() != "get service error" {
		t.Fatalf("err should be 'get service error'")
	}

}

func TestConsulServicesGetServicesInfo(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services: services,
		service:  getEntries(),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	c.(*ConsulServices).update()

	servicesInfo := c.GetServicesInfo()

	if _, ok := servicesInfo["service1"]; ok == false {
		t.Fatalf("service1 should be in services")
	}

	if _, ok := servicesInfo["service2"]; ok == false {
		t.Fatalf("service2 should be in services")
	}

	if servicesInfo["service1"].Name != "service1" {
		t.Fatalf("service info name should be service1")
	}

	if servicesInfo["service2"].Name != "service2" {
		t.Fatalf("service info name should be service1")
	}

	if len(servicesInfo["service1"].BackendInfos) != 2 {
		t.Fatalf("service1 backendinfos length should be 2")
	}

	if len(servicesInfo["service2"].BackendInfos) != 1 {
		t.Fatalf("service2 backendinfos length should be 1")
	}

	if servicesInfo["service1"].BackendInfos[0].ID != "id1" {
		t.Fatalf("service1 backendinfos[0].ID should be 'id1'")
	}
	if servicesInfo["service1"].BackendInfos[1].ID != "id3" {
		t.Fatalf("service1 backendinfos[1].ID should be 'id3'")
	}
	if servicesInfo["service2"].BackendInfos[0].ID != "id2" {
		t.Fatalf("service2 backendinfos[0].ID should be 'id2'")
	}

}

func TestConsulServicesGetBackends(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services: services,
		service:  getEntries(),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	c.(*ConsulServices).update()

	res := c.GetBackends("service1")

	if len(res) != 2 {
		t.Fatalf("service1 length should be 2")
	}

	if res[0].GetID() != "id1" || res[1].GetID() != "id3" {
		t.Fatalf("res[0].GetID() should be 'id1', res[1].GetID() should be 'id3'")
	}

	res = c.GetBackends("service2")

	if len(res) != 1 {
		t.Fatalf("service2 length should be 1")
	}

	if res[0].GetID() != "id2" {
		t.Fatalf("res[0].GetID() should be 'id2'")
	}

	res = c.GetBackends("service3")

	if res != nil {
		t.Fatalf("res should be nil")
	}

}

func TestConsulServicesGetBackend(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services: services,
		service:  getEntries(),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	c.(*ConsulServices).update()

	res := c.GetBackend("service1")

	if res == nil {
		t.Fatalf("GetBackend should be not return nil")
	}

	if res.GetID() != "id1" && res.GetID() != "id3" {
		t.Fatalf("res should be id1 or id3")
	}

	res = c.GetBackend("service3")

	if res != nil {
		t.Fatalf("res should be nil")
	}

}

func TestConsulServicesRegister(t *testing.T) {
	consul := &TestConsul{
		registerError: fmt.Errorf("error"),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err := c.Register(service.BackendInfo{}, &ConsulCheck{})

	if err.Error() != "error" {
		t.Fatalf("err should be error")
	}

	consul = &TestConsul{}

	c = NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err = c.Register(service.BackendInfo{}, &ConsulCheck{})

	if err != nil {
		t.Fatalf("err should be nil")
	}

}

func TestConsulServicesDeregister(t *testing.T) {
	consul := &TestConsul{
		deregisterError: fmt.Errorf("error"),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err := c.Deregister("")

	if err.Error() != "error" {
		t.Fatalf("err should be error")
	}

	consul = &TestConsul{}

	c = NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	err = c.Deregister("")

	if err != nil {
		t.Fatalf("err should be nil")
	}

}

func TestConsulServicesStart(t *testing.T) {
	services := make([]string, 2)
	services[0] = "service1"
	services[1] = "service2"

	consul := &TestConsul{
		services: services,
		service:  getEntries(),
	}

	c := NewConsulServices(nil, make([]string, 0))
	c.(*ConsulServices).consul = consul

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c.(*ConsulServices).start(ctx, time.Second)

	ctx, cancel = context.WithCancel(context.Background())
	go func(c service.Services) {
		for {
			res := c.GetBackend("service1")
			if res != nil {
				cancel()
			}
		}
	}(c)
	c.(*ConsulServices).start(ctx, time.Second)
}

func TestNewCOnsulService(t *testing.T) {
	NewConsulService("test")
}

func TestConsulServiceGetName(t *testing.T) {
	c := NewConsulService("test")
	res := c.GetName()

	if res != "test" {
		t.Fatalf("res should be test")
	}
}

func TestConsulServiceGetByRoundRobin(t *testing.T) {
	c := NewConsulService("test")

	backends := make([]service.Backend, 2)
	backends[0] = &ConsulBackend{
		ID: "test1",
	}
	backends[1] = &ConsulBackend{
		ID: "test2",
	}

	res := c.(*ConsulService).getByRoundRobin()

	if res != nil {
		t.Fatalf("res should be nil")
	}

	c.(*ConsulService).backends = backends

	res = c.(*ConsulService).getByRoundRobin()

	if res.GetID() != "test1" {
		t.Fatalf("res.ID should be test1")
	}

	res = c.(*ConsulService).getByRoundRobin()

	if res.GetID() != "test2" {
		t.Fatalf("res.ID should be test2")
	}

	res = c.(*ConsulService).getByRoundRobin()

	if res.GetID() != "test1" {
		t.Fatalf("res.ID should be test1")
	}

}

func TestConsulServiceGetBackends(t *testing.T) {
	c := NewConsulService("test")

	backends := make([]service.Backend, 2)
	backends[0] = &ConsulBackend{
		ID: "test1",
	}
	backends[1] = &ConsulBackend{
		ID: "test2",
	}

	c.(*ConsulService).backends = backends

	res := c.GetBackends()

	if len(res) != 2 {
		t.Fatalf("res length should be 2")
	}

	if res[0].GetID() != "test1" {
		t.Fatalf("res[0].GetID() should be test1")
	}

	if res[1].GetID() != "test2" {
		t.Fatalf("res[1].GetID() should be test2")
	}

}

func TestConsulServiceGetBackend(t *testing.T) {
	c := NewConsulService("test")

	backends := make([]service.Backend, 2)
	backends[0] = &ConsulBackend{
		ID: "test1",
	}
	backends[1] = &ConsulBackend{
		ID: "test2",
	}

	res := c.GetBackend("roundrobin")

	if res != nil {
		t.Fatalf("res should be nil")
	}

	c.(*ConsulService).backends = backends

	res = c.GetBackend("roundrobin")

	if res.GetID() != "test1" {
		t.Fatalf("res.ID should be test1")
	}

	res = c.GetBackend("roundrobin")

	if res.GetID() != "test2" {
		t.Fatalf("res.ID should be test2")
	}

	res = c.GetBackend("roundrobin")

	if res.GetID() != "test1" {
		t.Fatalf("res.ID should be test1")
	}

}

func TestConsulServicePut(t *testing.T) {
	c := NewConsulService("test")

	backends := make([]service.Backend, 2)
	backends[0] = &ConsulBackend{
		ID: "test1",
	}
	backends[1] = &ConsulBackend{
		ID: "test2",
	}

	c.Put(backends)

	res := c.GetBackends()

	if len(res) != 2 {
		t.Fatalf("res length should be 2")
	}

	if res[0].GetID() != "test1" {
		t.Fatalf("res[0].GetID() should be test1")
	}

	if res[1].GetID() != "test2" {
		t.Fatalf("res[1].GetID() should be test2")
	}

	if len(c.(*ConsulService).blackHouse) != 0 {
		t.Fatalf("blackHouse length should be 0")
	}

}
