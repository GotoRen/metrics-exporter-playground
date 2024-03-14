package main

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type HandleApplication struct {
	ApplicationName string
}

// CustomCollector implements the prometheus.Collector interface
type CustomCollector struct {
	prometheus.Collector
}

// Describe implements the Describe method of the prometheus.Collector interface
func (c *CustomCollector) Describe(ch chan<- *prometheus.Desc) {
	// Describe method implementation
}

// Collect implements the Collect method of the prometheus.Collector interface
func (c *CustomCollector) Collect(ch chan<- prometheus.Metric) {
	// Collect method implementation
}

// RegisterMetrics registers all the metrics collectors
func RegisterMetrics() []prometheus.Collector {
	return []prometheus.Collector{
		cpuUsageMetric,
		memoryUsageMetric,
	}
}

// PushMetrics pushes metrics to Prometheus Pushgateway
func PushMetrics(endpoint string, jobName string, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName)
	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.Push()
}

var (
	cpuUsageMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "CPU usage of the application",
		},
		[]string{"application_name"},
	)

	memoryUsageMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage",
			Help: "Memory usage of the application",
		},
		[]string{"application_name"},
	)
)

func Register(r prometheus.Registerer) {
	r.MustRegister(
		cpuUsageMetric,
		memoryUsageMetric,
	)
}

// UpdateMetrics 関数から jobName を削除し、main 関数でジョブ名を指定する
func UpdateMetrics(h *HandleApplication) {
	// Simulated CPU and memory usage
	cpuUsage := 0.75   // 75% CPU usage
	memoryUsage := 512 // 512 MB memory usage

	// Update additional metrics
	cpuUsageMetric.WithLabelValues(h.ApplicationName).Set(cpuUsage)
	memoryUsageMetric.WithLabelValues(h.ApplicationName).Set(float64(memoryUsage))
}

func main() {
	// Register metrics
	Register(prometheus.DefaultRegisterer)

	applicationName := "hogehoge"
	jobName := "hoge"

	// Start routine to update metrics and push to Pushgateway every 5 seconds
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Simulate data to update metrics
				handleApp := &HandleApplication{
					ApplicationName: applicationName,
				}
				UpdateMetrics(handleApp)

				fmt.Println("確認")

				// Push metrics to Prometheus Pushgateway with the specified job name
				if err := PushMetrics("http://localhost:9091", jobName, RegisterMetrics()...); err != nil {
					fmt.Println("Failed to push metrics:", err)
				} else {
					log.Println("Metrics pushed successfully")
				}
			}
		}
	}()

	// Set timer to exit after 1 minute
	exitTimer := time.NewTimer(1 * time.Minute)
	<-exitTimer.C
	fmt.Println("Server shutdown...")
}
