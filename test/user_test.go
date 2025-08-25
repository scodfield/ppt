package test

import (
	"go.uber.org/zap"
	"ppt/dao"
	"ppt/dao/db"
	"ppt/log"
	"testing"
)

func TestUserRegIP(t *testing.T) {
	regIP := "127.0.0.1"
	userID := uint64(1002)
	_ = db.UpdateIPReg(dao.MongoClient, userID, regIP)

	ipRegInfo, _ := db.GetIPReg(dao.MongoClient, regIP)
	log.Info("get ip reg info", zap.Any("reg_info", ipRegInfo))
}
