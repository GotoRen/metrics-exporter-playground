package pushmetric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

// TestNewCollector verifies that a new Collector instance is initialized correctly.
func TestNewCollector(t *testing.T) {
	collector := NewCollector()
	assert.NotNil(t, collector)
	assert.Empty(t, collector.collectors, "Collectors list should be empty initially")
}

// TestWithDefaultMetrics verifies that default metrics are added correctly.
func TestWithDefaultMetrics(t *testing.T) {
	collector := NewCollector()

	// Add default metrics.
	collector.WithDefaultMetrics()

	expectedMetrics := []prometheus.Collector{
		cpuUtilizationMetric,
		memoryUtilizationMetric,
		pushCountMetric,
	}

	assert.Equal(t, len(expectedMetrics), len(collector.collectors), "Collectors list should contain default metrics")

	for _, metric := range expectedMetrics {
		assert.Contains(t, collector.collectors, metric, "Default metric should be in collectors list")
	}
}

// TestWithCustomMetrics verifies that custom metrics are added correctly.
func TestWithCustomMetrics(t *testing.T) {
	customMetric1 := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "custom_metric_1",
			Help: "This is a custom metric 1",
		},
		[]string{"test-app", "test-instance"},
	)

	customMetric2 := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "custom_metric_2",
			Help: "This is a custom metric 2",
		},
		[]string{"test-app", "test-instance"},
	)

	collector := NewCollector()

	// Add any custom metrics.
	collector.WithCustomMetrics(customMetric1, customMetric2)

	assert.Equal(t, 2, len(collector.collectors), "Collectors list should contain custom metrics")
	assert.Contains(t, collector.collectors, customMetric1, "Custom metric 1 should be in collectors list")
	assert.Contains(t, collector.collectors, customMetric2, "Custom metric 2 should be in collectors list")
}
