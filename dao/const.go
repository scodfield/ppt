package dao

import "context"

const (
	UserLockKey  = "ppt:user:lock:%d"  // 用户Redis锁
	UserCacheKey = "ppt:user:cache:%d" // 用户UserCache
)

var (
	Ctx = context.Background()
)
