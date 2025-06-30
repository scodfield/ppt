package mq

import (
	"crypto/tls"
	"fmt"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
	"time"
)

const (
	TaskDBNum            = 15
	TaskQueueTypeInstant = "task_queue_instant" // 实时队列
	TaskQueueTypeLatency = "task_queue_latency" // 延时队列
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

// EnqueueTaskInstant 实时发送
func EnqueueTaskInstant(task *asynq.Task) error {
	info, err := asynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Queue(TaskQueueTypeInstant))
	if err != nil {
		logger.Error("EnqueueTaskInstant enqueue fail", zap.Any("asynq_task", task), zap.Any("err", err))
		return err
	}
	logger.Info("EnqueueTaskInstant enqueue success", zap.Any("info", info))
	return nil
}

// EnqueueTaskLatency 延时发送
func EnqueueTaskLatency(task *asynq.Task, sendTime time.Time) error {
	info, err := asynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Queue(TaskQueueTypeLatency), asynq.ProcessAt(sendTime))
	if err != nil {
		logger.Error("EnqueueTaskLatency enqueue fail", zap.Any("asynq_task", task), zap.Any("err", err))
		return err
	}
	logger.Info("EnqueueTaskLatency enqueue success", zap.Any("info", info))
	return nil
}
