package grpchelper

import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"

    balancer "github.com/servicekit/servicekit-go/balancer/consul"
    "github.com/servicekit/servicekit-go/coordinator"
    "github.com/servicekit/servicekit-go/logger"
)

func BalanceDial(credPath string, c coordinator.Coordinator, service string, tag string, log *logger.Logger) (*grpc.ClientConn, error) {
    creds, err := credentials.NewClientTLSFromFile(credPath, service)
    if err != nil {
        return nil, err
    }

    resolver, err := balancer.GetResolver(c, service, tag, log)
    if err != nil {
        return nil, err
    }

    conn, err := grpc.Dial(
        "",
        grpc.WithTransportCredentials(creds),
        grpc.WithBalancer(grpc.RoundRobin(resolver)))
    if err != nil {
        return nil, err
    }

    return conn, nil
}
