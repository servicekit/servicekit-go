package consul

import (
	"testing"

	"github.com/hashicorp/consul/api"

	coordinator "github.com/servicekit/servicekit-go/coordinator/consul"
	"github.com/servicekit/servicekit-go/logger"
	"github.com/servicekit/servicekit-go/service"
)

func TestNewConsulBalancer(t *testing.T) {

	tc := &coordinator.TestConsul{
		GetServicesServices: make([]*service.Service, 0),
		GetServicesMeta:     &api.QueryMeta{},
		GetServicesError:    nil,
	}

	NewConsulBalancer(tc, "", "", &logger.Logger{})
}
