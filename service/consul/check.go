package consul

import (
	"time"
)

type ConsulCheck struct {
	http          string
	interval      time.Duration
	timeout       time.Duration
	tlsSkipVerify bool
}

func (c *ConsulCheck) GetHTTPURL() string {
	return c.http
}

func (c *ConsulCheck) GetInterval() time.Duration {
	return c.interval
}

func (c *ConsulCheck) GetTimeout() time.Duration {
	return c.timeout
}

func (c *ConsulCheck) GetTLSSkipVerify() bool {
	return c.tlsSkipVerify
}
