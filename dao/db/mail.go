package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ppt/dao"
	"ppt/logger"
	model2 "ppt/model"
	"time"
)

type UserMailDao struct {
	db *gorm.DB
}

func NewUserMailDao(db *gorm.DB) *UserMailDao {
	return &UserMailDao{db: db}
}

// GetUserMailByID 获取指定邮件
func (m *UserMailDao) GetUserMailByID(id string) (*model2.UserMail, error) {
	userMail := &model2.UserMail{}
	if err := m.db.First(userMail, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return userMail, nil
}

// ReadMail 读取邮件
func (m *UserMailDao) ReadMail(id string) error {
	if err := m.db.Model(&model2.UserMail{}).Where("id = ?", id).Update("read_status", model2.MailReadStatusRead).Error; err != nil {
		return err
	}
	return nil
}

// ReceiveMailAccessory 领取邮件附件
func (m *UserMailDao) ReceiveMailAccessory(id string) error {
	if err := m.db.Model(&model2.UserMail{}).Where("id = ?", id).Updates(model2.UserMail{ReadStatus: model2.MailReadStatusRead, AccessoryStatus: model2.MailAccessoryStatusReceived}).Error; err != nil {
		return err
	}
	return nil
}

// GetUserMails 获取用户所有可见邮件(未删除/未撤销/未过期)
func (m *UserMailDao) GetUserMails(userID uint64) ([]*model2.UserMail, error) {
	var userMails []*model2.UserMail
	if err := m.db.Where("user_id = ? and visible_status = ?", userID, model2.MailVisibleStatusDefault).Find(&userMails).Error; err != nil {
		return nil, err
	}
	return userMails, nil
}

// CreateMailsInBatch 批量插入邮件
func (m *UserMailDao) CreateMailsInBatch(userMail []*model2.UserMail) error {
	return m.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "template_id"},
		},
		DoNothing: true,
	}).CreateInBatches(&userMail, 1000).Error
}

// CreateMailsByFirstOrCreate 插入邮件
func (m *UserMailDao) CreateMailsByFirstOrCreate(userMail []*model2.UserMail) error {
	var err error
	for _, mail := range userMail {
		result := m.db.FirstOrCreate(&mail, model2.UserMail{UserID: mail.UserID, TemplateID: mail.TemplateID})
		if result.Error != nil {
			err = result.Error
			logger.Error("UserMailDao.CreateMailsByFirstOrCreate first_or_create", zap.Error(result.Error))
			continue
		}
	}
	return err
}

// DeleteUserMailsByExpiredTime 删除邮件
func (m *UserMailDao) DeleteUserMailsByExpiredTime(now time.Time) ([]*model2.UserMail, error) {
	var userMails []*model2.UserMail
	expireTime := now.AddDate(0, 0, -dao.UserMailsExpiredDeleteDays).UnixMilli()
	if err := m.db.Clauses(clause.Returning{Columns: []clause.Column{{Name: "user_id"}, {Name: "template_id"}, {Name: "awards"}}}).Where("expired_time <= ?", expireTime).Delete(&userMails).Error; err != nil {
		logger.Error("UserMailDao.DeleteUserMailsByExpiredTime", zap.Error(err))
		return nil, err
	}
	return userMails, nil
}

// DeleteUserMailsByExpiredTimeAndBatch 批量删除邮件
func (m *UserMailDao) DeleteUserMailsByExpiredTimeAndBatch(now time.Time, limit int32) ([]*model2.UserMail, error) {
	var userMails []*model2.UserMail
	expireTime := now.AddDate(0, 0, -dao.UserMailsExpiredDeleteDays).UnixMilli()
	raw := `
		WITH batch_delete AS (
			SELECT id FROM user_mail WHERE expired_time <= ? ORDER BY id DESC LIMIT ?
		)
		DELETE FROM user_mail WHERE id IN (SELECT id FROM batch_delete) RETURNING *
	`
	if err := m.db.Debug().Raw(raw, expireTime, limit).Scan(&userMails).Error; err != nil {
		logger.Error("UserMailDao.DeleteUserMailsByExpiredTimeAndBatch", zap.String("raw_sql", raw), zap.Error(err))
		return nil, err
	}
	return userMails, nil
}
