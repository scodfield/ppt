package mq

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
	"ppt/log"
	"ppt/model"
)

func HandlePptTask(ctx context.Context, t *asynq.Task) error {
	var pptAsynqTask model.PptAsynqTask
	var err error
	if err = json.Unmarshal(t.Payload(), &pptAsynqTask); err != nil {
		log.Error("HandlePptTask json unmarshal t.Payload error", zap.Error(err))
		return err
	}
	log.Info("HandlePptTask receive asynq task", zap.Any("asynq_task", pptAsynqTask))
	/*
		todo handle specific logics
	*/
	return nil
}
