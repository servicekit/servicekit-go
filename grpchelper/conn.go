package grpchelper

import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/naming"
)

func BalanceDial(creds credentials.TransportCredentials, resolver naming.Resolver) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(
        "",
        grpc.WithTransportCredentials(creds),
        grpc.WithBalancer(grpc.RoundRobin(resolver)))
    if err != nil {
        return nil, err
    }

    return conn, nil
}
