package main

import (
	"context"
	"log"
	"net/http"
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
	pushInterval = 3 * time.Second // 3 秒毎に PushGateway にメトリクスを送信する例
	lifeTime     = 1 * time.Minute // 1 分後に Job を終了
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
	log.Println("CronJob starting...")

	// create context
	ctx, cancel := context.WithTimeout(context.Background(), lifeTime)
	defer cancel()

	// Collector を定義
	collector := pushmetric.NewCollector()

	// カスタムメトリクスを追加
	collector.WithDefaultMetrics().WithCustomMetrics(bytesSentCounter, bytesRecvCounter) // 任意のメトリクスを追加できる
	// // もしデフォルトのメトリクスを使用しない場合
	// collector.WithCustomMetrics(bytesSentCounter, bytesRecvCounter)

	// エクスポートメトリクス情報を登録する
	// config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector)

	// カスタムクライアントを使用する場合
	client := WithCustomClient()
	config := pushmetric.New(jobName, applicationName, pushInterval, pushGatewayEndPoint, collector).WithClient(client)

	errCh := make(chan error, 1)
	go func() {
		errCh <- config.RoutineSequentialExporter(ctx) // PushGateway にシーケンシャルにメトリクスをエクスポートする
	}()

	// メトリクスを任意のタイミングで更新する
	go func() {
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
				bytesSent, bytesRecv := monitorNetworkSpeed()
				bytesSentCounter.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(bytesSent)
				bytesRecvCounter.WithLabelValues(applicationNameLabelValue, instanceNameLavelValue).Set(bytesRecv)
			}
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-ctx.Done():
		log.Println("CronJob completed")
	}
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
