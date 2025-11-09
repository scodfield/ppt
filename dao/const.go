package dao

import (
	"context"
	"time"
)

const (
	UserLockKey                  = "ppt:user:lock:%d"             // 用户Redis锁
	UserCacheKey                 = "ppt:user:cache:%d"            // 用户UserCache
	UserFuncSwitchKey            = "ppt:user:func_switch:%d"      // 用户功能开关缓存
	UserLoginTimeQueueKey        = "ppt:user:login_time_queue:%d" // 用户最近登录时间
	UserSettleSetKey             = "ppt:user:settle_set"
	UserNameRegisterKey          = "ppt:user:name_registered" // 用户名注册
	UserActiveKey                = "ppt:user:active:%s"       // 活跃用户
	MongoDBPPT                   = "ppt"
	MongoCollUsers               = "users"
	MongoCollUserCredit          = "user_credit"
	MongoCollUserLogin           = "user_login"
	MongoCollFriendVisit         = "friend_visit"
	MongoCollIPReg               = "ip_reg" // IP注冊表
	UserMailExpiredKey           = "ppt:user:mail_expired"
	UserMailExpiredKeyExpire     = 6 * time.Hour
	UserMailExpiredDeleteDays    = 7                  // 过期删除天数
	UserMailExpiredMaxDeleteDays = 15                 // 最大过期删除天数
	UserIDKey                    = "ppt:user:user_id" // 用户UserID key
	UserIDMin                    = 100000000          // 最小UserID
	UserIDMax                    = 999999999          // 最大UserID
)

var (
	Ctx                        = context.Background()
	UserLoginTimeQueueMax      = 5 // 最近5次登录
	UserCacheDefaultExpiration = time.Minute * 30
	UserCacheDefaultCleanUp    = time.Minute * 30
	UserMailExpiredDeleteBatch = 20000
)
