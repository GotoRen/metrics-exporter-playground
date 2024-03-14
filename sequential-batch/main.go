package main

import (
	"log"
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

				pushmetric.Export(jobName)
			}
		}
	}()

	exitTimer := time.NewTimer(lifeTime)
	<-exitTimer.C

	log.Println("CronJob completed")
}
