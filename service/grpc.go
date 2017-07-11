package service

import (
    "fmt"
    "time"

    "golang.org/x/net/context"

    "github.com/servicekit/servicekit-go/coordinator"
    "github.com/servicekit/servicekit-go/logger"
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

    TTL time.Duration

    c coordinator.Coordinator

    log *logger.Logger

    errorChan chan error
}

func NewGRPCService(id string, service string, tags []string, address string, port int, server GRPCServer, ttl time.Duration, c coordinator.Coordinator, log *logger.Logger) *GRPCService {
    return &GRPCService{
        ID:      id,
        Service: service,
        Tags:    tags,
        Address: address,
        Port:    port,

        server: server,

        TTL: ttl,
        c:   c,
        log: log,

        errorChan: make(chan error),
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
        err = c.Register(ctx, s)
    }

    return err
}
