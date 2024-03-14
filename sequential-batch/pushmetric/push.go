package pushmetric

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

const pushGatewayEndPoint = "http://localhost:9091"

// Export pushes metrics to Prometheus Pushgateway with the given job name.
func Export(ctx context.Context, client *http.Client, jobName string) {
	if err := pushMetrics(ctx, pushGatewayEndPoint, jobName, client, RegisterMetrics()...); err != nil {
		log.Println("Failed to push metrics:", err)
	} else {
		log.Println("Metrics pushed successfully")
	}
}

// pushMetrics pushes the provided collectors to the specified endpoint with the given job name.
func pushMetrics(ctx context.Context, endpoint string, jobName string, client *http.Client, collectors ...prometheus.Collector) error {
	pusher := push.New(endpoint, jobName).Client(client)

	for _, collector := range collectors {
		pusher = pusher.Collector(collector)
	}

	return pusher.PushContext(ctx)
}

// GetInstanceName returns the hostname of the current instance.
func GetInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname:", err)
	}

	return hostname
}
