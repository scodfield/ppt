package db

import (
	"errors"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"ppt/cache"
	"ppt/dao"
	"ppt/logger"
	"ppt/model"
	"time"
)

// SetAccountInfo 设置用户账号信息
func SetAccountInfo(name string, user model.User) error {
	//user.AccId = get_account_id()
	// set cache info & mark registered
	//SetLoginCache(user)
	//redisCache.Command("HSET", regHash, name, 1)
	//_, insertErr := o.Insert(&user)
	//if insertErr != nil {
	//	fmt.Println("mys insert err, ", insertErr)
	//}

	var err error
	now := time.Now().UnixMilli()
	exists := false
	if exists, err = dao.RedisDB.HSetNX(dao.Ctx, dao.UserNameRegisterKey, name, now).Result(); err != nil {
		logger.Error("SetAccountInfo HSetNX user name error", zap.Uint64("user_id", user.UserID), zap.String("user_name", name))
		return err
	}
	if !exists {
		logger.Warn("SetAccountInfo HSetNX user name already exists", zap.Uint64("user_id", user.UserID), zap.String("user_name", name))
		return errors.New("user name already exists")
	}

	return nil
}

// WhetherUserNameRegistered 用户Name是否已注册(true-已注册,false-未注册)
func WhetherUserNameRegistered(name string) bool {
	result := dao.RedisDB.HGet(dao.Ctx, dao.UserNameRegisterKey, name)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.ErrNil) {
			return false
		}
		logger.Error("WhetherUserNameRegistered redis HGet error", zap.String("user_name", name), zap.Error(result.Err()))
		return true
	}
	regMilli, err := result.Int64()
	if err != nil {
		logger.Error("WhetherUserNameRegistered value assert error", zap.String("result_value", result.Val()), zap.Error(err))
		return true
	}
	if regMilli > 0 {
		return true
	}
	return false
}

// SetUserCache 设置用户缓存
func SetUserCache(user model.User) error {
	cache.UserCache.Set(user.UserID, user, dao.UserCacheDefaultExpiration)
	return nil
}

// GetUserCache 获取用户缓存
func GetUserCache(userID uint64) (*model.User, error) {
	userAny, err := cache.UserCache.Get(userID)
	if err != nil {
		logger.Error("GetUserCache Get user cache error", zap.Uint64("user_id", userID), zap.Error(err))
		return nil, err
	}
	userCache := userAny.(model.User)
	return &userCache, nil
}
