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

func RoutineSequentialExport(ctx context.Context, client *http.Client, jobName, applicationName string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			labels := &CustomLabels{
				ApplicationName: applicationName,
				InstanceName:    getInstanceName(),
			}

			m := getMetrics() // getMetrics

			UpdateGaugeMetric(labels, m.CpuUsage, m.MemoryUsage)
			IncrementCounterMetric(labels)

			export(ctx, client, jobName)
		}
	}
}

func export(ctx context.Context, client *http.Client, jobName string) {
	if err := pushMetrics(ctx, pushGatewayEndPoint, jobName, client, RegisterMetrics()...); err != nil {
		log.Println("Failed to push metrics:", err)
	} else {
		log.Println("Metrics pushed successfully")
	}
}

// pushMetrics pushes the provided collectors to the specified endpoint with the given job name.
func pushMetrics(ctx context.Context, endpoint string, jobName string, client *http.Client, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName).Client(client)

	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.PushContext(ctx)
}
