package pushmetric

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type Exporter struct {
	jobName         string
	applicationName string
	pushInterval    time.Duration
	collector       *Collector
	client          *http.Client
	endPoint        string
}

// New creates a new Exporter instance with the provided parameters.
func New(jobName, applicationName string, pushInterval time.Duration, endPoint string, collector *Collector) *Exporter {
	return &Exporter{
		jobName:         jobName,
		applicationName: applicationName,
		pushInterval:    pushInterval,
		client:          &http.Client{},
		endPoint:        endPoint,
		collector:       collector,
	}
}

// RuntineSequentialExporter continuously collects metrics and pushes them to the Pushgateway.
func (e *Exporter) RuntineSequentialExporter(ctx context.Context) {
	ticker := time.NewTicker(e.pushInterval)
	defer ticker.Stop()

	// Set label values
	applicationNameLabelValue := e.applicationName
	instanceNameLabelValue := GetInstanceName()

	for {
		select {
		case <-ticker.C:
			e.updateDefaultnMetric(applicationNameLabelValue, instanceNameLabelValue)
			e.Export(ctx)
		}
	}
}

// Export exports the metrics to the Pushgateway.
func (e *Exporter) Export(ctx context.Context) {
	if err := export(ctx, e.endPoint, e.jobName, e.client, e.collector.collectors); err != nil {
		log.Println("Failed to push metrics:", err)
	} else {
		log.Println("Metrics pushed successfully")
	}
}

// export exports the provided collectors to the specified endpoint with the given job name.
func export(ctx context.Context, endpoint string, jobName string, client *http.Client, collectors []prometheus.Collector) error {
	pusher := push.New(endpoint, jobName)
	if client != nil {
		pusher = pusher.Client(client)
	}

	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.PushContext(ctx)
}

// WithClient sets the HTTP client for the Exporter and returns the modified Exporter instance.
func (e *Exporter) WithClient(client *http.Client) *Exporter {
	e.client = client
	return e
}
