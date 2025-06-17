package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MailReadStatusType 邮件读取状态定义
type MailReadStatusType int32

const (
	MailReadStatusUnRead MailReadStatusType = 0  // 未读
	MailReadStatusRead   MailReadStatusType = 10 // 已读
)

// MailAccessoryStatusType 邮件附件状态定义
type MailAccessoryStatusType int32

const (
	MailAccessoryStatusUnReceive MailAccessoryStatusType = 0  // 附件未领取
	MailAccessoryStatusReceived  MailAccessoryStatusType = 10 // 已领取
)

// MailVisibleStatusType 邮件可见性状态定义
type MailVisibleStatusType int32

const (
	MailVisibleStatusDefault MailVisibleStatusType = 0  // 默认状态（初始创建）
	MailVisibleStatusDeleted MailVisibleStatusType = 10 // 已删除
	MailVisibleStatusRevoke  MailVisibleStatusType = 20 // 已撤销（后台撤销）
	MailVisibleStatusExpired MailVisibleStatusType = 30 // 已过期
)

// BaseModel 通用model
type BaseModel struct {
	ID        string `gorm:"type:uuid;primary_key" json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (b *BaseModel) BeforeCreate(db *gorm.DB) error {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return err
	}
	b.ID = uuidV7.String()
	return nil
}

type UserFuncSwitchT struct {
	UserID   uint64
	SwitchID int32 // 开关ID
	IsActive bool  // 是否开放
}
