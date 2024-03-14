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
	requestTimeoutLimit = time.Second * 5 // 5 秒間 レスポンスがない場合にタイムアウトエラー
	pushInterval        = 3 * time.Second // 3 秒毎に PushGateway に送信
	lifeTime            = 1 * time.Minute // 1 分後にJobを終了
)

func main() {
	log.Println("CronJob starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &http.Client{
		Timeout:   requestTimeoutLimit,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	go func() {
		ticker := time.NewTicker(pushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				labels := &pushmetric.CustomLabels{
					ApplicationName: applicationName,
					InstanceName:    pushmetric.GetInstanceName(),
				}

				cpuUsage := 0.85
				memoryUsage := 512

				pushmetric.UpdateGaugeMetric(labels, cpuUsage, memoryUsage)
				pushmetric.IncrementCounterMetric(labels)

				pushmetric.Export(ctx, client, jobName)
			}
		}
	}()

	exitTimer := time.NewTimer(lifeTime)
	<-exitTimer.C

	log.Println("CronJob completed")
}
