package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
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
	UserLoginCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_login_count",
			Help: "count of user login through different modes",
		},
		[]string{"login_mode"})
)

func InitProm() {
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestStatusCount)
	prometheus.MustRegister(RequestMethodCount)
	prometheus.MustRegister(UserLoginCount)
	prometheus.MustRegister(collectors.NewGoCollector())
}
