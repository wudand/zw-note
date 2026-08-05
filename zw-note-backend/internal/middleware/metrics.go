package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests partitioned by method, route, and status code.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "http_request_duration_seconds",
			Help:      "Histogram of HTTP request latencies partitioned by method and route.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "app",
			Name:      "http_requests_in_flight",
			Help:      "Current number of HTTP requests being processed.",
		},
	)
)

// Metrics records Prometheus counters, histograms, and gauges per request.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		httpRequestsInFlight.Inc()

		c.Next()

		httpRequestsInFlight.Dec()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			route,
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			route,
		).Observe(time.Since(start).Seconds())
	}
}
