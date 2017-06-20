package consul

import (
	"testing"
)

func TestConsulBackendGetID(t *testing.T) {
	c := &ConsulBackend{
		ID: "test",
	}

	if c.GetID() != "test" {
		t.Fatalf("id should be test")
	}
}

func TestConsulBackendGetBackendInfo(t *testing.T) {
	c := &ConsulBackend{
		ID: "test",
	}

	res := c.GetBackendInfo()

	if res.ID != "test" {
		t.Fatalf("res.ID should be test")
	}
}
