package service

type BackendStatus string

var BackendStatusPassing BackendStatus = "passing"
var BackendStatusWarning BackendStatus = "warning"
var BackendStatusCritical BackendStatus = "critical"

type ServiceInfo struct {
	Name         string         `json: "name"`
	BackendInfos []*BackendInfo `json: "backends"`
}

type BackendInfo struct {
	ID      string        `json: "id"`
	Service string        `json: "service"`
	Tags    []string      `json: "tags"`
	Port    int           `json: "port"`
	Address string        `json: "string"`
	Status  BackendStatus `json: "status"`
}
