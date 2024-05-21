package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/GotoRen/metrics-exporter-playground/metrics-exporter/pushmetric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/net"
)

const (
	applicationName     = "sample_apps"
	jobName             = "sample_job"
	pushGatewayEndPoint = "http://localhost:9091"
)

const (
	pushInterval = 3 * time.Second  // PushGateway にメトリクスを送信する間隔
	lifeTime     = 30 * time.Second // プロセスの実行時間（CronJob の実行時間に相当）
	gracePeriod  = 5 * time.Second  // プロセス終了までの待機時間
)

// カスタムメトリクスを定義
var (
	bytesSentCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_bytes_sent",
			Help: "Current bytes sent over the network",
		},
		[]string{"application_name", "instance"},
	)

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

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), lifeTime)
	defer cancel()

	// シグナルを捕捉するチャネルを設定
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Collector を定義
	collector := pushmetric.NewCollector()

	// デフォルトメトリクスを追加
	collector.WithDefaultMetrics()

	// カスタムメトリクスを追加
	collector.WithCustomMetrics(bytesSentCounter, bytesRecvCounter) // 任意のメトリクスを追加できる

	// カスタムクライアントを使用する場合
	client := WithCustomClient()
	config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector).WithClient(client)

	errCh := make(chan error, 1)
	go func() {
		errCh <- config.RoutineSequentialExporter(ctx) // PushGateway にシーケンシャルにメトリクスを出力
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1) // Add a WaitGroup counter

	// メトリクスを任意のタイミングで更新する
	go func() {
		defer wg.Done() // 終了時にこのゴルーチンが完了したことを WaitGroup に通知する

		ticker := time.NewTicker(1 * time.Second) // 1 秒毎にメトリクスが更新される場合
		defer ticker.Stop()

		// Set label values
		applicationNameLabelValue := applicationName
		instanceNameLavelValue := pushmetric.GetInstanceName()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Println("[DEBUG] call: Update default metrics")

				bytesSent, bytesRecv := monitorNetworkSpeed()
				bytesSentCounter.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(bytesSent)
				bytesRecvCounter.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(bytesRecv)
			}
		}
	}()

	select {
	case err := <-errCh: // エラーが発生した場合
		if err != nil {
			log.Fatalf("Exporter error: %v\n", err)
		}
	case <-signalChan: // シグナルが送信された場合
		log.Println("Received os.Signal. Initiating graceful shutdown...")
	case <-ctx.Done(): // 正常に終了する場合
	}

	cancel()  // コンテキストをキャンセルしてメトリクス更新の goroutine を停止
	wg.Wait() // メトリクス更新 goroutine の終了を待機

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
