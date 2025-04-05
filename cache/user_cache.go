package cache

import (
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"ppt/dao"
	"ppt/model"
	"time"
)

var (
	UserCache *Cache[uint64, any]
)

type UserCacheT struct {
	pgSql *gorm.DB
	redis redis.UniversalClient
}

func NewUserCache(pgSql *gorm.DB, redis redis.UniversalClient, defaultExpiration, cleanupExpiration time.Duration) {
	UserCache = NewCache[uint64, any](defaultExpiration, cleanupExpiration, &UserCacheT{
		pgSql: pgSql,
		redis: redis,
	}, false)
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
