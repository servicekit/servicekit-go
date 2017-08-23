package trace

import (
	"fmt"
	"net/http"

	"github.com/servicekit/servicekit-go/coordinator"
	"github.com/servicekit/servicekit-go/logger"
	"github.com/servicekit/servicekit-go/spec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Trace struct {
	addr string

	prom *prom

	log *logger.Logger
}

func NewTrace(c coordinator.Coordinator, id string, name string, tags []string, host string, port int, ttl time.Duration, log *logger.Logger) *Trace {
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

	err = c.Register(ctx, s, ttl)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (h *Trace) Serve() {
	err := http.ListenAndServe(h.addr, nil)
	if err != nil {
		panic(err)
	}
}

func (h *Trace) InitPrometheus(vecs ...PrometheusVec) {
	h.prom.init(vecs...)
}

func (h *Trace) GetCounter(name string) *prometheus.CounterVec {
	return h.prom.getCounter(name)
}

func (h *Trace) GetSummary(name string) *prometheus.SummaryVec {
	return h.prom.getSummary(name)
}

func (h *Trace) GetHistogram(name string) *prometheus.HistogramVec {
	return h.prom.getHistogram(name)
}

func (h *Trace) GetGauge(name string) *prometheus.GaugeVec {
	return h.prom.getGauge(name)
}
