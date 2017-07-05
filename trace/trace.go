package trace

import (
    "net/http"

    "github.com/servicekit/servicekit-go/logger"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Trace struct {
    addr string

    prom *prom

    log *logger.Logger
}

func NewTrace(addr string, log *logger.Logger) *Trace {
    t := &Trace{
        addr: addr,
    }

    http.Handle("/metrics", promhttp.Handler())

    t.prom = &prom{
        path:       "/metrics",
        collectors: make(map[string]prometheus.Collector),
        log:        log,
    }

    return t
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
