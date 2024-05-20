package pushmetric

import (
	"context"
	"fmt"
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

// RoutineSequentialExporter continuously collects metrics and pushes them to the Pushgateway.
func (e *Exporter) RoutineSequentialExporter(ctx context.Context) error {
	ticker := time.NewTicker(e.pushInterval)
	defer ticker.Stop()

	// Set label values
	applicationNameLabelValue := e.applicationName
	instanceNameLabelValue := GetInstanceName()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := updateDefaultMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf("error updating metrics: %w", err)
			}
			if err := e.Export(ctx); err != nil {
				return fmt.Errorf("error exporting metrics: %w", err)
			}
		}
	}

	return nil
}

// Export exports the metrics to the Pushgateway.
func (e *Exporter) Export(ctx context.Context) error {
	if err := export(ctx, e.endPoint, e.jobName, e.client, e.collector.collectors); err != nil {
		return fmt.Errorf("failed to push metrics: %w", err)
	}

	// fmt.Println("[DEBUG] call: Push metrics")

	return nil
}

// export exports the provided collectors to the specified endpoint with the given job name.
func export(ctx context.Context, ep string, name string, client *http.Client, collectors []prometheus.Collector) error {
	pusher := push.New(ep, name)
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

// Shutdown exports the final metrics and shuts down the exporter.
func (e *Exporter) Shutdown(ctx context.Context) error {
	if err := e.Export(ctx); err != nil {
		return fmt.Errorf("error occurred while exporting final metrics: %w", err)
	}

	return nil
}
