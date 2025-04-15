package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/logger"
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
			logger.Error("GetActiveUsers redis scan error", zap.String("scan_key", key), zap.Error(err))
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
