package mq

import (
	"crypto/tls"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
)

const (
	TaskDBNum = 15
)

var (
	asynqClient    *asynq.Client
	asynqInspector *asynq.Inspector
)

func InitAsynq(redisCfg *dao.RedisConfig) (res error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("InitAsynq panic recover", zap.Any("err", err))
			res = err.(error)
		}
	}()

	redisAddr := fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port)
	var tlsCfg *tls.Config
	if redisCfg.SSLVerify {
		tlsCfg = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	clientOpts := asynq.RedisClientOpt{
		Addr:      redisAddr,
		Password:  redisCfg.Password,
		DB:        TaskDBNum,
		TLSConfig: tlsCfg,
	}
	asynqClient = asynq.NewClient(clientOpts)
	if err := asynqClient.Ping(); err != nil {
		logger.Error("InitAsynq ping fail", zap.Any("err", err))
		return err
	}
	asynqInspector = asynq.NewInspector(clientOpts)

	return nil
}

func CloseAsynq() {
	if asynqClient != nil {
		if err := asynqClient.Close(); err != nil {
			logger.Error("CloseAsynq close asynq client fail", zap.Any("err", err))
		}
		asynqClient = nil
	}
	if asynqInspector != nil {
		if err := asynqInspector.Close(); err != nil {
			logger.Error("CloseAsynq close inspector fail", zap.Any("err", err))
		}
		asynqInspector = nil
	}
}
