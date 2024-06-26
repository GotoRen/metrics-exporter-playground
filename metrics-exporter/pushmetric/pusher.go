package pushmetric

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// The Exporter configures information for sending metrics to PushGateway.
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

// RoutineSequentialExporter continuously pushes metrics to the Pushgateway.
func (e *Exporter) RoutineSequentialExporter(ctx context.Context) error {
	ticker := time.NewTicker(e.pushInterval)
	defer ticker.Stop()

	// Set label values
	applicationNameLabelValue := e.applicationName
	instanceNameLabelValue := getInstanceName()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// CPU utilization
			if err := updateCPUMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf("failed to update CPU metric: %w", err)
			}

			// Memory utilization
			if err := updateMemoryMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf("failed to update Memory metric: %w", err)
			}

			// Push count
			if err := updatePushCountMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf("failed to update Push Count metric: %w", err)
			}

			if err := e.AsyncExport(ctx); err != nil {
				return fmt.Errorf("failed to export metrics asynchronously: %w", err)
			}
		}
	}

	return nil
}

// AsyncExport passes asyncCollectors to export function to push custom metrics.
func (e *Exporter) AsyncExport(ctx context.Context) error {
	if err := export(ctx, e.endPoint, e.jobName, e.client, e.collector.asynCollectors); err != nil {
		return fmt.Errorf("failed to export async collectors: %w", err)
	}
	return nil
}

// SyncExport passes syncCollectors to export function to push custom metrics.
func (e *Exporter) SyncExport(ctx context.Context) error {
	if err := export(ctx, e.endPoint, e.jobName, e.client, e.collector.syncCollectors); err != nil {
		return fmt.Errorf("failed to export sync collectors: %w", err)
	}
	return nil
}

// export exports the provided collector (asynCollectors or syncCollectors) to the specified endpoint.
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
func (e *Exporter) Shutdown(ctx context.Context, gracePeriod time.Duration) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracePeriod)
	defer shutdownCancel()

	if err := e.AsyncExport(shutdownCtx); err != nil {
		return fmt.Errorf("failed to export async collectors during shutdown: %w", err)
	}

	if err := e.SyncExport(shutdownCtx); err != nil {
		return fmt.Errorf("failed to export sync collectors during shutdown: %w", err)
	}

	return nil
}
