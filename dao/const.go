package dao

import "context"

const (
	UserLockKey                = "ppt:user:lock:%d"  // 用户Redis锁
	UserCacheKey               = "ppt:user:cache:%d" // 用户UserCache
	MongoDB                    = "ppt"
	MongoUsers                 = "users"
	UserMailsExpiredDeleteDays = 7 // 过期删除天数
)

var (
	Ctx = context.Background()
)
