package pushmetric

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Define CPU usage GaugeMetric.
var (
	CpuUtilizationMetric = prometheus.NewGaugeVec( //nolint:gochecknoglobals
		prometheus.GaugeOpts{
			Name: "cpu_utilization",
			Help: "CPU usage of the application",
		},
		[]string{"application_name", "instance"},
	)
)

// Define memory usage GaugeMetric.
var (
	MemoryUtilizationMetric = prometheus.NewGaugeVec( //nolint:gochecknoglobals
		prometheus.GaugeOpts{
			Name: "memory_utilization",
			Help: "Memory usage of the application",
		},
		[]string{"application_name", "instance"},
	)
)

// Define push count CounterMetric.
var (
	PushCountMetric = prometheus.NewCounterVec( //nolint:gochecknoglobals
		prometheus.CounterOpts{
			Name: "push_count",
			Help: "Number of times sent to PushGateway",
		},
		[]string{"application_name", "instance"},
	)
)

// updateCPUMetric updates CPU metrics.
func updateCPUMetric(lvs ...string) error {
	currentCPUUtilization, err := getCPUUtilization()
	if err != nil {
		return fmt.Errorf("error getting CPU usage: %w", err)
	}
	CpuUtilizationMetric.WithLabelValues(lvs...).Set(currentCPUUtilization)
	return nil
}

// updateMemoryMetric updates memory metrics.
func updateMemoryMetric(lvs ...string) error {
	currentMemoryUtilization, err := getMemoryUtilization()
	if err != nil {
		return fmt.Errorf("error getting memory usage: %w", err)
	}
	MemoryUtilizationMetric.WithLabelValues(lvs...).Set(currentMemoryUtilization)
	return nil
}

// updatePushCountMetric updates push count metrics.
func updatePushCountMetric(lvs ...string) error {
	PushCountMetric.WithLabelValues(lvs...).Inc()
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
