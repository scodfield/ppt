package db

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
	"ppt/model"
	"strconv"
	"time"
)

type UserLockRedis struct {
	client redis.UniversalClient
}

func NewUserLockRedis(client redis.UniversalClient) *UserLockRedis {
	return &UserLockRedis{client: client}
}

func (u *UserLockRedis) SetUserLock(userID uint64, expire int64, du time.Duration) (bool, error) {
	key := fmt.Sprintf(dao.UserLockKey, userID)
	return u.client.SetNX(dao.Ctx, key, expire, du).Result()
}

func GetActiveUsers(client redis.UniversalClient, key string) ([]uint64, error) {
	var cursor uint64
	var users []uint64
	scanSize := 10000
	for {
		var err error
		var result []string
		result, cursor, err = client.SScan(dao.Ctx, key, cursor, "*", int64(scanSize)).Result()
		if err != nil {
			logger.Error("GetActiveUsers redis scan error", zap.String("redis_key", key), zap.Error(err))
			return nil, err
		}
		for _, v := range result {
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				logger.Error("GetActiveUsers error", zap.String("parse_v", v), zap.Error(err))
				continue
			}
			users = append(users, id)
		}
		if cursor == 0 {
			break
		}
	}
	return users, nil
}

func SetActiveUsers(client redis.UniversalClient, key string, users []uint64) error {
	usersInterface := make([]interface{}, len(users))
	for i, v := range users {
		usersInterface[i] = v
	}
	_, err := client.SAdd(dao.Ctx, key, usersInterface...).Result()
	if err != nil {
		logger.Error("SetActiveUsers redis SADD error", zap.String("redis_key", key), zap.Uint64s("users", users), zap.Error(err))
		return err
	}
	return nil
}

func IsActiveUser(client redis.UniversalClient, key string, userID uint64) (bool, error) {
	exists, err := client.SIsMember(dao.Ctx, key, userID).Result()
	if err != nil {
		logger.Error("IsActiveUser redis SIsMember error", zap.String("redis_key", key), zap.Uint64("user_id", userID), zap.Error(err))
		return false, err
	}
	return exists, nil
}

func NewDynamicNotice(client redis.UniversalClient, key string, userID uint64, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("NewDynamicNotice json marshal error", zap.String("publish_channel", key), zap.Uint64("user_id", userID), zap.Any("notice_data", data), zap.Error(err))
		return err
	}
	err = client.Publish(dao.Ctx, key, jsonData).Err()
	if err != nil {
		logger.Error("NewDynamicNotice publish error", zap.String("publish_channel", key), zap.Uint64("user_id", userID), zap.Error(err))
		return err
	}
	return nil
}

func SendUserUpdateStream(client redis.UniversalClient, stream string, updates map[uint64]interface{}) error {
	begin := time.Now()
	pipe := client.Pipeline()
	for userID, update := range updates {
		data, err := json.Marshal(update)
		if err != nil {
			logger.Error("SendUserUpdateStream marshal user update data error", zap.String("stream", stream), zap.Uint64("user_id", userID), zap.Any("user_update_data", update), zap.Error(err))
			return err
		}
		pipe.XAdd(dao.Ctx, &redis.XAddArgs{
			Stream: stream,
			Values: map[string]interface{}{
				"data":    data,
				"user_id": userID,
			},
			MaxLen: 10000,
		})
	}
	_, err := pipe.Exec(dao.Ctx)
	if err != nil {
		return err
	}
	logger.Info("SendUserUpdateStream success to send stream", zap.String("stream", stream), zap.Duration("elapsed", time.Since(begin)))
	return nil
}

func SetUserFuncSwitch(client redis.UniversalClient, userID uint64, funcSwitches []model.UserFuncSwitchT) error {
	key := fmt.Sprintf(dao.UserFuncSwitchKey, userID)
	pipe := client.Pipeline()
	for _, funcSwitch := range funcSwitches {
		data, err := json.Marshal(funcSwitch)
		if err != nil {
			logger.Error("SetUserFuncSwitch json marshal error", zap.Any("func_switch", funcSwitch), zap.Error(err))
			return err
		}
		pipe.HSet(dao.Ctx, key, data)
	}
	_, err := pipe.Exec(dao.Ctx)
	if err != nil {
		logger.Error("SetUserFuncSwitch pipe exec error", zap.Error(err))
		return err
	}
	return nil
}
