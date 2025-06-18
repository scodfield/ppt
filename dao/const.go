package dao

import "context"

const (
	UserLockKey                = "ppt:user:lock:%d"        // 用户Redis锁
	UserCacheKey               = "ppt:user:cache:%d"       // 用户UserCache
	UserFuncSwitchKey          = "ppt:user:func_switch:%d" // 用户功能开关缓存
	MongoDB                    = "ppt"
	MongoCollUsers             = "users"
	MongoCollFriendVisit       = "friend_visit"
	UserMailsExpiredDeleteDays = 7 // 过期删除天数
)

var (
	Ctx = context.Background()
)
