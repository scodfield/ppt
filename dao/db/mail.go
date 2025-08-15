package db

import (
	"database/sql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"ppt/dao"
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
			log.Error("UserMailDao.CreateMailsByFirstOrCreate first_or_create", zap.Error(result.Error))
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
		log.Error("UserMailDao.DeleteUserMailsByExpiredTime", zap.Error(err))
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
		log.Error("UserMailDao.DeleteUserMailsByExpiredTimeAndBatch", zap.String("raw_sql", raw), zap.Error(err))
		return nil, err
	}
	return userMails, nil
}

// UpdateUserMailByTx UserMail事务更新
func (m *UserMailDao) UpdateUserMailByTx(userID uint64, updates map[string]interface{}) error {
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("UpdateUserMailByTx recover panic", zap.Uint64("user_id", userID), zap.Any("mail_updates", updates), zap.Any("recover_panic", r), zap.Stack("stack"))
			tx.Rollback()
		}
	}()

	tx.Model(&model2.UserMail{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).Updates(updates)
	if err := tx.Commit().Error; err != nil {
		log.Error("UpdateUserMailByTx commit update error", zap.Error(err))
		// 一般会自动回滚,安全起见
		tx.Rollback()
		return err
	}
	return nil
}

// UpdateUserMailByTransaction UserMail自动事务更新
func (m *UserMailDao) UpdateUserMailByTransaction(userID uint64, updates map[string]interface{}) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model2.UserMail{}).Where("user_id = ?", userID).Updates(updates).Error; err != nil {
			log.Error("UpdateUserMailByTransaction updates error", zap.Error(err))
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}
