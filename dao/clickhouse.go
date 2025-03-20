package dao

import (
	"fmt"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"ppt/config"
	"sync"
	"time"
)

var (
	CKSqlSession *gorm.DB
	ckOnce       sync.Once
)

func InitClickHouse(cfg *config.CKConfig) error {
	ckOnce.Do(func() {
		var err error
		CKSqlSession, err = initCKGorm(cfg)
		if err != nil {
			panic(err)
		}

	})

	return nil
}

func initCKGorm(cfg *config.CKConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("http://%s:%s@%s/%s?dial_timeout=10s&read_timeout=20s", cfg.UserName, cfg.Password, cfg.Url, cfg.Database)
	db, err := gorm.Open(clickhouse.New(clickhouse.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Second * 10)
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	db = db.Debug()
	return db, nil
}

func CloseClickHouse() {
	if CKSqlSession != nil {
		sqlDB, _ := CKSqlSession.DB()
		sqlDB.Close()
		CKSqlSession = nil
	}
}
