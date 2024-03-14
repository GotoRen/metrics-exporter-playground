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
	client          *http.Client
}

func New(jobName, applicationName string, interval time.Duration) *Exporter {
	return &Exporter{
		jobName:         jobName,
		applicationName: applicationName,
		interval:        interval,
		client:          &http.Client{},
	}
}

/* optional method */
// WithClient sets the HTTP client for the Exporter and returns the modified Exporter instance.
func (e *Exporter) WithClient(client *http.Client) *Exporter {
	e.client = client
	return e
}

// RuntineSequentialExporter continuously collects metrics and pushes them to the Pushgateway.
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

			m := getMetrics()

			UpdateGaugeMetric(labels, m.CpuUtilization, m.MemoryUtilization)
			IncrementCounterMetric(labels)

			export(ctx, e.jobName, e.client)
		}
	}
}

// export pushes the metrics to the Pushgateway.
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
