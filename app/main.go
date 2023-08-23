package main

import (
	"github.com/GotoRen/metrics-exporter-playground/app/internal"
	"github.com/GotoRen/metrics-exporter-playground/app/model"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := echo.New()

	internal.Register(prometheus.DefaultRegisterer)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	h1 := &model.HandleApplication{
		ApplicationName: "hoge-test-1",
		DashboardUID:    "uid_100",
	}

	h2 := &model.HandleApplication{
		ApplicationName: "hoge-test-2",
		DashboardUID:    "uid_200",
	}

	internal.UpdateEyesCustomMetric(h1)
	internal.UpdateEyesCustomMetric(h2)

	e.Start(":8080")
}
