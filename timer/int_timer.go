package timer

import (
	"go.uber.org/zap"
	"ppt/config"
	"ppt/dao"
	"ppt/dao/db"
	"ppt/log"
	"time"
)

func InitTimer() error {
	err := initUserMailExpireTimer()
	if err != nil {
		return err
	}
	return nil
}

func initUserMailExpireTimer() error {
	spec := "0 0 3 * * *"
	err := CreateCron("userMailExpireTimer", spec, config.TimeZone, userMailExpire)
	if err != nil {
		log.Error("initUserMailExpireTimer init userMailExpireTimer error", zap.Error(err))
		return err
	}
	return nil
}

func userMailExpire() {
	ok, err := db.GetUserMailExpiredLock(dao.RedisDB, dao.UserMailExpiredKey, dao.UserMailExpiredKeyExpire)
	if err != nil {
		log.Error("userMailExpireTimer get userMailExpiredLock error", zap.Error(err))
		return
	}
	if !ok {
		log.Info("userMailExpireTimer not get userMailExpiredLock")
		return
	}
	go func() {
		doUserMailExpireDelete()
	}()
}

func doUserMailExpireDelete() {
	begin := time.Now()
	flag := true
	for flag {
		delMails, err := db.NewUserMailDao(dao.PgDB).DeleteUserMailsByExpiredTimeAndBatch(time.Now(), int32(dao.UserMailExpiredDeleteBatch))
		if err != nil {
			log.Error("doUserMailExpireDelete delete userMailExpiredLock error", zap.Error(err))
			flag = false
			continue
		}
		if len(delMails) <= 0 {
			log.Info("doUserMailExpireDelete delete success")
			flag = false
			continue
		}
		// todo 归档处理
	}

	deleteCost := time.Since(begin).Seconds()
	log.Info("doUserMailExpireDelete cost seconds", zap.Float64("delete_cost", deleteCost))
}
