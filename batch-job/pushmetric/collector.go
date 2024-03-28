package pushmetric

import "github.com/prometheus/client_golang/prometheus"

type Collector struct {
	collectors []prometheus.Collector
}

// NewCollector initializes a new Collector instance.
func NewCollector() *Collector {
	return &Collector{
		collectors: make([]prometheus.Collector, 0),
	}
}

// WithDefaultMetrics adds default metrics to collectors.
func (c *Collector) WithDefaultMetrics() *Collector {
	c.collectors = append(c.collectors, cpuUtilizationMetric, requestCountMetric)
	return c
}

// WithCustomMetrics adds custom metrics to collectors.
func (c *Collector) WithCustomMetrics(cms ...prometheus.Collector) {
	for _, cm := range cms {
		c.collectors = append(c.collectors, cm)
	}
}
