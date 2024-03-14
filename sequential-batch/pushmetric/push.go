package pushmetric

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const pushGatewayEndPoint = "http://localhost:9091"

type HandleApplication struct {
	ApplicationName string
	InstanceName    string
}

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

// RegisterMetrics registers all the metrics collectors
func RegisterMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		cpuUsageMetric,
		memoryUsageMetric,
	}
}

func InitRegister() {
	Register(prometheus.DefaultRegisterer)
}

func Register(r prometheus.Registerer) {
	r.MustRegister(
		cpuUsageMetric,
		memoryUsageMetric,
	)
}

// // CustomCollector implements the prometheus.Collector interface
// type CustomCollector struct {
// 	prometheus.Collector
// }
// // Describe implements the Describe method of the prometheus.Collector interface
// func (c *CustomCollector) Describe(ch chan<- *prometheus.Desc) {
// 	// Describe method implementation
// }
// // Collect implements the Collect method of the prometheus.Collector interface
// func (c *CustomCollector) Collect(ch chan<- prometheus.Metric) {
// 	// Collect method implementation
// }

func UpdateMetrics(h *HandleApplication, cpuUsage float64, memoryUsage int) {
	// Update additional metrics
	cpuUsageMetric.WithLabelValues(h.ApplicationName, h.InstanceName).Set(cpuUsage)
	memoryUsageMetric.WithLabelValues(h.ApplicationName, h.InstanceName).Set(float64(memoryUsage))
}

// PushMetrics pushes metrics to Prometheus Pushgateway
func PushMetrics(endpoint string, jobName string, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName)
	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.Push()
}

func Exports(handleApp *HandleApplication, jobName string) {

	if err := PushMetrics(pushGatewayEndPoint, jobName, RegisterMetrics()...); err != nil {
		fmt.Println("Failed to push metrics:", err)
	} else {
		log.Println("Metrics pushed successfully")
	}
}
