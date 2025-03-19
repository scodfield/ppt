package dao

import (
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"sync"
	"time"
)

var (
	RedisDB   redis.UniversalClient
	redisOnce sync.Once
)

func InitRedis(redisCfg *RedisConfig) error {
	var err error
	redisOnce.Do(func() {
		url := fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port)
		host, _ := os.Hostname()
		clientName := host + "-" + fmt.Sprintf("%d", os.Getpid())
		if redisCfg.IsClustered {
			client := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:      []string{url},
				Password:   redisCfg.Password,
				Username:   redisCfg.UserName,
				ClientName: clientName,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			})
			_, err = client.Ping(Ctx).Result()
			if err != nil {
				panic(err)
			}
			RedisDB = client
		} else {
			client := redis.NewClient(&redis.Options{
				Addr:            url,
				Password:        redisCfg.Password,
				Username:        redisCfg.UserName,
				DB:              redisCfg.DBIndex,
				ClientName:      clientName,
				MaxActiveConns:  20,
				MaxIdleConns:    10,
				ConnMaxIdleTime: time.Minute * 2,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			})
			_, err = client.Ping(Ctx).Result()
			if err != nil {
				panic(err)
			}
			RedisDB = client
		}
	})
	return err
}

func CloseRedis() {
	if RedisDB != nil {
		_ = RedisDB.Close()
	}
}
