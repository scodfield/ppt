package dao

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"os"
	"ppt/log"
	"ppt/nacos/wrapper"
	"sync"
	"time"
)

var (
	RedisDB   redis.UniversalClient
	redisOnce sync.Once
)

func InitRedis(redisCfg *wrapper.RedisConfig) error {
	var err error
	url := fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port)
	host, _ := os.Hostname()
	clientName := host + "-" + fmt.Sprintf("%d", os.Getpid())
	if redisCfg.IsCluster {
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
		r := redis.NewClient(&redis.Options{
			Addr:            "127.0.0.1:6379",
			Password:        redisCfg.Password,
			Username:        redisCfg.UserName,
			DB:              redisCfg.DBIndex,
			ClientName:      clientName,
			MaxActiveConns:  20,
			MaxIdleConns:    10,
			ConnMaxIdleTime: time.Minute * 2,
			//TLSConfig:       &tls.Config{
			//	//InsecureSkipVerify: false,
			//},
			DialTimeout:  5 * time.Second, // 建立连接的超时时间
			ReadTimeout:  5 * time.Second, // 读超时
			WriteTimeout: 5 * time.Second, // 写超时
			PoolTimeout:  5 * time.Second, // 连接池获取连接的超时时间
		})
		ctx := context.Background()
		pingCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
		defer cancel()
		result := r.Ping(pingCtx)
		if result.Err() != nil {
			log.Error("InitRedis ping error", zap.Error(result.Err()))
			panic(err)
		}
		RedisDB = r
	}
	err = initUserID()
	return err
}

func initUserID() error {
	exists, err := RedisDB.Exists(Ctx, UserIDKey).Result()
	if err != nil {
		log.Error("initUserID RedisDB.Exists error", zap.Error(err))
		return err
	}
	if exists == 0 {
		if err = RedisDB.Set(Ctx, UserIDKey, UserIDMin, 0).Err(); err != nil {
			log.Error("initUserID RedisDB.Set error", zap.Error(err))
			return err
		}
	}
	return nil
}

func CloseRedis() {
	if RedisDB != nil {
		_ = RedisDB.Close()
		RedisDB = nil
	}
}
