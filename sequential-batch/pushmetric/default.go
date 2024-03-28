package pushmetric

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// デフォルトの GaugeMetric を定義：任意で呼び出す
var (
	cpuUtilizationMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_utilization",
			Help: "CPU usage of the application",
		},
		[]string{"application_name", "instance"},
	)
)

// デフォルトの CounterMetric を定義：任意で呼び出す
var (
	requestCountMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Number of requests processed",
		},
		[]string{"application_name", "instance"},
	)
)

func (collector *Collector) UpdateDefaultnMetric(lvs ...string) {
	value := GetCpuUtilization()
	SetGaugeMetric(cpuUtilizationMetric, value, lvs...)

	IncrementCounterMetric(requestCountMetric, lvs...)
}

// GetInstanceName returns the hostname of the current instance.
func GetInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	return hostname
}

func GetCpuUtilization() float64 {
	cpuUtilization, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("Error getting CPU usage:", err)
	}

	return cpuUtilization[0]
}

func GetMemoryUtilization() float64 {
	memoryUtilization, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory usage:", err)
	}

	return memoryUtilization.UsedPercent
}
