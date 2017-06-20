package consul

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"

	"github.com/servicekit/servicekit-go/service"
)

func TestConsulNewConsul(t *testing.T) {
	tc := &TestClient{}

	c := NewConsul(tc)

	if c.(*Consul).c != tc {
		t.Fatalf("Consul.c should be TestClient")
	}
}

func TestConsulGetServices(t *testing.T) {
	tc := &TestClient{
		catalog: &TestCatalog{
			err: fmt.Errorf("error"),
		},
	}

	c := NewConsul(tc)

	_, err := c.GetServices()
	if err.Error() != "error" {
		t.Fatalf("err.Error() should be error")
	}

	services := make(map[string][]string)
	services["service1"] = []string{"tag1", "tag2"}
	services["service2"] = []string{"tag1", "tag2"}
	queryMeta := &api.QueryMeta{}

	tc = &TestClient{
		catalog: &TestCatalog{
			services:  services,
			queryMeta: queryMeta,
		},
	}

	c = NewConsul(tc)

	res, _ := c.GetServices()

	if len(res.([]string)) != 2 {
		t.Fatalf("res length should be 2")
	}

	if res.([]string)[0] != "service1" && res.([]string)[1] != "service1" {
		t.Fatalf("service1 should be in res")
	}

	if res.([]string)[0] != "service2" && res.([]string)[1] != "service2" {
		t.Fatalf("service2 should be in res")
	}
}

func TestConsulGetService(t *testing.T) {
	tc := &TestClient{
		health: &TestHealth{
			err: fmt.Errorf("error"),
		},
	}

	c := NewConsul(tc)

	_, err := c.GetService("")
	if err.Error() != "error" {
		t.Fatalf("err.Error() should be error")
	}

	serviceEntries := []*api.ServiceEntry{&api.ServiceEntry{}, &api.ServiceEntry{}}
	queryMeta := &api.QueryMeta{}

	tc = &TestClient{
		health: &TestHealth{
			serviceEntries: serviceEntries,
			queryMeta:      queryMeta,
		},
	}

	c = NewConsul(tc)

	res, _ := c.GetService("")

	if len(res.([]*api.ServiceEntry)) != 2 {
		t.Fatalf("res length should be 2")
	}

	if res.([]*api.ServiceEntry)[0] != serviceEntries[0] {
		t.Fatalf("res[0] should be serviceEntries[0]")
	}

	if res.([]*api.ServiceEntry)[1] != serviceEntries[1] {
		t.Fatalf("res[1] should be serviceEntries[1]")
	}

}

func TestConsulRegister(t *testing.T) {
	tc := &TestClient{
		agent: &TestAgent{
			registerError: fmt.Errorf("error"),
		},
	}

	c := NewConsul(tc)

	err := c.Register(service.BackendInfo{}, &ConsulCheck{})

	if err.Error() != "error" {
		t.Fatalf("err should be error")
	}

	tc = &TestClient{
		agent: &TestAgent{},
	}

	c = NewConsul(tc)

	err = c.Register(service.BackendInfo{}, &ConsulCheck{})

	if err != nil {
		t.Fatalf("err should be nil")
	}

}

func TestConsulDeregister(t *testing.T) {
	tc := &TestClient{
		agent: &TestAgent{
			deregisterError: fmt.Errorf("error"),
		},
	}

	c := NewConsul(tc)

	err := c.Deregister("")

	if err.Error() != "error" {
		t.Fatalf("err should be error")
	}

	tc = &TestClient{
		agent: &TestAgent{},
	}

	c = NewConsul(tc)

	err = c.Deregister("")

	if err != nil {
		t.Fatalf("err should be nil")
	}

}
