package middleware

import (
	"github.com/gin-gonic/gin"
	"ppt/monitor"
	"strconv"

	"time"
)

func Prom() gin.HandlerFunc {
	return func(c *gin.Context) {
		begin := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		cost := time.Since(begin)
		monitor.RequestDuration.WithLabelValues(path).Observe(cost.Seconds())
		monitor.RequestCount.WithLabelValues(path).Inc()
		monitor.RequestStatusCount.WithLabelValues(path, strconv.Itoa(c.Writer.Status())).Inc()
		monitor.RequestMethodCount.WithLabelValues(path, method).Inc()
	}
}
