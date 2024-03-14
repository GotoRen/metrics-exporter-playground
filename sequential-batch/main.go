package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GotoRen/metrics-exporter-playground/sequential-batch/pushmetric"
)

func main() {
	pushmetric.InitRegister()

	applicationName := "hogehoge"
	jobName := "hoge123"
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second) // 3 秒毎に PushGateway に送信
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				handleApp := &pushmetric.HandleApplication{
					ApplicationName: applicationName,
					InstanceName:    hostname,
				}

				cpuUsage := 0.75
				memoryUsage := 512
				pushmetric.UpdateMetrics(handleApp, cpuUsage, memoryUsage)
				pushmetric.Exports(handleApp, jobName)
			}
		}
	}()

	// Set timer to exit after 1 minute
	exitTimer := time.NewTimer(1 * time.Minute)
	<-exitTimer.C
	fmt.Println("Server shutdown...")
}
