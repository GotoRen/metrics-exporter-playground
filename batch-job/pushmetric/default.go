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
)

// Define default CounterMetric
var (
	requestCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests processed",
		},
		[]string{"application_name", "instance"},
	)
)

func (e *Exporter) updateDefaultnMetric(lvs ...string) {
	currentCpuUtilization := GetCpuUtilization()

	cpuUtilizationMetric.WithLabelValues(lvs...).Set(currentCpuUtilization)
	requestCountMetric.WithLabelValues(lvs...).Inc()

}

// GetInstanceName returns the hostname of the current instance.
func GetInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	return hostname
}

// GetCpuUtilization retrieves the current CPU utilization percentage.
func GetCpuUtilization() float64 {
	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Error getting CPU usage:", err)
	}

	return cpuUtilization[0]
}

// GetMemoryUtilization retrieves the current memory utilization percentage.
func GetMemoryUtilization() float64 {
	memoryUtilization, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory usage:", err)
	}

	return memoryUtilization.UsedPercent
}
