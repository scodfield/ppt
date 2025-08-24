package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// MailTemplate 邮件模板
type MailTemplate struct {
	BaseModel
	SenderID   uint64         `gorm:"not null;column:sender_id;comment:发送者ID" json:"sender_id"`
	SenderName datatypes.JSON `gorm:"not null;type:jsonb;column:sender_name;comment:发送者" json:"sender_name"`
	Title      datatypes.JSON `gorm:"not null;type:jsonb;column:title;comment:标题" json:"title"`
	Content    datatypes.JSON `gorm:"not null;type:jsonb;column:content;comment:内容" json:"content"`
	Awards     datatypes.JSON `gorm:"type:jsonb;column:awards;comment:附件奖励" json:"awards"`
	Type       int32          `gorm:"not null;column:type;comment:邮件类型" json:"type"`
	Status     int32          `gorm:"not null;column:status;comment:邮件状态" json:"status"`
	ValidDays  int32          `gorm:"column:valid_days;comment:有效天数" json:"valid_days"`
}

func (MailTemplate) TableName() string {
	return "mail_template"
}

func MigrateMailTemplate(db *gorm.DB) error {
	return db.AutoMigrate(&MailTemplate{})
}

// UserMail 用户邮件
type UserMail struct {
	BaseModel
	UserID          uint64                  `gorm:"not null;column:user_id;comment:用户UserID" json:"user_id"`
	TemplateID      string                  `gorm:"not null;column:template_id;comment:邮件模板ID" json:"template_id"`
	Awards          datatypes.JSON          `gorm:"column:awards;comment:奖励附件" json:"awards"`
	ExpiredTime     int64                   `gorm:"not null;column:expired_time;comment:过期时间" json:"expired_time"`
	ReadStatus      MailReadStatusType      `gorm:"not null;type:integer;default:0;column:read_status;comment:读取状态" json:"read_status"`
	AccessoryStatus MailAccessoryStatusType `gorm:"not null;type:integer;default:0;column:accessory_status;comment:附件状态" json:"accessory_status"`
	VisibleStatus   MailVisibleStatusType   `gorm:"not null;type:integer;default:0;column:visible_status;comment:可见性状态" json:"visible_status"`
	CreateTime      int64                   `gorm:"autoCreateTime:milli;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime      int64                   `gorm:"autoUpdateTime:milli;column:update_time;comment:更新时间" json:"update_time"`
	Operator        string                  `gorm:"not null;column:operator;comment:操作人" json:"operator"`
}

func (UserMail) TableName() string {
	return "user_mail"
}

func MigrateUserMail(db *gorm.DB) error {
	return db.AutoMigrate(&UserMail{})
}
