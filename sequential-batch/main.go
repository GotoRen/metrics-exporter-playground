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
	requestTimeoutLimit = 5 * time.Second // 5 秒間 レスポンスがない場合にタイムアウトエラー
	pushInterval        = 3 * time.Second // 3 秒毎に PushGateway に送信
	lifeTime            = 1 * time.Minute // 1 分後にJobを終了
)

func main() {
	log.Println("CronJob starting...")

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create custom http client
	client := &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	// export metrics
	// TODO: Http client を実装する場合とそうでない場合を分離する
	go pushmetric.RoutineSequentialExport(ctx, client, jobName, applicationName, pushInterval)

	// wait main routine
	time.Sleep(lifeTime)
	log.Println("CronJob completed")
}
