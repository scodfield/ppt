package cache

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"ppt/dao"
	"ppt/model"
)

var (
	UserCache *Cache[uint64, any]
)

type UserCacheT struct {
	pgSql *gorm.DB
	redis redis.UniversalClient
}

func InitUserCache() error {
	var err error
	err = initUserCache()
	return err
}

func initUserCache() error {
	UserCache = NewCache[uint64, any](dao.UserCacheDefaultExpiration, dao.UserCacheDefaultCleanUp, &UserCacheT{
		pgSql: dao.PgDB,
		redis: dao.RedisDB,
	}, false)
	return nil
}

func (u *UserCacheT) Load(userID uint64) (interface{}, error) {
	var user model.User
	key := fmt.Sprintf(dao.UserCacheKey, userID)
	if cacheBytes, err := u.redis.Get(dao.Ctx, key).Bytes(); err == nil {
		err = json.Unmarshal(cacheBytes, &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	return model.GetUserByID(u.pgSql, userID)
}
