package pushmetric

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Define default GaugeMetric
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

// Define default CounterMetric
var (
	pushCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "push_count",
			Help: "Number of times sent to PushGateway",
		},
		[]string{"application_name", "instance"},
	)
)

func (e *Exporter) updateDefaultnMetric(lvs ...string) {
	currentCpuUtilization := getCpuUtilization()
	currentMemoryUtilization := getMemoryUtilization()

	cpuUtilizationMetric.WithLabelValues(lvs...).Set(currentCpuUtilization)
	memoryUtilizationMetric.WithLabelValues(lvs...).Set(currentMemoryUtilization)
	pushCountMetric.WithLabelValues(lvs...).Inc()
}

// GetInstanceName returns the hostname of the current instance.
func GetInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	return hostname
}

// getCpuUtilization retrieves the current CPU utilization percentage.
func getCpuUtilization() float64 {
	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Error getting CPU usage:", err)
	}

	return cpuUtilization[0]
}

// getMemoryUtilization retrieves the current memory utilization percentage.
func getMemoryUtilization() float64 {
	memoryUtilization, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory usage:", err)
	}

	return memoryUtilization.UsedPercent
}
