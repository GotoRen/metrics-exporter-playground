package pushmetric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Collector
		expectedLen int
	}{
		{
			name: "Test WithDefaultMetrics",
			setup: func() *Collector {
				c := NewCollector()
				return c
			},
			expectedLen: 3, // Expecting three default metrics
		},
		{
			name: "Test WithCustomMetrics",
			setup: func() *Collector {
				c := NewCollector()
				cm1 := prometheus.NewCounter(prometheus.CounterOpts{Name: "custom_metric_1"})
				cm2 := prometheus.NewCounter(prometheus.CounterOpts{Name: "custom_metric_2"})
				c.WithCustomMetrics(cm1, cm2)
				return c
			},
			expectedLen: 2, // Expecting two custom metrics
		},
		{
			name: "Test NewCollector",
			setup: func() *Collector {
				c := NewCollector()
				return c
			},
			expectedLen: 0, // Expecting empty collectors slice initially
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Setup
			c := test.setup()

			// Test
			collectorLen := len(c.collectors)

			// Assertion
			assert.Equal(t, test.expectedLen, collectorLen)
		})
	}
}
