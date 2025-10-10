package model

import (
	"github.com/google/uuid"
)

// CouponType 优惠券类型枚举
type CouponType int

const (
	CouponTypeFullReduction CouponType = iota * 10 // 0 = 满减
	CouponTypeFullGift                             // 10 = 满赠
)

// CouponStatus 优惠券状态枚举
type CouponStatus int

const (
	CouponStatusAvailable   CouponStatus = iota * 100 // 0 = 可用
	CouponStatusUsedSuccess                           // 100 = 成功使用
	CouponStatusUsedFailed                            // 110 = 使用失败
	CouponStatusExpired                               // 200 = 已过期
	CouponStatusDeleted                               // 300 = 已删除
)

// UserCoupon 用户优惠券表结构
type UserCoupon struct {
	BaseModel
	UserID       uint64       `gorm:"column:user_id;type:bigint;not null;primaryKey", json:"user_id"`
	CouponID     uuid.UUID    `gorm:"column:coupon_id;type:uuid;not null;primaryKey", json:"coupon_id"`
	CouponType   CouponType   `gorm:"column:coupon_type;type:integer;not null", json:"coupon_type"`
	CouponStatus CouponStatus `gorm:"column:coupon_status;type:integer;not null;default:0", json:"coupon_status"`
	ReceivedAt   int64        `gorm:"column:received_at;not null;", json:"received_at"`
	ExpiredAt    int64        `gorm:"column:expired_at;not null", json:"expired_at"`
}

// TableName 设置表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}
