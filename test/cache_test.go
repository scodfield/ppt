package test

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"ppt/log"
	"ppt/login/db"
	"ppt/model"
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
	_ = db.SetUserCache(userCache)

	userCache2, _ := db.GetUserCache(userID)
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
