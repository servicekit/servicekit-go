package consul

import (
	"github.com/hashicorp/consul/api"

	"github.com/servicekit/servicekit-go/service"
)

// see github.com/hashicorp/consul/api/agent.go
// type AgentService struct {
//  ID                string
//  Service           string
//  Tags              []string
//  Port              int
//  Address           string
//  EnableTagOverride bool
//  CreateIndex       uint64
//  ModifyIndex       uint64
// }
type ConsulBackend api.AgentService

func (c *ConsulBackend) GetID() string {
	return c.ID
}

func (c *ConsulBackend) GetBackendInfo() *service.BackendInfo {
	return &service.BackendInfo{
		ID:      c.ID,
		Service: c.Service,
		Tags:    c.Tags,
		Port:    c.Port,
		Address: c.Address,
		Status:  service.BackendStatusPassing,
	}

}
