package consul

import (
	"testing"
)

func TestConsulCheckGetHTTPURL(t *testing.T) {
	c := &ConsulCheck{
		http: "test",
	}

	res := c.GetHTTPURL()

	if res != "test" {
		t.Fatalf("res should be test")
	}
}

func TestConsulCheckGetInterval(t *testing.T) {
	c := &ConsulCheck{
		interval: 1,
	}

	res := c.GetInterval()

	if res != 1 {
		t.Fatalf("res should be 1")
	}
}

func TestConsulCheckGetTimeout(t *testing.T) {
	c := &ConsulCheck{
		timeout: 1,
	}

	res := c.GetTimeout()

	if res != 1 {
		t.Fatalf("res should be 1")
	}
}

func TestConsulCheckGetTLSSkipVerify(t *testing.T) {
	c := &ConsulCheck{
		tlsSkipVerify: false,
	}

	res := c.GetTLSSkipVerify()

	if res != false {
		t.Fatalf("res should be false")
	}

	c = &ConsulCheck{
		tlsSkipVerify: true,
	}

	res = c.GetTLSSkipVerify()

	if res != true {
		t.Fatalf("res should be true")
	}
}
