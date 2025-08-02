package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"ppt/logger"
	"ppt/middleware"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func (s *HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Fatal("http server shutdown error", zap.Error(err))
	}
}

func NewHttpServer(port int) *HttpServer {
	r := initRouter()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}
	return &HttpServer{server: httpServer}
}

func initRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{}))
	router.Use(middleware.GinRecover(&logger.Logger, true))
	router.Use(middleware.Prom(), middleware.Cors())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	router.GET("/metrics", gin.WrapF(middleware.HttpCheckIP([]string{}, promhttp.Handler()).ServeHTTP))
	return router
}
