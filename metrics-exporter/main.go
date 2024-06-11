package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/pushmetric"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	applicationName     = "sample_apps"           // 必須：アプリケーション名
	jobName             = "sample_job"            // 必須：Prometheus がエンドポイントグループを識別するための Job ラベル
	pushGatewayEndPoint = "http://localhost:9091" // 必須：PushGateway エンドポイント
	// pushGatewayEndPoint = "http://prometheus-pushgateway.monitoring.svc.cluster.local:9091" // 必須：PushGateway エンドポイント
	pushInterval = 1 * time.Second // 必須：PushGateway にメトリクスを送信する間隔
)

const (
	gracePeriod = 5 * time.Second // 必須：プロセス終了までの待機時間
)

// WithCustomClient define a custom HTTP client.
func WithCustomClient() *http.Client {
	requestTimeoutLimit := 3 * time.Minute // 5 秒間 レスポンスがない場合にタイムアウトエラー

	return &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
}

// カスタムメトリクスを追加
var (
	violationGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "abema",
		Name:      "gatekeeper_detailed_violations",
		Help:      "Number of detailed violations",
	}, []string{"name", "namespace", "group", "version", "kind", "violation_kind", "severity", "violation"})

	constraintGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "abema",
		Name:      "gatekeeper_detailed_constraints",
		Help:      "Detailed constraints info",
	}, []string{"name", "kind", "severity", "recommendation"})
)

func main() {
	log.Println("Job starting...")

	const deavyWorkDuration = 30 * time.Second // よしなに設定
	var count = 1

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// シグナルを捕捉するチャネルを設定
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	collector := pushmetric.NewCollector()

	collector.RegisterAsyncMetrics(
		pushmetric.CpuUtilizationMetric,    // CPU 使用率
		pushmetric.MemoryUtilizationMetric, // メモリ使用率
		pushmetric.PushCountMetric,         // Push 回数
	)

	// カスタムメトリクスを登録
	collector.RegisterSyncMetrics(
		violationGauge,
		constraintGauge,
	)

	client := WithCustomClient()
	config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector).WithClient(client)

	errCh := make(chan error, 1) // RoutineSequentialExporter からのエラーを補足するためのチャネルを作成
	go func() {
		errCh <- config.RoutineSequentialExporter(ctx)
		if err := <-errCh; err != nil {
			panic(err)
		}
	}()

	// シグナル受信すると Shutdown メソッドを呼び出してメトリクスをエクスポートする
	go func() {
		select {
		case <-signalChan:
			log.Println("Received shutdown signal")
			if err := config.Shutdown(ctx, gracePeriod); err != nil {
				log.Fatalf("failed to gracefully shutdown: %v\n", err)
			}
			os.Exit(0)
		case <-ctx.Done():
		}
	}()

	/*** CronJob が deavyWorkDuration 程度の重い処理をしていると仮定 ***/
	//---------------------------------------------------------------------------------//
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := time.After(deavyWorkDuration)

loop:
	for {
		select {
		case <-ticker.C:
			fmt.Printf("%d: Hello, Loop!\n", count) // deavyWorkDuration 間 1 秒毎に Hello, Loop! を出力
			count++

			// カスタムメトリクスを更新
			violationGauge.WithLabelValues(
				"test",           // name
				"monitoring",     // namespace
				"apps",           // group
				"v1",             // version
				"ReplicaSet",     // kind
				"MetadataFormat", // violation_kind
				"3",              // severity
				"warn-unconfigured-abema-required-labels", // violation
			).Add(1)

			constraintGauge.WithLabelValues(
				"test",            // name
				"apps",            // kind
				"8",               // severity
				"hoge を修正してください。", // recommendation
			).Add(1)

		case <-done:
			log.Println("Completed a heavy task!")
			break loop
		}
	}
	//---------------------------------------------------------------------------------//

	// カスタムメトリクスを PushGateway に送信
	if err := config.SyncExport(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Job completed!")
}
