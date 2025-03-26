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
	UserCache *UserCacheT
)

type UserCacheT struct {
	pgSql *gorm.DB
	redis redis.UniversalClient
}

func NewUserCache(pgSql *gorm.DB, redis redis.UniversalClient) *UserCacheT {
	return &UserCacheT{
		pgSql: pgSql,
		redis: redis,
	}
}

func (u *UserCacheT) Load(userID uint64) (*model.User, error) {
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
