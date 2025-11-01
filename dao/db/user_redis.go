package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/log"
	"ppt/model"
	"strconv"
	"time"
)

// UserLockRedis 用户关键信息强一致性锁
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

func (u *UserLockRedis) DelUserLock(userID uint64, oldExpire int64, du time.Duration) error {
	key := fmt.Sprintf(dao.UserLockKey, userID)
	currExpireStr, err := u.client.Get(dao.Ctx, key).Result()
	if err != nil {
		log.Error("UserLockRedis get error", zap.Uint64("user_id", userID), zap.String("key", key), zap.Error(err))
		return err
	}
	currExpire, err := strconv.ParseInt(currExpireStr, 10, 64)
	if err != nil {
		log.Error("UserLockRedis strconv ParseInt error", zap.Uint64("user_id", userID), zap.String("curr_value", currExpireStr), zap.Error(err))
		return err
	}
	if currExpire != oldExpire {
		// 其它服务已持有锁
		log.Info("UserLockRedis other service has get the lock")
		return nil
	}
	// 有较小的竟态窗口
	newExpireStr, err := u.client.GetDel(dao.Ctx, key).Result()
	if err != nil {
		// 锁已过期
		if errors.Is(err, redis.Nil) {
			return nil
		}
		log.Error("UserLockRedis DelUserLock error", zap.Uint64("user_id", userID), zap.String("key", key), zap.Error(err))
		return err
	}
	newExpire, err := strconv.ParseInt(newExpireStr, 10, 64)
	if err != nil {
		log.Error("UserLockRedis strconv ParseInt error", zap.Uint64("user_id", userID), zap.String("new_value", newExpireStr), zap.Error(err))
		return err
	}
	if newExpire == oldExpire {
		// 安全释放
		return nil
	}
	// 新服务占有锁,将newExpire设置回去,极端情况下可能又有其它服务占有锁,需使用SetNX
	succ, err := u.client.SetNX(dao.Ctx, key, newExpire, du).Result()
	if err != nil {
		log.Error("UserLockRedis SetNX set new lock expire error", zap.Uint64("user_id", userID), zap.Int64("new_value", newExpire), zap.Error(err))
		return err
	}
	if !succ {
		// 极端情况
		return errors.New("set new lock expire failed, other service get lock")
	}
	return nil
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
			log.Error("GetActiveUsers redis scan error", zap.String("redis_key", key), zap.Error(err))
			return nil, err
		}
		for _, v := range result {
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				log.Error("GetActiveUsers error", zap.String("parse_v", v), zap.Error(err))
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
		log.Error("SetActiveUsers redis SADD error", zap.String("redis_key", key), zap.Uint64s("users", users), zap.Error(err))
		return err
	}
	return nil
}

func IsActiveUser(client redis.UniversalClient, key string, userID uint64) (bool, error) {
	exists, err := client.SIsMember(dao.Ctx, key, userID).Result()
	if err != nil {
		log.Error("IsActiveUser redis SIsMember error", zap.String("redis_key", key), zap.Uint64("user_id", userID), zap.Error(err))
		return false, err
	}
	return exists, nil
}

func NewDynamicNotice(client redis.UniversalClient, key string, userID uint64, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Error("NewDynamicNotice json marshal error", zap.String("publish_channel", key), zap.Uint64("user_id", userID), zap.Any("notice_data", data), zap.Error(err))
		return err
	}
	err = client.Publish(dao.Ctx, key, jsonData).Err()
	if err != nil {
		log.Error("NewDynamicNotice publish error", zap.String("publish_channel", key), zap.Uint64("user_id", userID), zap.Error(err))
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
			log.Error("SendUserUpdateStream marshal user update data error", zap.String("stream", stream), zap.Uint64("user_id", userID), zap.Any("user_update_data", update), zap.Error(err))
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
	log.Info("SendUserUpdateStream success to send stream", zap.String("stream", stream), zap.Duration("elapsed", time.Since(begin)))
	return nil
}

