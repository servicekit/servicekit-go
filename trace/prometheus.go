package trace

import (
	"sync"

	"github.com/servicekit/servicekit-go/logger"

	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusVec defin a prometheus vector that have two method:
// GetName returns a vector name
// GetCollector return a prometheus.Collector
type PrometheusVec interface {
	GetName() string
	GetCollector() prometheus.Collector
}

// PrometheusCounter is a PrometheusVec that collect counter vector
type PrometheusCounter struct {
	Name   string
	Help   string
	Labels []string
}

// GetName return a count vector name
func (p *PrometheusCounter) GetName() string {
	return p.Name
}

// GetCollector return a collector with a counter options
func (p *PrometheusCounter) GetCollector() prometheus.Collector {
	opts := prometheus.CounterOpts{
		Name: p.Name,
		Help: p.Help,
	}

	return prometheus.NewCounterVec(opts, p.Labels)
}

// PrometheusGauge is a PrometheusVec that collect gauge vector
type PrometheusGauge struct {
	Name   string
	Help   string
	Labels []string
}

// GetName returns a gauge vector name
func (p *PrometheusGauge) GetName() string {
	return p.Name
}

// GetCollector returns a collector with a gauge options
func (p *PrometheusGauge) GetCollector() prometheus.Collector {
	opts := prometheus.GaugeOpts{
		Name: p.Name,
		Help: p.Help,
	}

	return prometheus.NewGaugeVec(opts, p.Labels)
}

// PrometheusHistogram is a PrometheusVec that collect histogram vector
type PrometheusHistogram struct {
	Name   string
	Help   string
	Labels []string
}

// GetName returns a histogram vector name
func (p *PrometheusHistogram) GetName() string {
	return p.Name
}

// GetCollector returns a collector with a histogram options
func (p *PrometheusHistogram) GetCollector() prometheus.Collector {
	opts := prometheus.HistogramOpts{
		Name: p.Name,
		Help: p.Help,
	}

	return prometheus.NewHistogramVec(opts, p.Labels)
}

type prom struct {
	path       string
	collectors map[string]prometheus.Collector

	inited bool

	log *logger.Logger

	sync.Mutex
}

func (p *prom) init(vecs ...PrometheusVec) {
	p.Lock()

	if p.inited == true {
		return
	}

	for _, v := range vecs {
		p.collectors[v.GetName()] = v.GetCollector()
		prometheus.MustRegister(p.collectors[v.GetName()])
	}

	p.inited = true

	p.Unlock()
}

func (p *prom) getCounter(name string) *prometheus.CounterVec {
	v, ok := p.collectors[name]
	if ok == false {
		return nil
	}

	c, ok := v.(*prometheus.CounterVec)
	if ok == false {
		return nil
	}

	return c
}

func (p *prom) getSummary(name string) *prometheus.SummaryVec {
	v, ok := p.collectors[name]
	if ok == false {
		return nil
	}

	c, ok := v.(*prometheus.SummaryVec)
	if ok == false {
		return nil
	}

	return c
}

func (p *prom) getHistogram(name string) *prometheus.HistogramVec {
	v, ok := p.collectors[name]
	if ok == false {
		return nil
	}

	c, ok := v.(*prometheus.HistogramVec)
	if ok == false {
		return nil
	}

	return c
}

func (p *prom) getGauge(name string) *prometheus.GaugeVec {
	v, ok := p.collectors[name]
	if ok == false {
		return nil
	}

	c, ok := v.(*prometheus.GaugeVec)
	if ok == false {
		return nil
	}

	return c
}
