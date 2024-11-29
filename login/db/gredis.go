package db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

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
