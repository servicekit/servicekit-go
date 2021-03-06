package consul

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/consul/api"
	"github.com/servicekit/servicekit-go/coordinator"
	"github.com/servicekit/servicekit-go/logger"
	"google.golang.org/grpc/naming"
)

type contextKey string

func (c contextKey) String() string {
	return "mypackage context key " + string(c)
}

// Resolver implements the gRPC Resolver interface using a Consul backend.
//
// See the gRPC load balancing documentation for details about Balancer and
// Resolver: https://github.com/grpc/grpc/blob/master/doc/load-balancing.md.
type Resolver struct {
	consul      coordinator.Coordinator
	service     string
	tag         string
	passingOnly bool

	quitc    chan struct{}
	updatesc chan []*naming.Update

	log *logger.Logger
}

// NewResolver initializes and returns a new Resolver.
//
// It resolves addresses for gRPC connections to the given service and tag.
// If the tag is irrelevant, use an empty string.
func newResolver(consul coordinator.Coordinator, service, tag string, log *logger.Logger) (*Resolver, error) {
	r := &Resolver{
		consul:      consul,
		service:     service,
		tag:         tag,
		passingOnly: true,
		quitc:       make(chan struct{}),
		updatesc:    make(chan []*naming.Update, 1),

		log: log,
	}

	// Retrieve instances immediately
	instances, index, err := r.getInstances(0)
	if err != nil {
		r.log.Warnf("Resolver: error retrieving instances from Consul: %v", err)
	}
	fmt.Println(instances)
	updates := r.makeUpdates(nil, instances)
	if len(updates) > 0 {
		r.updatesc <- updates
	}

	// Start updater
	go r.updater(instances, index)

	return r, nil
}

// Resolve creates a watcher for target. The watcher interface is implemented
// by Resolver as well, see Next and Close.
func (r *Resolver) Resolve(target string) (naming.Watcher, error) {
	return r, nil
}

// Next blocks until an update or error happens. It may return one or more
// updates. The first call will return the full set of instances available
// as NewConsulResolver will look those up. Subsequent calls to Next() will
// block until the resolver finds any new or removed instance.
//
// An error is returned if and only if the watcher cannot recover.
func (r *Resolver) Next() ([]*naming.Update, error) {
	return <-r.updatesc, nil
}

// Close closes the watcher.
func (r *Resolver) Close() {
	select {
	case <-r.quitc:
	default:
		close(r.quitc)
		close(r.updatesc)
	}
}

// updater is a background process started in NewResolver. It takes
// a list of previously resolved instances (in the format of host:port, e.g.
// 192.168.0.1:1234) and the last index returned from Consul.
func (r *Resolver) updater(instances []string, lastIndex uint64) {
	var err error
	var oldInstances = instances
	var newInstances []string

	// TODO Cache the updates for a while, so that we don't overwhelm Consul.
	for {
		select {
		case <-r.quitc:
			break
		default:
			newInstances, lastIndex, err = r.getInstances(lastIndex)
			if err != nil {
				r.log.Debugf("grpc/lb: error retrieving instances from Consul: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			updates := r.makeUpdates(oldInstances, newInstances)
			if len(updates) > 0 {
				r.updatesc <- updates
			}
			oldInstances = newInstances
			time.Sleep(1 * time.Second)
		}
	}
}

// getInstances retrieves the new set of instances registered for the
// service from Consul.
func (r *Resolver) getInstances(lastIndex uint64) ([]string, uint64, error) {
	ctx := context.Background()

	context.WithValue(ctx, contextKey("passingOnly"), true)
	context.WithValue(ctx, contextKey("queryOptions"), &api.QueryOptions{
		WaitIndex: lastIndex,
	})
	services, meta, err := r.consul.GetServices(ctx, r.service, r.tag)
	if err != nil {
		return nil, lastIndex, err
	}

	_m, ok := meta.(*api.QueryMeta)
	if ok == false {
		return nil, lastIndex, fmt.Errorf("invalid meta data")
	}

	var instances []string
	for _, service := range services {
		s := service.Address
		if len(s) == 0 {
			s = service.NodeAddress
		}
		addr := net.JoinHostPort(s, strconv.Itoa(service.Port))
		instances = append(instances, addr)
	}
	return instances, _m.LastIndex, nil
}

// makeUpdates calculates the difference between and old and a new set of
// instances and turns it into an array of naming.Updates.
func (r *Resolver) makeUpdates(oldInstances, newInstances []string) []*naming.Update {
	oldAddr := make(map[string]struct{}, len(oldInstances))
	for _, instance := range oldInstances {
		oldAddr[instance] = struct{}{}
	}
	newAddr := make(map[string]struct{}, len(newInstances))
	for _, instance := range newInstances {
		newAddr[instance] = struct{}{}
	}

	var updates []*naming.Update
	for addr := range newAddr {
		if _, ok := oldAddr[addr]; !ok {
			updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr})
		}
	}
	for addr := range oldAddr {
		if _, ok := newAddr[addr]; !ok {
			updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr})
		}
	}

	return updates

}
