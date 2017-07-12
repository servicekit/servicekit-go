package config

import (
    "fmt"
    "log/syslog"

    "github.com/kelseyhightower/envconfig"
)

type ServiceENV string
type ServiceVersion string

func NewCustomConfigPrefix(s string) string {
    return fmt.Sprintf("%s_%s", "SERVICE_CUSTOM", s)
}

func newConfigPrefix(s string) string {
    return fmt.Sprintf("%s_%s", "SERVICE", s)
}

const (
    ServiceENVDev     ServiceENV = "dev"
    ServiceENVStaging ServiceENV = "staging"
    ServiceENVProd    ServiceENV = "production"
)

type ServiceConfig struct {
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

func NewServiceConfig() *ServiceConfig {
    var serviceConfig ServiceConfig
    envconfig.Process(newConfigPrefix("SERVICE"), &serviceConfig)
    fmt.Println(serviceConfig)
    return &serviceConfig
}
