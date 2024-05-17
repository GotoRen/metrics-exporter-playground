package pushmetric

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Define default GaugeMetric.
var (
	cpuUtilizationMetric = prometheus.NewGaugeVec( //nolint:gochecknoglobals
		prometheus.GaugeOpts{
			Name: "cpu_utilization",
			Help: "CPU usage of the application",
		},
		[]string{"application_name", "instance"},
	)

	memoryUtilizationMetric = prometheus.NewGaugeVec( //nolint:gochecknoglobals
		prometheus.GaugeOpts{
			Name: "memory_utilization",
			Help: "Memory usage of the application",
		},
		[]string{"application_name", "instance"},
	)
)

// Define default CounterMetric.
var (
	pushCountMetric = prometheus.NewCounterVec( //nolint:gochecknoglobals
		prometheus.CounterOpts{
			Name: "push_count",
			Help: "Number of times sent to PushGateway",
		},
		[]string{"application_name", "instance"},
	)
)

// updateDefaultMetric updates default metrics.
func updateDefaultMetric(lvs ...string) error {
	currentCPUUtilization, err := getCPUUtilization()
	if err != nil {
		return fmt.Errorf("error getting CPU usage: %w", err)
	}

	currentMemoryUtilization, err := getMemoryUtilization()
	if err != nil {
		return fmt.Errorf("error getting memory usage: %w", err)
	}

	cpuUtilizationMetric.WithLabelValues(lvs...).Set(currentCPUUtilization)
	memoryUtilizationMetric.WithLabelValues(lvs...).Set(currentMemoryUtilization)
	pushCountMetric.WithLabelValues(lvs...).Inc()

	return nil
}

// GetInstanceName returns the hostname of the current instance.
func GetInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}

	return hostname
}

// getCPUUtilization retrieves the current CPU utilization percentage.
func getCPUUtilization() (float64, error) {
	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get CPU utilization: %w", err)
	}

	return cpuUtilization[0], nil
}

// getMemoryUtilization retrieves the current memory utilization percentage.
func getMemoryUtilization() (float64, error) {
	memoryUtilization, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to get memory utilization: %w", err)
	}

	return memoryUtilization.UsedPercent, nil
}
