package db

import (
	"errors"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ppt/cache"
	"ppt/dao"
	"ppt/log"
	"ppt/model"
	"time"
)

// RegAccountInfo 注册用户账号
func RegAccountInfo(name string, user model.User) error {
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
		log.Error("SetAccountInfo HSetNX user name error", zap.Uint64("user_id", user.UserID), zap.String("user_name", name))
		return err
	}
	if !exists {
		log.Warn("SetAccountInfo HSetNX user name already exists", zap.Uint64("user_id", user.UserID), zap.String("user_name", name))
		return errors.New("user name already exists")
	}

	return nil
}

// WhetherUserNameRegistered 用户Name是否已注册(true-已注册,false-未注册)
func WhetherUserNameRegistered(name string) (bool, error) {
	exists, err := dao.RedisDB.SIsMember(dao.Ctx, dao.UserNameRegisterKey, name).Result()
	if err != nil {
		log.Error("WhetherUserNameRegistered SIsMember error", zap.Error(err))
		return false, err
	}
	return exists, nil
}

// RegUserName 注册账户名
func RegUserName(name string) error {
	tx := func(tx *redis.Tx) error {
		exists, err := tx.SIsMember(dao.Ctx, dao.UserNameRegisterKey, name).Result()
		if err != nil {
			log.Error("RegUserName SIsMember error", zap.Error(err))
			return err
		}
		if exists {
			log.Warn("RegUserName SIsMember already exists", zap.String("user_name", name))
			return errors.New("user name already exists")
		}
		_, err = tx.SAdd(dao.Ctx, dao.UserNameRegisterKey, name).Result()
		if err != nil {
			log.Error("RegUserName SAdd  error", zap.String("user_name", name), zap.Error(err))
			return err
		}
		return nil
	}

	err := dao.RedisDB.Watch(dao.Ctx, tx, dao.UserNameRegisterKey)
	if err != nil {
		log.Error("RegUserName Watch error", zap.String("user_name", name), zap.Error(err))
		return err
	}
	return nil
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
		log.Error("GetUserCache Get user cache error", zap.Uint64("user_id", userID), zap.Error(err))
		return nil, err
	}
	userCache := userAny.(model.User)
	return &userCache, nil
}
