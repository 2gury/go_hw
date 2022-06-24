package middleware

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitoring struct {
	Hits       *prometheus.CounterVec
	Duration   *prometheus.HistogramVec
	HitsVk     *prometheus.CounterVec
	DurationVk *prometheus.HistogramVec
}

func NewMonitoring(mx *echo.Echo) *Monitoring {
	monitoring := &Monitoring{
		Hits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hits",
				Help: "number of requests",
			},
			[]string{"method", "path", "status"},
		),
		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "duration",
				Help: "duration of requests",
			},
			[]string{"method", "path", "status"},
		),
		HitsVk: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hits_vk",
				Help: "number of requests to vk services",
			},
			[]string{"method", "url", "status"},
		),
		DurationVk: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "duration_vk",
				Help: "duration of requests to vk services",
			},
			[]string{"method", "url", "status"},
		),
	}

	prometheus.MustRegister(monitoring.Hits, monitoring.Duration, monitoring.HitsVk, monitoring.DurationVk)
	mx.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	return monitoring
}
