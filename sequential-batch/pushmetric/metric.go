package pushmetric

import "github.com/prometheus/client_golang/prometheus"

// Define custom metric labels.
type CustomLabels struct {
	ApplicationName string
	InstanceName    string
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

// Initialize metrics.
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
func UpdateGaugeMetric(l *CustomLabels, cpuUsage float64, memoryUsage int) {
	cpuUsageMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Set(cpuUsage)
	memoryUsageMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Set(float64(memoryUsage))
}

// IncrementCounterMetric increments counter metric.
func IncrementCounterMetric(l *CustomLabels) {
	requestCountMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Inc()
}
