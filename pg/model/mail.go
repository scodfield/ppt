package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MailTemplate struct {
	BaseModel
	SenderID   uint64         `gorm:"not null;column:sender_id;comment:发送者ID" json:"sender_id"`
	SenderName datatypes.JSON `gorm:"type:jsonb;column:sender_name;comment:发送者" json:"sender_name"`
	Title      datatypes.JSON `gorm:"type:jsonb;column:title;comment:标题" json:"title"`
	Content    datatypes.JSON `gorm:"type:jsonb;column:content;comment:内容" json:"content"`
	Type       int32          `gorm:"not null;column:type;comment:邮件类型" json:"type"`
	Status     int32          `gorm:"not null;column:status;comment:邮件状态" json:"status"`
}

func (MailTemplate) TableName() string {
	return "mail_template"
}

func MigrateMailTemplate(db *gorm.DB) error {
	return db.AutoMigrate(&MailTemplate{})
}
