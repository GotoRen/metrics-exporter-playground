package pushmetric

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const pushGatewayEndPoint = "http://localhost:9091"

type Exporter struct {
	jobName         string
	applicationName string
	interval        time.Duration
	client          *http.Client // optional
}

func ExportConfig(jobName, applicationName string, interval time.Duration) *Exporter {
	return &Exporter{
		jobName:         jobName,
		applicationName: applicationName,
		interval:        interval,
		client:          &http.Client{},
	}
}

func (e *Exporter) WithClient(client *http.Client) *Exporter {
	e.client = client
	return e
}

func (e *Exporter) RuntineSequentialExporter(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			labels := &CustomLabels{
				ApplicationName: e.applicationName,
				InstanceName:    getInstanceName(),
			}

			m := getMetrics() // getMetrics

			UpdateGaugeMetric(labels, m.CpuUsage, m.MemoryUsage)
			IncrementCounterMetric(labels)

			export(ctx, e.jobName, e.client)
		}
	}
}

func export(ctx context.Context, jobName string, client *http.Client) {
	if err := pushMetrics(ctx, pushGatewayEndPoint, jobName, client, RegisterMetrics()...); err != nil {
		log.Println("Failed to push metrics:", err)
	} else {
		log.Println("Metrics pushed successfully")
	}
}

// pushMetrics pushes the provided collectors to the specified endpoint with the given job name.
func pushMetrics(ctx context.Context, endpoint string, jobName string, client *http.Client, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName)
	if client != nil {
		pusher = pusher.Client(client)
	}

	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.PushContext(ctx)
}

func pushMetricsWithCustomClient(ctx context.Context, endpoint string, jobName string, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName)

	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.PushContext(ctx)
}
