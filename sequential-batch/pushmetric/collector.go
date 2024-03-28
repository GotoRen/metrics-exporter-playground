package pushmetric

import "github.com/prometheus/client_golang/prometheus"

type Collector struct {
	collectors []prometheus.Collector
}

// NewCollector は Collector を初期化する
func NewCollector() *Collector {
	return &Collector{
		collectors: make([]prometheus.Collector, 0),
	}
}

// WithDefaultMetricsは、デフォルトのメトリクスを collectors に追加する
func (c *Collector) WithDefaultMetrics() *Collector {
	c.collectors = append(c.collectors, cpuUtilizationMetric, requestCountMetric)
	return c
}

// WithCustomMetricsは、カスタムメトリクスを collectors に追加する
func (c *Collector) WithCustomMetrics(cms ...prometheus.Collector) {
	for _, cm := range cms {
		c.collectors = append(c.collectors, cm)
	}
}

// Register registers all the metrics collectors.
func Register(c *Collector) {
	r := prometheus.DefaultRegisterer
	r.MustRegister(
		c.collectors...,
	)
}
