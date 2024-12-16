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
)

// startPrometheus: 启动监控
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
