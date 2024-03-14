package pushmetric

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const pushGatewayEndPoint = "http://localhost:9091"

// Export pushes metrics to Prometheus Pushgateway with the given job name.
func Export(jobName string) {
	if err := pushMetrics(pushGatewayEndPoint, jobName, RegisterMetrics()...); err != nil {
		log.Println("Failed to push metrics:", err)
	} else {
		log.Println("Metrics pushed successfully")
	}
}

// pushMetrics pushes the provided collectors to the specified endpoint with the given job name.
func pushMetrics(endpoint string, jobName string, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName)

	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.Push()
}

// GetInstanceName returns the hostname of the current instance.
func GetInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	return hostname
}
