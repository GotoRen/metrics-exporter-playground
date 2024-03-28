package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GotoRen/metrics-exporter-playground/sequential-batch/pushmetric"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	applicationName     = "sample_apps"
	jobName             = "sample_job"
	pushGatewayEndPoint = "http://localhost:9091"
)

const (
	pushInterval = 3 * time.Second // 3 秒毎にメトリクスを送信する例
	lifeTime     = 3 * time.Minute // 1 分後にJobを終了
)

func main() {
	log.Println("CronJob starting...")

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ここで Collector を定義
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

	// ここで エクスポートメトリクス 情報を登録
	config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector)

	// // カスタムクライアントを使用する場合
	// client := WithCustomClient()
	// config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint).WithClient(client)

	go config.RuntineSequentialExporter(ctx) // push ルーチンを回す

	// ここで任意のタイミング（とりあえず 5 秒毎に）でメトリクス情報を更新
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				value := pushmetric.GetMemoryUtilization()
				label1 := applicationName
				label2 := pushmetric.GetInstanceName()
				pushmetric.SetGaugeMetric(memoryUtilizationMetric, value, label1, label2)

				collector.UpdateDefaultnMetric(label1, label2)
			}
		}
	}()

	// wait main routine
	time.Sleep(lifeTime)
	log.Println("CronJob completed")
}

// 任意: カスタムクライアントを定義
func WithCustomClient() *http.Client {
	requestTimeoutLimit := 5 * time.Second // 5 秒間 レスポンスがない場合にタイムアウトエラー

	return &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
}
