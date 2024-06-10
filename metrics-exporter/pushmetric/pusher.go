package pushmetric

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// Exporter はPushGateway にメトリクスを送信する際の情報を登録します。
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
			// CPU
			if err := updateCPUMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf(": %w", err)
			}

			// Memory
			if err := updateMemoryMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf(": %w", err)
			}

			// Push Count
			if err := updatePushCountMetric(applicationNameLabelValue, instanceNameLabelValue); err != nil {
				return fmt.Errorf(": %w", err)
			}

			fmt.Println("[DEBUG] 非同期メトリクスを更新しました")

			if err := e.AsyncExport(ctx); err != nil {
				return fmt.Errorf("hoge: %w", err)
			}
		}
	}

	return nil
}

// 非同期的でメトリクスを出力する場合
func (e *Exporter) AsyncExport(ctx context.Context) error {
	if err := export(ctx, e.endPoint, e.jobName, e.client, e.collector.asynCollectors); err != nil {
		return fmt.Errorf("hoge: %w", err)
	} else {
		log.Println("[DEBUG] 非同期メトリクスを出力しました。")
	}
	return nil
}

// 任意のタイミングでメトリクスを出力する場合
func (e *Exporter) SyncExport(ctx context.Context) error {
	if err := export(ctx, e.endPoint, e.jobName, e.client, e.collector.syncCollectors); err != nil {
		return fmt.Errorf("hoge: %w", err)
	} else {
		log.Println("[DEBUG] 同期メトリクスを出力しました。")
	}
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
func (e *Exporter) Shutdown(gracePeriod time.Duration) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracePeriod)
	defer shutdownCancel()

	fmt.Println("Graceful Shutdown...")

	if err := e.AsyncExport(shutdownCtx); err != nil {
		return fmt.Errorf("hoge: %w", err)
	}

	if err := e.SyncExport(shutdownCtx); err != nil {
		return fmt.Errorf("hoge: %w", err)
	}

	return nil
}
