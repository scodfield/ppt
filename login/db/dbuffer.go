package db

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

const DefaultMaxID int64 = 10000
const DefaultStep int32 = 100

func GetFunctionMaxID(idType string) (int64, int32) {
	var maxID int64 = -1
	lockKey := DISTRIBUTED_LOCK_DBUFFER + idType
	if locked, expireTime := GetDistributedLock(lockKey); locked {
		defer DelDistributedLock(lockKey, expireTime)
		curMaxID, err := redisC.Get(ctx, idType).Result()
		if errors.Is(err, redis.Nil) {
			redisC.Set(ctx, idType, DefaultMaxID+int64(DefaultStep), 0)
			return DefaultMaxID, DefaultStep
		} else if err != nil {
			log.Fatal("Get max id error: ", err)
			return -1, 0
		}
		maxID, _ = strconv.ParseInt(curMaxID, 10, 64)
		redisC.Set(ctx, idType, maxID+DefaultMaxID, 0)
	}
	return maxID, DefaultStep
}
