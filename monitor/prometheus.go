package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"runtime"
	"time"
)

var (
	counter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "actor2",
		Name:      "thd_op_count_test",
		Help:      "The total number of processed events",
	})
	gauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "actor2",
		Name:      "thd_gauge_test",
		Help:      "The total number of gauge",
	})
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "req_duration_seconds",
			Help:    "request process duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"})
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "req_count",
			Help: "count of requests",
		},
		[]string{"path"})
	RequestStatusCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "req_status_count",
			Help: "count of requests status code",
		},
		[]string{"path", "code"})
	RequestMethodCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "req_method_count",
			Help: "count of requests method",
		},
		[]string{"path", "method"})
)

// StartPrometheus 启动监控
func StartPrometheus() {
	reg := prometheus.NewRegistry()
	initPromRegistry(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	opProcessedMetrics()
	server := http.NewServeMux()
	server.Handle("/metrics", promHandler)
	http.ListenAndServe(":8081", server)
}

func initPromRegistry(reg *prometheus.Registry) {
	reg.MustRegister(counter)
	reg.MustRegister(gauge)
	reg.MustRegister(RequestDuration)
	reg.MustRegister(RequestCount)
	reg.MustRegister(RequestStatusCount)
	reg.MustRegister(RequestMethodCount)
	reg.MustRegister(collectors.NewGoCollector())
}

func opProcessedMetrics() {
	go func() {
		for {
			counter.Inc()
			gauge.Add(float64(runtime.NumGoroutine()))
			time.Sleep(10 * time.Second)
		}
	}()
}
