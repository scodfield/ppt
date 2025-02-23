package db

import (
	"gorm.io/gorm"
	"ppt/pg/model"
)

type UserMailDao struct {
	db *gorm.DB
}

func NewUserMailDao(db *gorm.DB) *UserMailDao {
	return &UserMailDao{db: db}
}

// GetUserMailByID 获取指定邮件
func (m *UserMailDao) GetUserMailByID(id string) (*model.UserMail, error) {
	userMail := &model.UserMail{}
	if err := m.db.First(userMail, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return userMail, nil
}

// ReadMail 读取邮件
func (m *UserMailDao) ReadMail(id string) error {
	if err := m.db.Model(&model.UserMail{}).Where("id = ?", id).Update("read_status", model.MailReadStatusRead).Error; err != nil {
		return err
	}
	return nil
}

// ReceiveMailAccessory 领取邮件附件
func (m *UserMailDao) ReceiveMailAccessory(id string) error {
	if err := m.db.Model(&model.UserMail{}).Where("id = ?", id).Updates(model.UserMail{ReadStatus: model.MailReadStatusRead, AccessoryStatus: model.MailAccessoryStatusReceived}).Error; err != nil {
		return err
	}
	return nil
}

// GetUserMails 获取用户所有可见邮件(未删除/未撤销/未过期)
func (m *UserMailDao) GetUserMails(userID uint64) ([]*model.UserMail, error) {
	var userMails []*model.UserMail
	if err := m.db.Where("user_id = ? and visible_status = ?", userID, model.MailVisibleStatusDefault).Find(&userMails).Error; err != nil {
		return nil, err
	}
	return userMails, nil
}
