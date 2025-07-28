package dao

import "context"

const (
	UserLockKey                = "ppt:user:lock:%d"             // 用户Redis锁
	UserCacheKey               = "ppt:user:cache:%d"            // 用户UserCache
	UserFuncSwitchKey          = "ppt:user:func_switch:%d"      // 用户功能开关缓存
	UserLoginTimeQueueKey      = "ppt:user:login_time_queue:%d" // 用户最近登录时间
	UserSettleSetKey           = "ppt:user:settle_set"
	MongoDBPTT                 = "ppt"
	MongoCollUsers             = "users"
	MongoCollUserCredit        = "user_credit"
	MongoCollUserLogin         = "user_login"
	MongoCollFriendVisit       = "friend_visit"
	UserMailsExpiredDeleteDays = 7 // 过期删除天数
)

var (
	Ctx                   = context.Background()
	UserLoginTimeQueueMax = 5 // 最近5次登录
)
