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
var redisC, pubSubRedis *redis.Client
var redisOnce sync.Once

func GetRedis() *redis.Client {
	return redisC
}

func GetPubSubRedis() *redis.Client {
	return pubSubRedis
}

func InitRedis() *redis.Client {
	redisOnce.Do(func() {
		redisC, _ = initRedisClint("127.0.0.1", 6379, "", 0, 10)
		pubSubRedis, _ = initRedisClint("127.0.0.1", 6379, "", 0, 0)
	})

	initSubscribe()
	fmt.Println("redis connect success.")
	return redisC
}

func initRedisClint(host string, port int, password string, db, poolSize int) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	r := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  time.Second * 5,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		PoolSize:     poolSize,
		PoolTimeout:  time.Second * 5,
	})
	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Fail to connect redis:", err)
		return nil, err
	}
	return r, nil
}

func GenerateUserID() (newID int64) {
	if locked, expireTime := GetDistributedLock(DISTRIBUTED_LOCK_USERID); locked {
		newID, _ = redisC.Incr(context.Background(), "uid").Result()
		DelDistributedLock(DISTRIBUTED_LOCK_USERID, expireTime)
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

func initSubscribe() {
	sub := &subscriber{
		channels: []string{REDIS_CHANNEL_LOGIN_NOTICE},
	}
	go subloop(sub)
}

func subloop(sub *subscriber) {
	go func() {
		for {
			ps := pubSubRedis.Subscribe(ctx, sub.Channels()...)
			closed := false
			for {
				msg, err := ps.ReceiveMessage(ctx)
				if err != nil {
					closed = sub.OnError(err)
					break
				}
				sub.OnMessage(msg)
			}
			_ = ps.Unsubscribe(ctx, sub.Channels()...)
			_ = ps.Close()
			if closed {
				break
			}
		}
	}()
}

const (
	REDIS_CHANNEL_LOGIN_NOTICE = "login_notice"
	REDIS_STREAM_TEST_STREAM   = "test_stream"
)

type subscriber struct {
	channels []string
}

func (s *subscriber) Channels() []string {
	return s.channels
}

func (s *subscriber) OnError(err error) bool {
	log.Println("redis pub sub error: ", err)
	if errors.Is(err, redis.ErrClosed) {
		return true
	}
	return false
}

func (s *subscriber) OnMessage(msg *redis.Message) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Fatal("redis on pub_sub message error: ", err)
			}
		}()

		switch msg.Channel {
		case REDIS_CHANNEL_LOGIN_NOTICE:
			log.Printf("recv login_notice, payload:%s\n", msg.Payload)
		default:
			log.Println("redis on pub_sub message error: unknown channel: ", msg.Channel)
		}
	}()
}

// TestStream
func TestStream() {
	redisC.XAdd(ctx, &redis.XAddArgs{
		Stream: REDIS_STREAM_TEST_STREAM,
		Values: map[string]interface{}{
			"name": "ppt_001",
		},
	})
}
