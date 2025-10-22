package mq

import (
	"crypto/tls"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"ppt/log"
	"ppt/nacos/wrapper"
	"sync"
	"syscall"
)

var (
	asynqServer *asynq.Server
	once        sync.Once
	sigs        = make(chan os.Signal, 1)
)

func InitAsynqServer(redisConfig *wrapper.RedisConfig) error {
	var err error
	once.Do(func() {
		redisAddr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)

		var tlsConfig *tls.Config
		if redisConfig.SSLVerify {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		var clientOpt asynq.RedisConnOpt
		clientOpt = asynq.RedisClientOpt{
			Addr:      redisAddr,
			Password:  redisConfig.Password,
			DB:        TaskDBNum,
			TLSConfig: tlsConfig,
		}
		if redisConfig.IsCluster {
			clientOpt = asynq.RedisClusterClientOpt{
				Addrs:     []string{redisAddr},
				Password:  redisConfig.Password,
				TLSConfig: tlsConfig,
			}
		}

		cfg := asynq.Config{
			Queues: map[string]int{
				TaskQueueTypeInstant: 6,
				TaskQueueTypeLatency: 3,
			},
		}
		srv := asynq.NewServer(clientOpt, cfg)
		if err = srv.Ping(); err != nil {
			log.Error("ppt InitAsynqServer ping redis server error", zap.Error(err))
			return
		}
		asynqServer = srv
	})
	return nil
}

func StartAsynqServer() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(PPTTaskType, HandlePptTask)
	if err := asynqServer.Start(mux); err != nil {
		log.Error("asynq server Start error", zap.Error(err))
		return
	}
	log.Info("asynq server started")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	_ = <-sigs
	log.Info("asynq server shutting down")
	asynqServer.Stop()
	asynqServer.Shutdown()
	asynqServer = nil
}

func CloseAsynqServer() {
	if asynqServer != nil {
		sigs <- syscall.SIGQUIT
	}
}
