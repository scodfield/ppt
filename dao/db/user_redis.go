package db

import (
	"encoding/json"
	"errors"
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

func PushUserLoginTime(client redis.UniversalClient, userID uint64, loginTime int64) error {
	key := fmt.Sprintf(dao.UserLoginTimeQueueKey, userID)
	pipe := client.Pipeline()
	pipe.LPush(dao.Ctx, key, loginTime)
	pipe.LTrim(dao.Ctx, key, 0, int64(dao.UserLoginTimeQueueMax))
	_, err := pipe.Exec(dao.Ctx)
	if err != nil {
		logger.Error("PushUserLoginTime pipe exec error", zap.Uint64("user_id", userID), zap.Int64("login_time", loginTime), zap.Error(err))
		return err
	}
	return nil
}

func GetUserLastLoginTime(client redis.UniversalClient, userID uint64) (int64, error) {
	key := fmt.Sprintf(dao.UserLoginTimeQueueKey, userID)
	lastLoginStr, err := client.LIndex(dao.Ctx, key, -1).Result()
	if err != nil {
		logger.Error("GetUserLastLoginTime redis LIndex error", zap.Uint64("user_id", userID), zap.Error(err))
		return 0, err
	}
	lastLogin, err := strconv.ParseInt(lastLoginStr, 10, 64)
	if err != nil {
		logger.Error("GetUserLastLoginTime ParseInt error", zap.Uint64("user_id", userID), zap.String("last_login_str", lastLoginStr), zap.Error(err))
		return 0, err
	}
	return lastLogin, nil
}

func SetUserSettle(client redis.UniversalClient, userID uint64, settleSec int64) error {
	userSettle := redis.Z{Score: float64(settleSec), Member: userID}
	_, err := client.ZAddNX(dao.Ctx, dao.UserSettleSetKey, userSettle).Result()
	if err != nil {
		logger.Error("SetUserSettle redis ZAdd error", zap.Uint64("user_id", userID), zap.Int64("settle_sec", settleSec), zap.Error(err))
		return err
	}
	return nil
}

func PopUserSettle(client redis.UniversalClient) (uint64, error) {
	result, err := client.ZPopMin(dao.Ctx, dao.UserSettleSetKey, 1).Result()
	if err != nil {
		logger.Error("PopUserSettle redis ZPopMin error", zap.Error(err))
		return 0, err
	}
	userID, ok := result[0].Member.(uint64)
	if !ok {
		logger.Error("PopUserSettle user_id assertion error", zap.Any("user_redis_z", result[0]))
		return 0, errors.New("user_id assertion error")
	}
	return userID, nil
}
