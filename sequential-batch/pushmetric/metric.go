package pushmetric

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Define custom metric labels.
type CustomLabels struct {
	ApplicationName string
	InstanceName    string
}

// Define metric values.
type Metrics struct {
	CpuUtilization    float64
	MemoryUtilization float64
	RequestCount      int64
}

// Define gauge metrics.
var (
	cpuUtilizationMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_utilization",
			Help: "CPU usage of the application",
		},
		[]string{"application_name", "instance"},
	)

	memoryUtilizationMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_utilization",
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
		cpuUtilizationMetric,
		memoryUtilizationMetric,
		requestCountMetric,
	)
}

// RegisterMetrics returns a slice of all the metrics collectors.
func RegisterMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		cpuUtilizationMetric,
		memoryUtilizationMetric,
		requestCountMetric,
	}
}

// UpdateGaugeMetric updates gauge metrics with given values.
func UpdateGaugeMetric(l *CustomLabels, cpuUtilization, memoryUtilization float64) {
	cpuUtilizationMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Set(cpuUtilization)
	memoryUtilizationMetric.WithLabelValues(l.ApplicationName, l.InstanceName).Set(memoryUtilization)
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
	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Error getting CPU usage:", err)
	}

	memoryUtilization, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory usage:", err)
	}

	return &Metrics{
		CpuUtilization:    cpuUtilization[0],
		MemoryUtilization: memoryUtilization.UsedPercent,
	}
}
