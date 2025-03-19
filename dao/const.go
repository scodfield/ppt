package dao

import "context"

const (
	UserLockKey = "ppt:user:lock:%d" // 用户Redis锁
)

var (
	Ctx = context.Background()
)
