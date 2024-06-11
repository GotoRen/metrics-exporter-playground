package pushmetric

import "github.com/prometheus/client_golang/prometheus"

// The Collector manages custom metrics to be sent to PushGateway.
type Collector struct {
	// asyncCollectors store custom metrics pushed continuously.
	// Custom metrics registered with this collector are sent to PushGateway asynchronously, separate from the main process.
	asynCollectors []prometheus.Collector

	// syncCollectors store custom metrics pushed at arbitrary times.
	// Custom metrics registered with this collector are sent to PushGateway at appropriate times.
	syncCollectors []prometheus.Collector
}

// NewCollector initializes a new Collector instance.
func NewCollector() *Collector {
	return &Collector{
		asynCollectors: make([]prometheus.Collector, 0),
		syncCollectors: make([]prometheus.Collector, 0),
	}
}

// RegisterAsyncMetrics registers custom metrics in asynCollectors.
// Registered custom metrics are pushed asynchronously and continuously based on the push interval value.
func (c *Collector) RegisterAsyncMetrics(cms ...prometheus.Collector) *Collector {
	c.asynCollectors = append(c.asynCollectors, cms...)
	return c
}

// RegisterSyncMetrics registers custom metrics in syncCollectors.
// Metrics are pushed by calling the export method at arbitrary times.
func (c *Collector) RegisterSyncMetrics(cms ...prometheus.Collector) *Collector {
	c.syncCollectors = append(c.syncCollectors, cms...)
	return c
}
