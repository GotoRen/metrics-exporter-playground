package pushmetric

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

/*** GaugeMetric ***/
// SetGaugeMetric sets gauge metrics with given values.
func SetGaugeMetric(collector prometheus.Collector, value float64, lvs ...string) {
	// gv := collector.(prometheus.GaugeVec)
	// gv.WithLabelValues(lvs...)

	if gv, ok := collector.(*prometheus.GaugeVec); ok {
		gv.WithLabelValues(lvs...).Set(value)
	} else {
		// 型アサーション失敗
		fmt.Println("collector is not of type *prometheus.GaugeVec")
	}
}

func IncrementGaugeMetric(collector prometheus.Collector, lvs ...string) {
	gv := collector.(prometheus.GaugeVec)
	gv.WithLabelValues(lvs...).Inc()
}

func DecrementGaugeMetric(collector prometheus.Collector, lvs ...string) {
	gv := collector.(prometheus.GaugeVec)
	gv.WithLabelValues(lvs...).Dec()
}

func AddGaugeMetric(collector prometheus.Collector, value float64, lvs ...string) {
	gv := collector.(prometheus.GaugeVec)
	gv.WithLabelValues(lvs...).Add(value)
}

func SubGaugeMetric(collector prometheus.Collector, value float64, lvs ...string) {
	gv := collector.(prometheus.GaugeVec)
	gv.WithLabelValues(lvs...).Sub(value)
}

func SetToCurrentTimeGaugeMetric(collector prometheus.Collector, value float64, lvs ...string) {
	gv := collector.(prometheus.GaugeVec)
	gv.WithLabelValues(lvs...).SetToCurrentTime()
}

/*** CounterMetrics ***/
// IncrementCounterMetric increments counter metric.
func IncrementCounterMetric(collector prometheus.Collector, lvs ...string) {
	// cv := collector.(prometheus.CounterVec)
	// cv.WithLabelValues(lvs...).Inc()

	if cv, ok := collector.(prometheus.CounterVec); ok {
		cv.WithLabelValues(lvs...).Inc()
	} else {
		// 型アサーション失敗
		fmt.Println("collector is not of type *prometheus.CounterVec")
	}

}

func DecrementCounterMetric(collector prometheus.Collector, value float64, lvs ...string) {
	cv := collector.(prometheus.CounterVec)
	cv.WithLabelValues(lvs...).Add(value)
}
