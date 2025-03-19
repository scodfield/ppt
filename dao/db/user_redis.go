package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"ppt/dao"
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
