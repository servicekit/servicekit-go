package consul

import (
    "fmt"
    "time"

    "github.com/hashicorp/consul/api"
    "golang.org/x/net/context"

    "github.com/servicekit/servicekit-go/logger"
    "github.com/servicekit/servicekit-go/spec"
    "github.com/servicekit/servicekit-go/version"
)

const (
    DefaultTTL = time.Minute
    EnableTLS  = true
)

type consul struct {
    c *api.Client

    log *logger.Logger
}

func checkService(checkID string, checks []*api.HealthCheck) bool {
    for _, c := range checks {
        if c.CheckID == checkID && c.Status == api.HealthPassing {
            return true
        }
    }

    return false
}

func NewConsul(addr, scheme, token string, log *logger.Logger) (*consul, error) {
    // create a reusable client
    c, err := api.NewClient(&api.Config{Address: addr, Scheme: scheme, Token: token})
    if err != nil {
        return nil, err
    }

    return &consul{
        c: c,

        log: log,
    }, nil
}

func (c *consul) GetServices(ctx context.Context, name string, tag string) ([]*spec.Service, interface{}, error) {
    var passingOnly bool
    var queryOptions *api.QueryOptions

    if v, ok := ctx.Value("passingOnly").(bool); ok == false {
        passingOnly = true
    } else {
        passingOnly = v
    }

    if v, ok := ctx.Value("queryOptions").(*api.QueryOptions); ok == false {
        queryOptions = nil
    } else {
        queryOptions = v
    }

    serviceEntries, meta, err := c.c.Health().Service(name, tag, passingOnly, queryOptions)
    if err != nil {
        return nil, nil, err
    }

    services := make([]*spec.Service, 0)

    for _, serviceEntry := range serviceEntries {
        if checkService(fmt.Sprintf("service:%s", serviceEntry.Service.ID), serviceEntry.Checks) == false {
            continue
        }

        services = append(services, &spec.Service{
            ID:          serviceEntry.Service.ID,
            Service:     serviceEntry.Service.Service,
            Tags:        serviceEntry.Service.Tags,
            Version:     version.GetVersion(serviceEntry.Service.Tags),
            Address:     serviceEntry.Service.Address,
            Port:        serviceEntry.Service.Port,
            CreateIndex: serviceEntry.Service.CreateIndex,
            ModifyIndex: serviceEntry.Service.ModifyIndex,
            NodeID:      serviceEntry.Node.ID,
            Node:        serviceEntry.Node.Node,
            NodeAddress: serviceEntry.Node.Address,
            Datacenter:  serviceEntry.Node.Datacenter,
        })
    }

    return services, meta, nil
}

func (c *consul) Register(ctx context.Context, serv *spec.Service) error {
    ttl := serv.TTL

    enableTLS, ok := ctx.Value("enabletls").(bool)
    if ok != true {
        enableTLS = EnableTLS
    }

    service := &api.AgentServiceRegistration{
        ID:      serv.ID,
        Name:    serv.Service,
        Address: serv.Address,
        Port:    serv.Port,
        Tags:    serv.Tags,
        Check: &api.AgentServiceCheck{
            TTL:           ttl.String(),
            TLSSkipVerify: enableTLS,
        },
    }

    if err := c.c.Agent().ServiceRegister(service); err != nil {
        return err
    } else {
        go func(ctx context.Context) {
            c.log.Infof("consul: service: %s update ttl started", serv.ID)
            for {
                select {
                case <-ctx.Done():
                    c.log.Infof("consul: service: %s update ttl stopped", serv.ID)
                    return
                default:
                    c.log.Infof("consul: service: %s updated ttl ", serv.ID)
                    c.c.Agent().UpdateTTL(fmt.Sprintf("service:%s", serv.ID), "", api.HealthPassing)
                    time.Sleep(ttl/2 - 1)
                }
            }
        }(ctx)
    }

    return nil
}

func (c *consul) Deregister(ctx context.Context, serviceID string) error {
    return c.c.Agent().ServiceDeregister(serviceID)
}