func SetUserFuncSwitch(client redis.UniversalClient, userID uint64, funcSwitches []model.UserFuncSwitchT) error {
	key := fmt.Sprintf(dao.UserFuncSwitchKey, userID)
	pipe := client.Pipeline()
	for _, funcSwitch := range funcSwitches {
		data, err := json.Marshal(funcSwitch)
		if err != nil {
			log.Error("SetUserFuncSwitch json marshal error", zap.Any("func_switch", funcSwitch), zap.Error(err))
			return err
		}
		pipe.HSet(dao.Ctx, key, data)
	}
	_, err := pipe.Exec(dao.Ctx)
	if err != nil {
		log.Error("SetUserFuncSwitch pipe exec error", zap.Error(err))
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
		log.Error("PushUserLoginTime pipe exec error", zap.Uint64("user_id", userID), zap.Int64("login_time", loginTime), zap.Error(err))
		return err
	}
	return nil
}

func GetUserLastLoginTime(client redis.UniversalClient, userID uint64) (int64, error) {
	key := fmt.Sprintf(dao.UserLoginTimeQueueKey, userID)
	lastLoginStr, err := client.LIndex(dao.Ctx, key, -1).Result()
	if err != nil {
		log.Error("GetUserLastLoginTime redis LIndex error", zap.Uint64("user_id", userID), zap.Error(err))
		return 0, err
	}
	lastLogin, err := strconv.ParseInt(lastLoginStr, 10, 64)
	if err != nil {
		log.Error("GetUserLastLoginTime ParseInt error", zap.Uint64("user_id", userID), zap.String("last_login_str", lastLoginStr), zap.Error(err))
		return 0, err
	}
	return lastLogin, nil
}

func SetUserSettle(client redis.UniversalClient, settles map[string]int64) error {
	pipe := client.Pipeline()
	for userID, settleMilli := range settles {
		userSettle := redis.Z{Score: float64(settleMilli), Member: userID}
		pipe.ZAddNX(dao.Ctx, dao.UserSettleSetKey, userSettle)
	}
	_, err := pipe.Exec(dao.Ctx)
	if err != nil {
		log.Error("SetUserSettle redis ZAdd error", zap.Any("user_settles", settles), zap.Error(err))
		return err
	}
	return nil
}

func PopUserSettle(client redis.UniversalClient) (uint64, int64, error) {
	result, err := client.ZPopMin(dao.Ctx, dao.UserSettleSetKey, 1).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Info("PopUserSettle no user settle in set.")
			return 0, 0, nil
		}
		log.Error("PopUserSettle redis ZPopMin error", zap.Error(err))
		return 0, 0, err
	}
	if len(result) <= 0 {
		log.Info("PopUserSettle no user settle in set.")
		return 0, 0, nil
	}
	userID, ok := result[0].Member.(uint64)
	if !ok {
		log.Error("PopUserSettle user_id assertion error", zap.Any("user_redis_z", result[0]))
		return 0, 0, errors.New("user_id assertion error")
	}
	settleSec := result[0].Score
	return userID, int64(settleSec), nil
}

// GenerateUserID 生成新的用户UserID
func GenerateUserID(client redis.UniversalClient) (uint64, error) {
	userID, err := client.Incr(dao.Ctx, dao.UserIDKey).Result()
	if err != nil {
		log.Error("GenerateUserID redis Incr error", zap.Error(err))
		return 0, err
	}
	if userID > dao.UserIDMax {
		log.Warn("GenerateUserID redis UserID > max", zap.Int64("user_id", userID))
		return 0, errors.New("user_id out of range")
	}
	return uint64(userID), nil
}

// GenerateMultipleUserIDs 生成多个用户UserID
func GenerateMultipleUserIDs(client redis.UniversalClient, num int) ([]uint64, error) {
	pipe := client.Pipeline()
	results := make([]*redis.IntCmd, num)
	for i := 0; i < num; i++ {
		results[i] = pipe.Incr(dao.Ctx, dao.UserIDKey)
	}
	_, err := pipe.Exec(dao.Ctx)
	if err != nil {
		log.Error("GenerateMultipleUserIDs pipe exec error", zap.Error(err))
		return nil, err
	}

	userIDs := make([]uint64, num)
	for i, result := range results {
		id, err := result.Result()
		if err != nil {
			log.Error("GenerateMultipleUserIDs result.Result() error", zap.Error(err))
			return nil, err
		}
		if id > dao.UserIDMax {
			log.Warn("GenerateMultipleUserIDs result id > max", zap.Int64("user_id", id))
			return nil, errors.New("user_id out of range")
		}
		userIDs[i] = uint64(id)
	}
	return userIDs, nil
}
