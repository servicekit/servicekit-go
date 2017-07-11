package service

import (
    "fmt"
    "time"

    "golang.org/x/net/context"

    "github.com/servicekit/servicekit-go/coordinator"
    "github.com/servicekit/servicekit-go/logger"
    "github.com/servicekit/servicekit-go/spec"
)

type GRPCServer interface {
    Serve(ctx context.Context, network, addr, pem, key string) error
}

type GRPCService struct {
    ID      string
    Service string
    Tags    []string
    Address string
    Port    int

    server GRPCServer
    pem    string
    key    string

    TTL time.Duration

    c coordinator.Coordinator

    log *logger.Logger

    errorChan chan error
}

func NewGRPCService(id string, service string, tags []string, address string, port int, server GRPCServer, pem, key string, ttl time.Duration, c coordinator.Coordinator, log *logger.Logger) *GRPCService {
    return &GRPCService{
        ID:      id,
        Service: service,
        Tags:    tags,
        Address: address,
        Port:    port,

        server: server,
        pem:    pem,
        key:    key,

        TTL: ttl,
        c:   c,
        log: log,

        errorChan: make(chan error),
    }
}

func (g *GRPCService) getService() spec.Service {
    return spec.Service{
        ID:      g.ID,
        Service: g.Service,
        Tags:    g.Tags,
        Address: g.Address,
        Port:    g.Port,
    }
}

func (g *GRPCService) Start(ctx context.Context, delayRegisterTime time.Duration) error {
    go func() {
        err := g.server.Serve(
            ctx,
            "tcp",
            fmt.Sprintf("%s:%d", g.Address, g.Port),
            g.pem,
            g.key,
        )
        g.errorChan <- err
    }()

    time.Sleep(delayRegisterTime)

    var err error

    select {
    case err = <-g.errorChan:
    default:
        err = g.c.Register(ctx, g.getService())
    }

    return err
}
