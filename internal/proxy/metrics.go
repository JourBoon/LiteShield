package proxy

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	reqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "liteshield_requests_total",
			Help: "Total des requêtes par client et status",
		},
		[]string{"client", "status", "cache"},
	)

	latencyHist = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "liteshield_latency_ms",
			Help:    "Latence des requêtes (ms)",
			Buckets: prometheus.LinearBuckets(1, 5, 10),
		},
		[]string{"client"},
	)
)

func init() {
	prometheus.MustRegister(reqTotal, latencyHist)
}

func recordMetrics(clientID string, status int, cacheHit bool, duration time.Duration) {
	cacheLabel := "miss"
	if cacheHit {
		cacheLabel = "hit"
	}
	reqTotal.WithLabelValues(clientID, toStatusLabel(status), cacheLabel).Inc()
	latencyHist.WithLabelValues(clientID).Observe(float64(duration.Milliseconds()))
}

func toStatusLabel(status int) string {
	if status >= 500 {
		return "5xx"
	}
	if status >= 400 {
		return "4xx"
	}
	return "2xx"
}
