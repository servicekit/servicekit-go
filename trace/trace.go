package trace

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/servicekit/servicekit-go/coordinator"
	"github.com/servicekit/servicekit-go/logger"
	"github.com/servicekit/servicekit-go/spec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Trace can init a Prometheus handler for provider Prometheus vectors
type Trace struct {
	addr string

	prom *prom

	log *logger.Logger
}

// NewTrace returns a trace
// Serve a http server that provider a http interface for push metrics
// Register itself to the coordinator
func NewTrace(c coordinator.Coordinator, id string, name string, tags []string, host string, port int, ttl time.Duration, log *logger.Logger) (*Trace, error) {
	t := &Trace{
		addr: fmt.Sprintf("%s:%d", host, port),
	}

	http.Handle("/metrics", promhttp.Handler())

	t.prom = &prom{
		path:       "/metrics",
		collectors: make(map[string]prometheus.Collector),
		log:        log,
	}

	s := &spec.Service{
		ID:      id,
		Service: name,
		Tags:    tags,
		Address: host,
		Port:    port,
	}

	err := c.Register(context.Background(), s, ttl)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Serve serve a http server
func (h *Trace) Serve() {
	err := http.ListenAndServe(h.addr, nil)
	if err != nil {
		panic(err)
	}
}

// InitPrometheus init a prometheus handler
func (h *Trace) InitPrometheus(vecs ...PrometheusVec) {
	h.prom.init(vecs...)
}

// GetCounter returns a prometheus count vector
func (h *Trace) GetCounter(name string) *prometheus.CounterVec {
	return h.prom.getCounter(name)
}

// GetSummary returns a prometheus summary vector
func (h *Trace) GetSummary(name string) *prometheus.SummaryVec {
	return h.prom.getSummary(name)
}

// GetHistogram returns a prometheus histogram vector
func (h *Trace) GetHistogram(name string) *prometheus.HistogramVec {
	return h.prom.getHistogram(name)
}

// GetGauge returns a prometheus gauge vector
func (h *Trace) GetGauge(name string) *prometheus.GaugeVec {
	return h.prom.getGauge(name)
}
