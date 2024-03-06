package main

import (
	"log"
	"net/http"

	"github.com/GotoRen/metrics-exporter-playground/app/internal"
	"github.com/GotoRen/metrics-exporter-playground/app/model"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
)

const pushGatewayEndPoint = "http://pushgateway:9091"

func main() {
	e := echo.New()

	internal.Register(prometheus.DefaultRegisterer)

	/*** Pull metrics ***/
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

	/*** Push metrics ***/
	e.GET("/push", func(c echo.Context) error {
		internal.MetricsCounter().Inc()

		if err := push.New(pushGatewayEndPoint, "eyes_metrics_counter_job").
			Collector(internal.MetricsCounter()).Push(); err != nil {

			errMsg := "Could not push to Pushgateway"
			log.Println(errMsg, err)

			return c.String(http.StatusInternalServerError, errMsg)
		}

		return c.String(http.StatusOK, "Hello, eyes!!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
