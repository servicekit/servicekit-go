package config

import (
	"fmt"
	"log/syslog"

	"github.com/kelseyhightower/envconfig"
)

// ServiceENV string
type ServiceENV string

// ServiceVersion string
type ServiceVersion string

// NewCustomConfigPrefix returns a string which start with "SERVICE_CUSTOM"
func NewCustomConfigPrefix(s string) string {
	return fmt.Sprintf("%s_%s", "SERVICE_CUSTOM", s)
}

// newConfigPrefix returns a string which start with "SERVICE"
func newConfigPrefix(s string) string {
	return fmt.Sprintf("%s_%s", "SERVICE", s)
}

const (
	// ServiceENVTesting used for testing env
	ServiceENVTesting ServiceENV = "testing"
	// ServiceENVDev used for dev env
	ServiceENVDev ServiceENV = "dev"
	// ServiceENVStaging used for staging env
	ServiceENVStaging ServiceENV = "staging"
	// ServiceENVProd used for prod env
	ServiceENVProd ServiceENV = "production"
)

// ServiceConfig is used to describe configs
type ServiceConfig struct {
	ServiceID      string
	ServiceName    string
	ServiceENV     ServiceENV
	ServiceVersion string
	ServiceTags    []string

	TraceHost string
	TracePort int

	LoggerNetwork  string
	LoggerADDR     string
	LoggerPriority syslog.Priority
}

// NewServiceConfig returns a ServiceConfig
func NewServiceConfig() *ServiceConfig {
	var serviceConfig ServiceConfig
	envconfig.Process(newConfigPrefix("SERVICE"), &serviceConfig)
	fmt.Println(serviceConfig)
	return &serviceConfig
}
