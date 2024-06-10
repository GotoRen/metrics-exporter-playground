package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GotoRen/metrics-exporter-playground/metrics-exporter/pushmetric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/net"
)

// 正常終了：CronJob が実行されたことにより終了（プログラムが終了する）=> ライブラリでは、これを ContextWithTimeOut として表現する
// - 終了に伴う処理
//   - Goroutine のキャンセル

// 異常終了：SIGTERM 受信による強制シャットダウン
// - 終了に伴う処理
//   - Goroutine のキャンセル
//   - メトリクスを出力して終了

const (
	applicationName     = "sample_apps"           // 必須：アプリケーション名
	jobName             = "sample_job"            // 必須：Prometheus がエンドポイントグループを識別するための Job ラベル
	pushGatewayEndPoint = "http://localhost:9091" // 必須：PushGateway エンドポイント
)

const (
	pushInterval = 3 * time.Second // 必須：PushGateway にメトリクスを送信する間隔
	gracePeriod  = 5 * time.Second // 必須：プロセス終了までの待機時間
)

const (
	lifeTime = 1 * time.Minute // 任意：プロセスの実行時間（CronJob の実行時間に相当）
)

// 任意：カスタムメトリクスを定義
var (
	bytesSentCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_bytes_sent",
			Help: "Current bytes sent over the network",
		},
		[]string{"application_name", "instance"},
	)
)

var (
	bytesRecvCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_bytes_received",
			Help: "Current bytes received over the network",
		},
		[]string{"application_name", "instance"},
	)
)

func main() {
	log.Println("Job starting...")

	// 一定時間経過後後にキャンセル（ctx.Done）を実行（CronJob の正常終了）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// シグナルを捕捉するチャネルを設定
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 必須：Collector を定義
	collector := pushmetric.NewCollector()

	// 任意：ライブラリ側で準備しているカスタムメトリクスを登録（AsyncMetrics として登録）
	collector.RegisterAsyncMetrics(
		pushmetric.CpuUtilizationMetric,    // 必要に応じてデフォルトで準備しているメトリクスを登録
		pushmetric.MemoryUtilizationMetric, // 必要に応じてデフォルトで準備しているメトリクスを登録
		pushmetric.PushCountMetric,         // 必要に応じてデフォルトで準備しているメトリクスを登録
		bytesSentCounter,
	)

	// 任意：ユーザ定義カスタムメトリクスを追加（SyncMetrics として登録）
	collector.RegisterSyncMetrics(
		bytesRecvCounter,
	)

	// 任意：カスタムクライアントを使用する場合
	client := WithCustomClient()

	// 必須：Exporter を定義
	config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector).WithClient(client)

	// 任意：定常的にメトリクスを出力するエクスポータを起動します
	errCh := make(chan error, 1) // RoutineSequentialExporter からのエラーを補足するためのチャネルを作成
	go func() {
		errCh <- config.RoutineSequentialExporter(ctx)
	}()

	fmt.Println("待機しています 1")

	// -------------------------------------------------------------------------------------------------------------- //
	// COMMENT: ライブラリの使用者はカスタムメトリクスの更新処理を書きます（任意のタイミング）
	// メトリクスを任意のタイミングで更新
	go func() {
		ticker := time.NewTicker(1 * time.Second) // 1 秒毎にメトリクスが更新される場合
		defer ticker.Stop()

		// カスタムメトリクスに対するラベルをセット
		applicationNameLabelValue := applicationName
		instanceNameLavelValue := pushmetric.GetInstanceName()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// ex. 毎秒 ネットワーク送受信データ量を取得してメトリクスを更新
				bytesSent, _ := monitorNetworkSpeed()
				bytesSentCounter.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(bytesSent)

				fmt.Println("[DEBUG] Sent:", bytesSent)
			}
		}
	}()

	fmt.Println("待機しています 2")

	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	case sig := <-signalChan: // シグナルが送信された場合
		log.Printf("Received os.Signal: %v. Initiating graceful shutdown...", sig)
	case <-time.After(lifeTime): // lifeTime 経過した場合（ex. CronJob の終了）
		log.Println("LifeTime elapsed.")
	}

	cancel() // メトリクス更新 の goroutine をキャンセル（ctx.Done を実行）する

	// -------------------------------------------------------------------------------------------------------------- //
	// COMMENT: ライブラリの使用者はカスタムメトリクスの更新処理を書きます（任意のタイミング）
	// 既存の Gouroutine 停止後に更新します
	applicationNameLabelValue := applicationName
	instanceNameLavelValue := pushmetric.GetInstanceName()
	_, bytesRecv := monitorNetworkSpeed()
	bytesRecvCounter.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(bytesRecv)

	fmt.Println("[DEBUG] Recv:", bytesRecv)
	// -------------------------------------------------------------------------------------------------------------- //

	fmt.Println("待機しています 4")

	// 必須：Shutdown を実行して更新したメトリクスを確実に出力します
	if err := config.Shutdown(gracePeriod); err != nil {
		log.Fatalf("failed to gracefully shutdown: %v\n", err)
	}

	log.Println("Job completed")
}

//--------------------------------------------------------------------------------------------------------------//
// sample function
//--------------------------------------------------------------------------------------------------------------//

// WithCustomClient define a custom HTTP client.
func WithCustomClient() *http.Client {
	requestTimeoutLimit := 5 * time.Second // 5 秒間 レスポンスがない場合にタイムアウトエラー

	return &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
}

var (
	prevNetIOCounters []net.IOCountersStat
	prevTime          time.Time
)

// monitorNetworkSpeed monitors the upload and download speed of a specified network interface
// at a specified interval. The speeds are measured in KB/s.
func monitorNetworkSpeed() (float64, float64) {
	var targetNIC string = "en0"

	// Get current network I/O counters
	netIOCounters, err := net.IOCounters(true)
	if err != nil {
		log.Println("Error getting network I/O info:", err)
		return 0, 0
	}

	currentTime := time.Now()

	// Skip processing if prevNetIOCounters has not been initialized
	if prevNetIOCounters == nil || len(prevNetIOCounters) == 0 {
		prevNetIOCounters = netIOCounters
		prevTime = currentTime
		return 0, 0 // Initial call, no prior data to compare
	}

	for i, counter := range netIOCounters {
		if counter.Name == targetNIC { // Check for specific network interface
			// Calculate the time difference between the current and previous measurements
			timeDiff := currentTime.Sub(prevTime).Seconds()

			// Ensure the index exists before dereferencing
			if i >= len(prevNetIOCounters) {
				log.Println("Index out of range for prevNetIOCounters")
				continue
			}

			// Calculate the difference in sent and received bytes
			sentBytesDiff := counter.BytesSent - prevNetIOCounters[i].BytesSent
			recvBytesDiff := counter.BytesRecv - prevNetIOCounters[i].BytesRecv

			// Calculate upload and download speed in Bytes/s
			uploadSpeed := float64(sentBytesDiff) / timeDiff
			downloadSpeed := float64(recvBytesDiff) / timeDiff

			// Convert Bytes/s to KB/s
			uploadSpeedKBps := uploadSpeed / 1024
			downloadSpeedKBps := downloadSpeed / 1024

			// fmt.Printf("Upload Speed: %.2f KB/s, Download Speed: %.2f KB/s\n", uploadSpeedKBps, downloadSpeedKBps)

			// Store current counters for the next comparison
			prevNetIOCounters = netIOCounters
			prevTime = currentTime

			return uploadSpeedKBps, downloadSpeedKBps
		}
	}

	// Store current counters for the next comparison
	prevNetIOCounters = netIOCounters
	prevTime = currentTime

	return 0, 0
}
