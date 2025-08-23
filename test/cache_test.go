package test

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"ppt/dao"
	"ppt/dao/db"
	"ppt/log"
	logindb "ppt/login/db"
	"ppt/model"
	"ppt/util"
	"testing"
	"time"
)

func TestUserCache(t *testing.T) {
	userID := uint64(1001)
	uuid, _ := uuid2.NewV7()
	now := time.Now().UnixMilli()
	userCache := model.User{
		ID:        uuid.String(),
		UserID:    userID,
		Username:  "ppt_001",
		Password:  "tdv23d8rf",
		Email:     "ppt_001@qq.com",
		BrandID:   1,
		Channel:   "ppt_thx",
		Lang:      "zh",
		CreatedAt: now,
		UpdateAt:  now,
	}
	_ = logindb.SetUserCache(userCache)

	userCache2, _ := logindb.GetUserCache(userID)
	fmt.Printf("userCache2: %+v\n", userCache2)

}

func TestNetDial(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 10*time.Second)
	if err != nil {
		log.Error("TestNetDial DialTimeout error", zap.Error(err))
		return
	}
	defer conn.Close()
	log.Info("TestNetDial success", zap.Any("conn", conn))

	_, err = conn.Write([]byte("PING\r\n"))
	if err != nil {
		log.Error("TestNetDial write ping error", zap.Error(err))
		return
	}

	buffer := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		log.Error("TestNetDial read ping error", zap.Error(err))
		return
	}
	log.Info("TestNetDial success", zap.Any("buffer", buffer[:n]))
}

func TestActiveUser(t *testing.T) {
	now := time.Now()
	userID := uint64(1001)
	activeKey := fmt.Sprintf(dao.UserActiveKey, util.TimeToDateStr(now))
	if err := db.SetActiveUsers(dao.RedisDB, activeKey, []uint64{userID}); err != nil {
		log.Error("TestActiveUser SetActiveUsers error", zap.Error(err))
		return
	}

	isActive, err := db.IsActiveUser(dao.RedisDB, activeKey, userID)
	if err != nil {
		log.Error("TestActiveUser IsActiveUser error", zap.Error(err))
		return
	}
	fmt.Printf("isActive: %+v\n", isActive)
}
