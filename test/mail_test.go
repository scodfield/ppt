package test

import (
	"encoding/json"
	"go.uber.org/zap"
	"ppt/code"
	"ppt/dao"
	"ppt/dao/db"
	"ppt/log"
	"ppt/model"
	"testing"
	"time"
)

func TestUserMail(t *testing.T) {
	err := model.MigrateMailTemplate(dao.PgDB)
	if err != nil {
		log.Error("TestUserMail MigrateMailTemplate error", zap.Error(err))
		return
	}

	err = model.MigrateUserMail(dao.PgDB)
	if err != nil {
		log.Error("TestUserMail MigrateUserMail error", zap.Error(err))
		return
	}

	senderID := uint64(101)
	senderNameMap := map[string]string{
		"zh": "系统邮件",
		"en": "System Mail",
	}
	senderName, _ := json.Marshal(senderNameMap)
	titleMap := map[string]string{
		"zh": "邮件标题1",
		"en": "Mail Title 1",
	}
	title, _ := json.Marshal(titleMap)
	contentMap := map[string]string{
		"zh": "欢迎欢迎",
		"en": "Welcome my friends",
	}
	content, _ := json.Marshal(contentMap)
	awardsMap := map[string]interface{}{
		code.CoinTypeGold:   700000,
		code.CoinTypeSilver: 1200000,
		code.CoinTypeCopper: 5000000,
	}
	awards, _ := json.Marshal(awardsMap)
	mailTemplate := &model.MailTemplate{
		SenderID:   senderID,
		SenderName: senderName,
		Title:      title,
		Content:    content,
		Awards:     awards,
		Type:       code.MailTypeSystem,
		Status:     code.MailTemplateStatusInUse,
		ValidDays:  code.MailDefaultValidDays,
	}
	_ = db.NewUserMailDao(dao.PgDB).CreateMailTemplate(mailTemplate)
	log.Info("mail template info", zap.Any("mail_template", mailTemplate))

	userID := uint64(1001)
	mailTemplateID := "0198dc6c-955d-7589-8ac9-cd9bb0a9d51c"
	expiredTime := time.Now().AddDate(0, 0, int(code.MailDefaultValidDays)).UnixMilli()
	userMail := &model.UserMail{
		UserID:          userID,
		TemplateID:      mailTemplateID,
		ExpiredTime:     expiredTime,
		ReadStatus:      code.MailReadStatusUnread,
		AccessoryStatus: code.MailAccessoryStatusUnDraw,
		VisibleStatus:   code.MailVisibleStatusVisible,
		Operator:        code.MailDefaultOperator,
	}
	err = db.NewUserMailDao(dao.PgDB).CreateMailsInBatch([]*model.UserMail{userMail})
	if err != nil {
		log.Error("TestUserMail CreateMailsInBatch error", zap.Error(err), zap.Any("userMail", *userMail))
		return
	}

	userMail2, _ := db.NewUserMailDao(dao.PgDB).GetUserMailByID(userMail.ID)
	log.Info("GetUserMailByID mail info", zap.Any("user_mail", userMail2))
}
