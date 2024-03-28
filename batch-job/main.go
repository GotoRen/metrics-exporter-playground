package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GotoRen/metrics-exporter-playground/batch-job/pushmetric"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	applicationName     = "sample_apps"
	jobName             = "sample_job"
	pushGatewayEndPoint = "http://localhost:9091"
)

const (
	pushInterval = 3 * time.Second // 3 秒毎に PushGateway にメトリクスを送信する例
	lifeTime     = 1 * time.Minute // 1 分後に Job を終了
)

func main() {
	log.Println("CronJob starting...")

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Collector を定義
	collector := pushmetric.NewCollector()

	// カスタムメトリクスを追加
	memoryUtilizationMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_utilization",
			Help: "Memory usage of the application",
		},
		[]string{"application_name", "instance"},
	)

	collector.WithDefaultMetrics().WithCustomMetrics(memoryUtilizationMetric)
	// // もしデフォルトのメトリクスを使用しない場合
	// collector.WithCustomMetrics(memoryUtilizationMetric)

	// エクスポートメトリクス情報を登録する
	config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector)

	// // カスタムクライアントを使用する場合
	// client := WithCustomClient()
	// config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint).WithClient(client)

	go config.RuntineSequentialExporter(ctx) // PushGateway にシーケンシャルにメトリクスをエクスポートする

	// メトリクスを任意のタイミングで更新する
	go func() {
		ticker := time.NewTicker(1 * time.Second) // 1 秒毎にメトリクスが更新される場合
		defer ticker.Stop()

		// Set label values
		applicationNameLabelValue := applicationName
		instanceNameLavelValue := pushmetric.GetInstanceName()

		for {
			select {
			case <-ticker.C:
				currentMemoryUtilization := pushmetric.GetMemoryUtilization()
				memoryUtilizationMetric.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(currentMemoryUtilization)
			}
		}
	}()

	// make the main routine wait.
	time.Sleep(lifeTime)
	log.Println("CronJob completed")
}

// Example: Define a custom HTTP client.
func WithCustomClient() *http.Client {
	requestTimeoutLimit := 5 * time.Second // 5 秒間 レスポンスがない場合にタイムアウトエラー

	return &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
}
