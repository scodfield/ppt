package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

var ctx = context.Background()
var redisC *redis.Client
var redisOnce sync.Once

func GetRedis() *redis.Client {
	return redisC
}

func InitRedis() *redis.Client {
	redisOnce.Do(func() {
		redisC = redis.NewClient(&redis.Options{
			Addr:         "127.0.0.1:6379",
			Password:     "",
			DB:           0,
			DialTimeout:  time.Second * 5,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			PoolSize:     10,
			PoolTimeout:  time.Second * 5,
		})
	})
	pong, err := redisC.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Fail to connect redis:", err)
		return nil
	}
	fmt.Printf("redis connect success, pong: %s\n", pong)
	return redisC
}

func GenerateUserID() (newID int64) {
	var locked bool
	var count int
	for !locked {
		count++
		locked, _ = redisC.SetNX(ctx, "uid_lock", "locked", 10*time.Second).Result()
		if locked || count >= 10 {
			break
		}
	}
	if locked {
		newID, _ = redisC.Incr(context.Background(), "uid").Result()
		redisC.Del(ctx, "uid_lock")
	}
	return
}

func UpdateToken(id int64, token string) {
	redisC.Set(ctx, formatTokenKey(id), token, 24*time.Hour)
}

func formatTokenKey(id int64) string {
	return fmt.Sprintf("token_%v", id)
}

func IsTokenOutOfDate(id int64) bool {
	_, err := redisC.Get(ctx, formatTokenKey(id)).Result()
	if errors.Is(err, redis.Nil) {
		return false
	} else if err != nil {
		log.Fatal("Fail to get token:", err)
		return false
	}
	return true
}
