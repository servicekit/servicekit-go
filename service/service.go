package service

import (
	"time"
)

type Services interface {
	GetServicesInfo() map[string]*ServiceInfo
	GetBackends(name string) []Backend
	GetBackend(name string) Backend
	Register(info BackendInfo, check Check) error
	Deregister(serviceID string) error
}

type Service interface {
	GetName() string
	GetBackend(Policy string, opt ...interface{}) Backend
	GetBackends() []Backend
	Put(backends []Backend) error
}

type Backend interface {
	GetID() string
	GetBackendInfo() *BackendInfo
}

type ServiceCoordinate interface {
	GetServices() (interface{}, error)
	GetService(name string) (interface{}, error)
	Register(info BackendInfo, check Check) error
	Deregister(serviceID string) error
}

type Check interface {
	GetHTTPURL() string
	GetInterval() time.Duration
	GetTimeout() time.Duration
	GetTLSSkipVerify() bool
}
