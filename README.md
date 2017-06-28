# servicekit for golang

[![Travis CI Status](https://travis-ci.org/servicekit/fabio.svg?branch=master)](https://travis-ci.org/servicekit/servicekit-go)
[![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/fabiolb/fabio/master/LICENSE)

servicekit is a simple, minimum viable service management library for
golang.

## Quick Start

* get library
  go get github.com/servicekit/servicekit-go              (>=go1.8)

* coordinator service
   ```
	co, err := coordinator.NewConsul("127.0.0.1:8500", "http", "", log)
	if err != nil {
		panic(err)
	}

	// Get Services
	services, meta, error:= co.GetServices(context.Background(), "service name", "tag string")
	...

	// Register service
	err := co.Register(context.Background(), service.Service{ID: "nginx1", Service: "nginx"})
	...

	// Deregister service
	err := co.Deregister(context.Background(), "service_id")
	...

   ```

* grpc load balance
   grpc load balance is implemented by https://github.com/grpc/grpc/blob/master/doc/load-balancing.md
   ```
	co, err := coordinator.NewConsul("127.0.0.1:8500", "http", "", log)
	if err != nil {
		panic(err)
	}

	b, err := balancer.NewConsulBalancer(co, "account_service", "", log)
	if err != nil {
		panic(err)
	}

	r := b.GetResolver(context.Background())

	conn, err := grpc.Dial(
		"",
		grpc.WithBalancer(grpc.RoundRobin(r)))
	...

* grpc invoke fault tolerant
  ```
	 i := invoker.NewFailoverInvoker(10, time.Second, invoker.NewFibDelay(time.Second))
	 err := i.Invoke(context.Background(), conn, "/Account/Auth", authRequest, authResponse)
	 ...
```
