package internal

import (
	"github.com/GotoRen/metrics-exporter-playground/app/model"
	"github.com/prometheus/client_golang/prometheus"
)

// eyes_custom_metric メトリクス の ラベル名
const (
	applicationNameKey = "application_name"
	dashboardUidKey    = "dashboard_uid"
)

var (
	eyesCustomMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eyes_custom_metric",
			Help: "Application metrics managed by eyes",
		},
		[]string{applicationNameKey, dashboardUidKey},
	)
)

func Register(r prometheus.Registerer) {
	r.MustRegister(
		eyesCustomMetric,
	)
}

func UpdateEyesCustomMetric(h *model.HandleApplication) {
	eyesCustomMetric.WithLabelValues(h.ApplicationName, h.DashboardUID)
}
