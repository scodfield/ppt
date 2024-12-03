package db

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

const MaxRetryTimes = 10
const (
	DISTRIBUTED_LOCK_USERID  = "dis_lock_uid"
	DISTRIBUTED_LOCK_DBUFFER = "dis_lock_dbuffer"
)

func GetDistributedLock(lockKey string) (bool, int64) {
	for i := 0; i < MaxRetryTimes; i++ {
		newExpireTime := time.Now().Unix() + 10
		if locked, _ := redisC.SetNX(ctx, lockKey, newExpireTime, 10*time.Second).Result(); locked {
			return true, newExpireTime
		}

		oldExpireTimeStr, err := redisC.Get(ctx, lockKey).Result()
		if errors.Is(err, redis.Nil) {
			continue
		} else if err != nil {
			log.Println("Fail to get uid_lock:", err)
			return false, 0
		}
		oldExpireTime, _ := strconv.ParseInt(oldExpireTimeStr, 10, 64)
		if oldExpireTime > time.Now().Unix() {
			continue
		}

		newExpireTime = time.Now().Unix() + 10
		currentExpireTimeStr, _ := redisC.GetSet(ctx, lockKey, newExpireTime).Result()
		currentExpireTime, _ := strconv.ParseInt(currentExpireTimeStr, 10, 64)
		if currentExpireTime == newExpireTime {
			return true, newExpireTime
		}
	}
	return false, 0
}

func DelDistributedLock(lockKey string, expireTime int64) {
	currentExpireTimeStr, err := redisC.Get(ctx, lockKey).Result()
	if errors.Is(err, redis.Nil) {
		return
	} else if err != nil {
		log.Println("Fail to get uid_lock:", err)
		return
	}
	currentExpireTime, _ := strconv.ParseInt(currentExpireTimeStr, 10, 64)
	if currentExpireTime == expireTime {
		redisC.Del(ctx, lockKey)
	}
}
