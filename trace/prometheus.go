package trace

import (
	"sync"

	"github.com/servicekit/servicekit-go/logger"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusVec interface {
	GetName() string
	GetCollector() prometheus.Collector
}

type PrometheusCounter struct {
	Name   string
	Help   string
	Labels []string
}

func (p *PrometheusCounter) GetName() string {
	return p.Name
}

func (p *PrometheusCounter) GetCollector() prometheus.Collector {
	opts := prometheus.CounterOpts{
		Name: p.Name,
		Help: p.Help,
	}

	return prometheus.NewCounterVec(opts, p.Labels)
}

type PrometheusGauge struct {
	Name   string
	Help   string
	Labels []string
}

func (p *PrometheusGauge) GetName() string {
	return p.Name
}

func (p *PrometheusGauge) GetCollector() prometheus.Collector {
	opts := prometheus.GaugeOpts{
		Name: p.Name,
		Help: p.Help,
	}

	return prometheus.NewGaugeVec(opts, p.Labels)
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
