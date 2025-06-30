package db

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
	"time"
)

const (
	TaskInfoCacheKey    = "asynq:task:cache:%s"
	TaskInfoCacheExpire = time.Hour * 24 * 30
)

func SetAsynqTaskCache(client redis.UniversalClient, key string, value []byte) (bool, error) {
	return client.SetNX(dao.Ctx, key, value, TaskInfoCacheExpire).Result()
}

func GetAsynqTaskCache(client redis.UniversalClient, key string) ([]byte, error) {
	result := client.Get(dao.Ctx, key)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return nil, nil
		}
		logger.Error("GetAsynqTaskCache Get error", zap.String("asynq_task_key", key), zap.Error(result.Err()))
		return nil, result.Err()
	}
	return result.Bytes()
}

func DelAsynqTaskCache(client redis.UniversalClient, key string) error {
	return client.Del(dao.Ctx, key).Err()
}
