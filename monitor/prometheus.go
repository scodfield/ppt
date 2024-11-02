package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "thd_op_count_test",
		Help: "The total number of processed events",
	})
)

// startPrometheus: 启动监控
func StartPrometheus() {
	opProcessedMetrics()
	server := http.NewServeMux()
	server.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", server)
}

func opProcessedMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(1 * time.Second)
		}
	}()
}
