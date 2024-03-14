package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GotoRen/metrics-exporter-playground/sequential-batch/pushmetric"
)

const (
	applicationName = "sample_apps"
	jobName         = "sample_job"
)

const (
	pushInterval = 3 * time.Second // 3 秒毎に PushGateway に送信
	lifeTime     = 1 * time.Minute // 1 分後にJobを終了
)

func main() {
	log.Println("CronJob starting...")

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := pushmetric.ExportConfig(jobName, applicationName, pushInterval)

	// // カスタムクライアントを使用する場合
	// client := NewCustomClient()
	// config := pushmetric.ExportConfig(jobName, applicationName, pushInterval).WithClient(client)

	go config.RuntineSequentialExporter(ctx)

	// wait main routine
	time.Sleep(lifeTime)
	log.Println("CronJob completed")
}

// カスタムクライアントを定義
func NewCustomClient() *http.Client {
	requestTimeoutLimit := 5 * time.Second // 5 秒間 レスポンスがない場合にタイムアウトエラー

	return &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}
}
