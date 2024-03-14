package pushmetric

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

// Define custom metric labels.
type CustomLabels struct {
	ApplicationName string
	InstanceName    string
}

// Define metric values.
type Metrics struct {
	CpuUsage     float64
	MemoryUsage  int64
	RequestCount int64
}

// Define gauge metrics.
var (
	cpuUsageMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "CPU usage of the application",
		},
		[]string{"application_name", "instance"},
	)

	memoryUsageMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage of the application",
		},
		[]string{"application_name", "instance"},
	)
)

// Define counter metric.
var (
	requestCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests processed",
		},
		[]string{"application_name", "instance"},
	)
)

// init initializes metrics.
func init() {
	Register(prometheus.DefaultRegisterer)
}

// Register registers all the metrics collectors.
func Register(r prometheus.Registerer) {
	r.MustRegister(
		cpuUsageMetric,
		memoryUsageMetric,
		requestCountMetric,
	)
}

// RegisterMetrics returns a slice of all the metrics collectors.
func RegisterMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		cpuUsageMetric,
		memoryUsageMetric,
		requestCountMetric,
	}
}

// UpdateGaugeMetric updates gauge metrics with given values.
func UpdateGaugeMetric(l *CustomLabels, cpuUsage float64, memoryUsage int64) {
	cpuUsageMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Set(cpuUsage)
	memoryUsageMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Set(float64(memoryUsage))
}

// IncrementCounterMetric increments counter metric.
func IncrementCounterMetric(l *CustomLabels) {
	requestCountMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Inc()
}

// getInstanceName returns the hostname of the current instance.
func getInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	return hostname
}

// getMetrics returns the collected metrics.
func getMetrics() *Metrics {
	m := new(Metrics)

	// TODO
	m.CpuUsage = 0.85
	m.MemoryUsage = 512

	return m
}
